package webui

import (
	"fmt"
	"time"

	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"

	"github.com/influxdata/influxdb/client/v2"
)

// NewAPICfgIfxServer API for IfxServer Catalog Management
func NewAPICfgIfxServer(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/ifxserver", func() {
		m.Get("/", reqSignedIn, GetIfxServer)
		m.Post("/", reqSignedIn, bind(config.IfxServerCfg{}), AddIfxServer)
		m.Put("/:id", reqSignedIn, bind(config.IfxServerCfg{}), UpdateIfxServer)
		m.Delete("/:id", reqSignedIn, DeleteIfxServer)
		m.Get("/:id", reqSignedIn, GetIfxServerCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetIfxServerAffectOnDel)
		m.Post("/ping", reqSignedIn, bind(config.IfxServerCfg{}), PingIfxServer)
		m.Post("/import", reqSignedIn, bind(config.IfxServerCfg{}), ImportIfxCatalog)
	})

	return nil
}

// ping internal
func ping(cfg *config.IfxServerCfg) (client.Client, time.Duration, string, error) {

	conf := client.HTTPConfig{
		Addr:      cfg.URL,
		Username:  cfg.AdminUser,
		Password:  cfg.AdminPasswd,
		UserAgent: "resistor",
		//		Timeout:   time.Duration(cfg.Timeout) * time.Second,
	}
	cli, err := client.NewHTTPClient(conf)

	if err != nil {
		return cli, 0, "", err
	}
	elapsed, message, err := cli.Ping(time.Duration(5) * time.Second)
	if err == nil {
		_, err = InfluxQuery(cli, "show databases", "", 0)
	}
	return cli, elapsed, message, err
}

// InfluxQuery does IQL queries
func InfluxQuery(cli client.Client, query string, db string, colidx int) ([]string, error) {
	defer cli.Close()
	q := client.NewQuery(query, db, "s")
	response, err := cli.Query(q)
	if err != nil {
		log.Errorf("Query RESULT ERROR: %+s", err)
		return nil, err
	}
	if response.Error() != nil {
		log.Errorf("Query RESULT ERROR: %+s", response.Error())
		return nil, fmt.Errorf("Query RESULT ERROR: %+s", response.Error())
	}
	var retarray []string
	for _, v1 := range response.Results {
		for _, v2 := range v1.Series {
			for _, v3 := range v2.Values {
				retarray = append(retarray, v3[colidx].(string))
			}
		}
	}
	return retarray, nil
}

func getInfluxDBs(cli client.Client) []string {
	ret, err := InfluxQuery(cli, "SHOW DATABASES", "", 0)
	if err != nil {
		log.Errorf("GetInfluxDB's Query: %s", err)
	}
	return ret
}

func getDBMeasurements(cli client.Client, db string) []string {
	ret, err := InfluxQuery(cli, "SHOW MEASUREMENTS", db, 0)
	if err != nil {
		log.Errorf("Get Measurements Query: %s", err)
	}
	return ret
}

func getDBRetention(cli client.Client, db string) []string {
	ret, err := InfluxQuery(cli, "SHOW RETENTION POLICIES ON \""+db+"\"", "", 0)
	if err != nil {
		log.Errorf("Get Retention Policies Query: %s", err)
	}
	return ret
}

func getMeasurementsFields(cli client.Client, db string, m string) []string {

	ret, err := InfluxQuery(cli, "SHOW FIELD KEYS FROM \""+m+"\"", db, 0)
	if err != nil {
		log.Errorf("Get Measurements Fields Query: %s", err)
	}
	return ret
}

func getMeasurementsTags(cli client.Client, db string, m string) []string {
	ret, err := InfluxQuery(cli, "SHOW TAG KEYS FROM \""+m+"\"", db, 0)
	if err != nil {
		log.Errorf("Get Measurements Tags Query: %s", err)
	}
	return ret
}

func getMeasurementTagsByKey(cli client.Client, db string, m string, tagkey string) []string {
	s := "SHOW TAG VALUES FROM \"" + m + "\" WITH KEY = \"" + tagkey + "\""
	ret, err := InfluxQuery(cli, s, db, 1)
	if err != nil {
		log.Errorf("getMeasurementTagsByKey. Error getting Measurement Tags with query: %s. Error: %s", s, err)
	}
	log.Debugf("getMeasurementTagsByKey. DB: %s. Query: %s. Returning: %+v", db, s, ret)
	return ret
}

