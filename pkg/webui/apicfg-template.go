package webui

import (
	"time"

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

// GetTemplate Return templates list to frontend
func GetTemplate(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetTemplateCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get templates: %+s", err)
		return
	}
	kapaserversarray, err := GetKapaServers("")
	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_ = GetKapaTemplates(devcfgarray, kapaserversarray)
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting templates %+v", &devcfgarray)
}

// AddTemplate Inserts new template into the internal DB and into the kapacitor servers
func AddTemplate(ctx *Context, dev config.TemplateCfg) {
	dev.Modified = time.Now().UTC()
	sKapaSrvsNotOK := make([]string, 0)
	kapaserversarray, err := GetKapaServers("")
	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_, _, sKapaSrvsNotOK = SetKapaTemplate(dev, kapaserversarray)
	}
	if len(sKapaSrvsNotOK) > 0 {
		log.Warningf("Error on inserting for template %s. Not inserted for kapacitor servers: %+v.", dev.ID, sKapaSrvsNotOK)
	}
	log.Printf("ADDING template %+v", dev)
	affected, err := agent.MainConfig.Database.AddTemplateCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for template %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateTemplate Updates template into the internal DB and into the kapacitor servers
func UpdateTemplate(ctx *Context, dev config.TemplateCfg) {
	dev.Modified = time.Now().UTC()
	sKapaSrvsNotOK := make([]string, 0)
	kapaserversarray, err := GetKapaServers("")
	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_, _, sKapaSrvsNotOK = SetKapaTemplate(dev, kapaserversarray)
	}
	if len(sKapaSrvsNotOK) > 0 {
		log.Warningf("Error on updating for template %s. Not updated for kapacitor servers: %+v.", dev.ID, sKapaSrvsNotOK)
	}
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateTemplateCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for template %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteTemplate delete template from database
func DeleteTemplate(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete template: %+v", id)
	sKapaSrvsNotOK := make([]string, 0)
	kapaserversarray, err := GetKapaServers("")
	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_, _, sKapaSrvsNotOK = DeleteKapaTemplate(id, kapaserversarray)
	}
	if len(sKapaSrvsNotOK) == 0 {
		affected, err := agent.MainConfig.Database.DelTemplateCfg(id)
		if err != nil {
			log.Warningf("Error on deleting for template %s  , affected : %+v , error: %s", id, affected, err)
			ctx.JSON(404, err.Error())
		} else {
			ctx.JSON(200, "deleted")
		}
	} else {
		log.Warningf("Error on deleting for template %s. Not deleted for kapacitor servers: %+v.", id, sKapaSrvsNotOK)
		ctx.JSON(404, "Error on deleting for template. Not deleted for all kapacitor servers.")
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
		kapaserversarray, err := GetKapaServers("")
		if err != nil {
			log.Warningf("Error getting kapacitor servers: %+s", err)
		} else {
			_, _, sKapaSrvsNotOK := GetKapaTemplate(&dev, kapaserversarray)
			dev.ServersWOLastDeployment = sKapaSrvsNotOK
		}
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
