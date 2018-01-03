package webui

import (
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

// GetAlertID Return snmpdevice list to frontend
func GetAlertID(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetAlertIDCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddAlertID Insert new snmpdevice to de internal BBDD --pending--
func AddAlertID(ctx *Context, dev config.AlertIDCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddAlertIDCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateAlertID --pending--
func UpdateAlertID(ctx *Context, dev config.AlertIDCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateAlertIDCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
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
	affected, err := agent.MainConfig.Database.DelAlertIDCfg(id)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
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
