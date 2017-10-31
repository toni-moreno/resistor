package webui

import (
	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgAlertId
func NewAPICfgAlertId(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/alertid", func() {
		m.Get("/", reqSignedIn, GetAlertID)
		m.Post("/", reqSignedIn, bind(config.AlertIdCfg{}), AddAlertId)
		m.Put("/:id", reqSignedIn, bind(config.AlertIdCfg{}), UpdateAlertId)
		m.Delete("/:id", reqSignedIn, DeleteAlertId)
		m.Get("/:id", reqSignedIn, GetAlertIdCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetAlertIDAffectOnDel)
	})

	return nil
}

// GetAlertID Return snmpdevice list to frontend
func GetAlertID(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetAlertIdCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddAlertId Insert new snmpdevice to de internal BBDD --pending--
func AddAlertId(ctx *Context, dev config.AlertIdCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddAlertIdCfg(dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateAlertId --pending--
func UpdateAlertId(ctx *Context, dev config.AlertIdCfg) {
	id := ctx.Params(":id")
	log.Debugf("Tying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateAlertIdCfg(id, dev)
	if err != nil {
		log.Warningf("Error on update for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteAlertId
func DeleteAlertId(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Tying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelAlertIdCfg(id)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//etAlertIdCfgByID --pending--
func GetAlertIdCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetAlertIdCfgByID(id)
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
	obarray, err := agent.MainConfig.Database.GetAlertIdCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