// ImportIfxCatalog Imports Influx Catalog to the internal database
func ImportIfxCatalog(ctx *Context, dev config.IfxServerCfg) {
	log.Infof("Importing catalog for Server: %s", dev.ID)
	cli, _, _, err := ping(&dev)
	if err != nil {
		log.Errorf("Error on Ping Influx Server %s: Err: %s", dev.ID, err)
		ctx.JSON(404, err.Error())
		return
	}

	dbs := getInfluxDBs(cli)
	for _, db := range dbs {
		rps := getDBRetention(cli, db)
		meas := getDBMeasurements(cli, db)

		var itemarray []*config.ItemComponent
		for _, m := range meas {
			//msmtid=ifxserver-dbname-msmtname
			msmtid := dev.ID + "-" + db + "-" + m
			tags := getMeasurementsTags(cli, db, m)
			fields := getMeasurementsFields(cli, db, m)
			mcfg := config.IfxMeasurementCfg{ID: msmtid, Name: m, Tags: tags, Fields: fields}
			log.Infof("ImportIfxCatalog. Got DATABASE [%s] with retention policy [%s] MEASUREMENT %s/%s TAGS [%+v] FIELDS [%+v]", db, rps, mcfg.ID, m, tags, fields)
			itemarray = append(itemarray, &config.ItemComponent{ID: mcfg.ID, Name: m})
		}

		//dbid=ifxserver-dbname
		dbid := dev.ID + "-" + db
		log.Debugf("ImportIfxCatalog. AddOrUpdateIfxDBCfg with dbid: %s", dbid)
		dbcfg := config.IfxDBCfg{ID: dbid, Name: db, IfxServer: dev.ID, Retention: rps, Measurements: itemarray}
		_, err := agent.MainConfig.Database.AddOrUpdateIfxDBCfg(&dbcfg)
		if err != nil {
			log.Errorf("ImportIfxCatalog. Error on Importing Influx DBs: %s Err: %s", dev.ID, err)
			ctx.JSON(404, err.Error())
			return
		}
		//Insert measurements
		for _, m := range meas {
			//msmtid=ifxserver-dbname-msmtname
			msmtid := dev.ID + "-" + db + "-" + m
			tags := getMeasurementsTags(cli, db, m)
			fields := getMeasurementsFields(cli, db, m)
			mcfg := config.IfxMeasurementCfg{ID: msmtid, Name: m, Tags: tags, Fields: fields}
			log.Debugf("ImportIfxCatalog. AddOrUpdateIfxMeasurementCfg with msmtid: %s", msmtid)
			_, err := agent.MainConfig.Database.AddOrUpdateIfxMeasurementCfg(&mcfg)
			if err != nil {
				log.Errorf("ImportIfxCatalog. Error on Importing Influx DBs: %s Err: %s", dev.ID, err)
				ctx.JSON(404, err.Error())
				return
			}
		}
	}

	ctx.JSON(200, &dbs)
}

// PingIfxServer Pings Influx Server and returns the result on context
func PingIfxServer(ctx *Context, dev config.IfxServerCfg) {
	_, elapsed, message, err := ping(&dev)
	if err != nil {
		log.Warningf("Error on Ping Influx Server %s: Err: %s", dev.ID, err)
		ctx.JSON(404, err.Error())
		return
	}
	ctx.JSON(200, &struct {
		Message string
		Elapsed time.Duration
	}{
		Message: message,
		Elapsed: elapsed,
	})

}

// GetIfxServer Returns Influx Server list to frontend
func GetIfxServer(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetIfxServerCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get IfxServer :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting IfxServer %+v", &devcfgarray)
}

