package webui

import (
	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgKapacitor Kapacitor ouput
func NewAPICfgKapacitor(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/kapacitor", func() {
		m.Get("/", reqSignedIn, GetKapacitor)
		m.Post("/", reqSignedIn, bind(config.KapacitorCfg{}), AddKapacitor)
		m.Put("/:id", reqSignedIn, bind(config.KapacitorCfg{}), UpdateKapacitor)
		m.Delete("/:id", reqSignedIn, DeleteKapacitor)
		m.Get("/:id", reqSignedIn, GetKapacitorCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetKapacitorAffectOnDel)
	})

	return nil
}

// GetKapacitor Return snmpdevice list to frontend
func GetKapacitor(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetKapacitorCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddKapacitor Insert new snmpdevice to de internal BBDD --pending--
func AddKapacitor(ctx *Context, dev config.KapacitorCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddKapacitorCfg(dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateKapacitor --pending--
func UpdateKapacitor(ctx *Context, dev config.KapacitorCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateKapacitorCfg(id, dev)
	if err != nil {
		log.Warningf("Error on update for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteKapacitor delete a backend config
func DeleteKapacitor(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelKapacitorCfg(id)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetKapacitorCfgByID --pending--
func GetKapacitorCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetKapacitorCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetKapacitorAffectOnDel --pending--
func GetKapacitorAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetKapacitorCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
