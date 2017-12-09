package webui

import (
	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
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
	})

	return nil
}

// GetIfxServer Return snmpdevice list to frontend
func GetIfxServer(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetIfxServerCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddIfxServer Insert new snmpdevice to de internal BBDD --pending--
func AddIfxServer(ctx *Context, dev config.IfxServerCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddIfxServerCfg(dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateIfxServer --pending--
func UpdateIfxServer(ctx *Context, dev config.IfxServerCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateIfxServerCfg(id, dev)
	if err != nil {
		log.Warningf("Error on update for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteIfxServer delete from the catalog database
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

//GetIfxServerCfgByID --pending--
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

//GetIfxServerAffectOnDel --pending--
func GetIfxServerAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetIfxServerCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
