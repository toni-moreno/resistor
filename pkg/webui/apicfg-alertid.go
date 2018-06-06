package webui

import (
	"time"

	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgAlertID config API for alerts
func NewAPICfgAlertID(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/alertid", func() {
		m.Get("/", reqSignedIn, GetAlertID)
		m.Post("/", reqSignedIn, bind(config.AlertIDCfg{}), AddAlertID)
		m.Put("/:id", reqSignedIn, bind(config.AlertIDCfg{}), UpdateAlertID)
		m.Delete("/:id", reqSignedIn, DeleteAlertID)
		m.Get("/:id", reqSignedIn, GetAlertIDCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetAlertIDAffectOnDel)
	})

	return nil
}

// GetAlertID Return alerts list to frontend
func GetAlertID(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetAlertIDCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	kapaserversarray, err := GetKapaServers("")
	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_ = GetKapaTasks(devcfgarray, kapaserversarray)
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddAlertID Inserts new alert into the internal DB and into the kapacitor servers
func AddAlertID(ctx *Context, dev config.AlertIDCfg) {
	dev.Modified = time.Now().UTC()
	sKapaSrvsNotOK := make([]string, 0)
	kapaserversarray, err := GetKapaServers(dev.KapacitorID)
	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_, _, sKapaSrvsNotOK = SetKapaTask(dev, kapaserversarray)
	}
	if len(sKapaSrvsNotOK) > 0 {
		log.Warningf("Error on inserting for alert %s. Not inserted for kapacitor server %s.", dev.ID, dev.KapacitorID)
	}
	log.Printf("ADDING alert %+v", dev)
	affected, err := agent.MainConfig.Database.AddAlertIDCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for alert %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateAlertID Updates alert into the internal DB and into the kapacitor servers
func UpdateAlertID(ctx *Context, dev config.AlertIDCfg) {
	dev.Modified = time.Now().UTC()
	sKapaSrvsNotOK := make([]string, 0)
	kapaserversarray, err := GetKapaServers(dev.KapacitorID)
	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_, _, sKapaSrvsNotOK = SetKapaTask(dev, kapaserversarray)
	}
	if len(sKapaSrvsNotOK) > 0 {
		log.Warningf("Error on updating for alert %s. Not updated for kapacitor server %s.", dev.ID, dev.KapacitorID)
	}
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateAlertIDCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for alert %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteAlertID removes alert from
func DeleteAlertID(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	sKapaSrvsNotOK := make([]string, 0)
	kapaserversarray, err := GetKapaServers("")
	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_, _, sKapaSrvsNotOK = DeleteKapaTask(id, kapaserversarray)
	}
	if len(sKapaSrvsNotOK) == 0 {
		affected, err := agent.MainConfig.Database.DelAlertIDCfg(id)
		if err != nil {
			log.Warningf("Error on deleting for alert %s  , affected : %+v , error: %s", id, affected, err)
			ctx.JSON(404, err.Error())
		} else {
			ctx.JSON(200, "deleted")
		}
	} else {
		log.Warningf("Error on deleting for alert %s", id)
		ctx.JSON(404, "Error on deleting for alert. Not deleted for all kapacitor servers.")
	}
}

//GetAlertIDCfgByID --pending--
func GetAlertIDCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetAlertIDCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		kapaserversarray, err := GetKapaServers("")
		if err != nil {
			log.Warningf("Error getting kapacitor servers: %+s", err)
		} else {
			_, _, sKapaSrvsNotOK := GetKapaTask(&dev, kapaserversarray)
			dev.ServersWOLastDeployment = sKapaSrvsNotOK
		}
		ctx.JSON(200, &dev)
	}
}

//GetAlertIDAffectOnDel --pending--
func GetAlertIDAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetAlertIDCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