// AddIfxServer Inserts Influx Server to the internal database and returns the result on context
func AddIfxServer(ctx *Context, dev config.IfxServerCfg) {
	log.Printf("ADDING IfxServer %+v", dev)
	affected, err := agent.MainConfig.Database.AddIfxServerCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for IfxServer %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateIfxServer Updates Influx Server into the internal database and returns the result on context
func UpdateIfxServer(ctx *Context, dev config.IfxServerCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateIfxServerCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for IfxServer %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteIfxServer Deletes Influx Server from the internal database and returns the result on context
func DeleteIfxServer(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelIfxServerCfg(id)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetIfxServerCfgByID Returns Influx Server to frontend
func GetIfxServerCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetIfxServerCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetIfxServerAffectOnDel Checks if Deletion of Influx Server affects on some other items of the internal database and returns the result on context
func GetIfxServerAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetIfxServerCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for Influx server with id: %s. Error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}

//GetDevicesByProductID Gets the list of devices related to the productid parameter and returns the result on context
func GetDevicesByProductID(ctx *Context) {
	productid := ctx.Params(":productid")
	log.Debugf("GetDevicesByProductID. Entering with productid: %s", productid)
	productcfg, err := agent.MainConfig.Database.GetProductCfgByID(productid)
	if err != nil {
		log.Warningf("GetDevicesByProductID. Error getting product with id: %s. Error: %s", productid, err)
		ctx.JSON(404, err.Error())
	} else {
		ifxdbcfgarray, err := agent.MainConfig.Database.GetIfxDBCfgArrayByMeasName(productcfg.Measurements[0])
		if err != nil {
			log.Warningf("GetDevicesByProductID. Error getting InfluxDB list for measurement name: %s. Error: %s", productcfg.Measurements[0], err)
			ctx.JSON(404, err.Error())
		} else {
			ifxdbcfg := ifxdbcfgarray[0]
			ifxservercfg, err := agent.MainConfig.Database.GetIfxServerCfgByID(ifxdbcfg.IfxServer)
			if err != nil {
				log.Warningf("GetDevicesByProductID. Error getting Influx server with id: %d. Error: %s", ifxdbcfg.IfxServer, err)
				ctx.JSON(404, err.Error())
			} else {
				cli, _, _, err := ping(&ifxservercfg)
				if err != nil {
					log.Warningf("GetDevicesByProductID. Error Pinging Influx Server with id: %d. Error: %s", ifxdbcfg.ID, err)
					ctx.JSON(404, err.Error())
				} else {
					devices := getMeasurementTagsByKey(cli, ifxdbcfg.Name, productcfg.Measurements[0], productcfg.ProductTag)
					ctx.JSON(200, devices)
				}
			}
		}
	}
}

//GetDevicesByAlertID Gets the list of devices related to the alertid parameter and returns the result on context
func GetDevicesByAlertID(ctx *Context) {
	alertid := ctx.Params(":alertid")
	log.Debugf("GetDevicesByAlertID. Entering with alertid: %s", alertid)
	alertidcfg, err := agent.MainConfig.Database.GetAlertIDCfgByID(alertid)
	if err != nil {
		log.Warningf("GetDevicesByAlertID. Error getting alert with id: %s. Error: %s", alertid, err)
		ctx.JSON(404, err.Error())
	} else {
		ifxdbid := alertidcfg.InfluxDB
		ifxdbcfg, err := agent.MainConfig.Database.GetIfxDBCfgByID(ifxdbid)
		if err != nil {
			log.Warningf("GetDevicesByAlertID. Error getting InfluxDB with id: %s. Error: %s", ifxdbid, err)
			ctx.JSON(404, err.Error())
		} else {
			ifxservercfg, err := agent.MainConfig.Database.GetIfxServerCfgByID(ifxdbcfg.IfxServer)
			if err != nil {
				log.Warningf("GetDevicesByAlertID. Error getting Influx server with id: %d. Error: %s", ifxdbcfg.IfxServer, err)
				ctx.JSON(404, err.Error())
			} else {
				cli, _, _, err := ping(&ifxservercfg)
				if err != nil {
					log.Warningf("GetDevicesByAlertID. Error Pinging Influx Server with id: %d. Error: %s", ifxdbcfg.ID, err)
					ctx.JSON(404, err.Error())
				} else {
					devices := getMeasurementTagsByKey(cli, ifxdbcfg.Name, alertidcfg.InfluxMeasurement, alertidcfg.ProductTag)
					ctx.JSON(200, devices)
				}
			}
		}
	}
}
