package webui

import (
	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgTemplate set API for the template management
func NewAPICfgTemplate(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/template", func() {
		m.Get("/", reqSignedIn, GetTemplate)
		m.Post("/", reqSignedIn, bind(config.TemplateCfg{}), AddTemplate)
		m.Put("/:id", reqSignedIn, bind(config.TemplateCfg{}), UpdateTemplate)
		m.Delete("/:id", reqSignedIn, DeleteTemplate)
		m.Get("/:id", reqSignedIn, GetTemplateCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetTemplateAffectOnDel)
	})

	return nil
}

// GetTemplate Return snmpdevice list to frontend
func GetTemplate(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetTemplateCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddTemplate Insert new snmpdevice to de internal BBDD --pending--
func AddTemplate(ctx *Context, dev config.TemplateCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddTemplateCfg(dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateTemplate --pending--
func UpdateTemplate(ctx *Context, dev config.TemplateCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateTemplateCfg(id, dev)
	if err != nil {
		log.Warningf("Error on update for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteTemplate delete template from database
func DeleteTemplate(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelTemplateCfg(id)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetTemplateCfgByID --pending--
func GetTemplateCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetTemplateCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetTemplateAffectOnDel --pending--
func GetTemplateAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetTemplateCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
