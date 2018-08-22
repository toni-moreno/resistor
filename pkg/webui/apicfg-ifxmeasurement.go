package webui

import (
	"strconv"
	"strings"

	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgIfxMeasurement API for IfxMeasurement Catalog Management
func NewAPICfgIfxMeasurement(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/ifxmeasurement", func() {
		m.Get("/", reqSignedIn, GetIfxMeasurement)
		m.Post("/", reqSignedIn, bind(config.IfxMeasurementCfg{}), AddIfxMeasurement)
		m.Put("/:id", reqSignedIn, bind(config.IfxMeasurementCfg{}), UpdateIfxMeasurement)
		m.Delete("/:id", reqSignedIn, DeleteIfxMeasurement)
		m.Get("/:id", reqSignedIn, GetIfxMeasurementCfgByID)
		m.Get("/getnames/", reqSignedIn, GetIfxMeasurementCfgDistinctNamesArray)
		m.Get("/gettags/:filter", reqSignedIn, GetIfxMeasurementTagsArray)
		m.Get("/bydbidmeasname/:filter", reqSignedIn, GetIfxMeasurementCfgByDbIDMeasName)
		m.Get("/checkondel/:id", reqSignedIn, GetIfxMeasurementAffectOnDel)
	})

	return nil
}

// GetIfxMeasurement Return IfxMeasurement list to frontend
func GetIfxMeasurement(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetIfxMeasurementCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get IfxMeasurements :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting IfxMeasurements %+v", &devcfgarray)
}

// GetIfxMeasurementCfgDistinctNamesArray Return IfxMeasurement list with different names to frontend
func GetIfxMeasurementCfgDistinctNamesArray(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetIfxMeasurementCfgDistinctNamesArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get IfxMeasurements :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting IfxMeasurements %+v", &devcfgarray)
}

// AddIfxMeasurement Insert new snmpdevice to de internal BBDD --pending--
func AddIfxMeasurement(ctx *Context, dev config.IfxMeasurementCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddIfxMeasurementCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateIfxMeasurement --pending--
func UpdateIfxMeasurement(ctx *Context, dev config.IfxMeasurementCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	nid, err := strconv.ParseInt(id, 10, 64)
	affected, err := agent.MainConfig.Database.UpdateIfxMeasurementCfg(nid, &dev)
	if err != nil {
		log.Warningf("Error on update for device %d/%s  , affected : %+v , error: %s", dev.ID, dev.Name, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteIfxMeasurement delete from the catalog database
func DeleteIfxMeasurement(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	nid, err := strconv.ParseInt(id, 10, 64)
	affected, err := agent.MainConfig.Database.DelIfxMeasurementCfg(nid)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetIfxMeasurementCfgByID --pending--
func GetIfxMeasurementCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Getting Influx Measurements with id: '%s'.", id)
	nid, err := strconv.ParseInt(id, 10, 64)
	dev, err := agent.MainConfig.Database.GetIfxMeasurementCfgByID(nid)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetIfxMeasurementCfgByDbIDMeasName Gets an Influx Measurement with all its information from DbID and Measurement Name
func GetIfxMeasurementCfgByDbIDMeasName(ctx *Context) {
	filter := ctx.Params(":filter")
	log.Debugf("Getting Influx Measurements with filter (DbID&MeasName): '%s'.", filter)
	sqlquery := "select * from ifx_measurement_cfg, ifx_db_meas_rel where ifxmeasid = ifx_measurement_cfg.id "
	if len(filter) > 0 {
		params := strings.Split(filter, "&")
		sqlquery = sqlquery + " and ifxdbid = " + params[0] + " and name = '" + params[1] + "'"
	}
	sqlquery = sqlquery + " order by name, id"
	dev, err := agent.MainConfig.Database.GetIfxMeasurementCfgBySQLQuery(sqlquery)
	if err != nil {
		log.Warningf("Error getting Measurement for sqlquery %s. Error: %s", sqlquery, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

/*GetIfxMeasurementTagsArray Gets the array of tags for the measurements passed in filter */
/*The filter contains a list of measurement names*/
/*then with these measurement names a list of tags is obtained*/
func GetIfxMeasurementTagsArray(ctx *Context) {
	filter := ctx.Params(":filter")
	log.Debugf("Getting Influx Measurement Tags with filter: '%s'.", filter)
	tagsarray, err := agent.MainConfig.Database.GetIfxMeasurementTagsArray(filter)
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error getting tags: %+s", err)
		return
	}
	log.Debugf("Getting Influx Measurement Tags with filter: '%s'. Returning size: %d", filter, len(tagsarray))
	ctx.JSON(200, &tagsarray)
}

//GetIfxMeasurementAffectOnDel --pending--
func GetIfxMeasurementAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	nid, err := strconv.ParseInt(id, 10, 64)
	obarray, err := agent.MainConfig.Database.GetIfxMeasurementCfgAffectOnDel(nid)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
