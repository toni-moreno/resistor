package webui

import (
	"fmt"
	"time"

	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"github.com/toni-moreno/resistor/pkg/kapa"
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
		m.Post("/deploy", reqSignedIn, bind(config.TemplateCfg{}), DeployTemplate)
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
	kapaserversarray, err := kapa.GetKapaServers("")
	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_ = kapa.GetKapaTemplates(devcfgarray, kapaserversarray)
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting templates %+v", &devcfgarray)
}

// GetTemplates Gets templates array
func GetTemplates(templateid string) ([]*config.TemplateCfg, error) {
	filter := ""
	if len(templateid) > 0 {
		filter = fmt.Sprintf("id = '%s'", templateid)
	}
	log.Debugf("Getting templates with filter: %s.", filter)
	devcfgarray, err := agent.MainConfig.Database.GetTemplateCfgArray(filter)
	if err != nil {
		log.Errorf("Error getting templates: %+s.", err)
	}
	return devcfgarray, err
}

// AddTemplate Inserts new template into the internal DB and into the kapacitor servers
func AddTemplate(ctx *Context, dev config.TemplateCfg) {
	dev.Modified = time.Now().UTC()
	sKapaSrvsNotOK := make([]string, 0)
	kapaserversarray, err := kapa.GetKapaServers("")
	errmsg := ""
	if err != nil {
		errmsg += fmt.Sprintf("Error getting kapacitor servers: %+s", err)
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_, _, sKapaSrvsNotOK, err = kapa.SetKapaTemplate(dev, kapaserversarray)
	}
	if err != nil && len(sKapaSrvsNotOK) > 0 {
		errmsg += " " + fmt.Sprintf("Error deploying template %s on kapacitor servers: %+v. Not deployed for kapacitor servers: %+v. Error: %s", dev.ID, dev.ServersWOLastDeployment, sKapaSrvsNotOK, err)
		log.Warningf("Error deploying template %s on kapacitor servers: %+v. Not deployed for kapacitor servers: %+v. Error: %s", dev.ID, dev.ServersWOLastDeployment, sKapaSrvsNotOK, err)
	} else if err != nil {
		errmsg += " " + fmt.Sprintf("Error deploying template %s. Error: %s", dev.ID, err)
		log.Warningf("Error deploying template %s. Error: %s", dev.ID, err)
	} else if len(sKapaSrvsNotOK) > 0 {
		errmsg += " " + fmt.Sprintf("Error deploying template %s. Not deployed for kapacitor servers: %+v.", dev.ID, sKapaSrvsNotOK)
		log.Warningf("Error deploying template %s. Not deployed for kapacitor servers: %+v.", dev.ID, sKapaSrvsNotOK)
	}
	log.Printf("ADDING template %+v", dev)
	affected, err := agent.MainConfig.Database.AddTemplateCfg(&dev)
	if err != nil || len(errmsg) > 0 {
		if err != nil {
			errmsg += " " + fmt.Sprintf("Error on insert for template %s  , affected : %+v , error: %s", dev.ID, affected, err)
			log.Warningf("Error on insert for template %s  , affected : %+v , error: %s", dev.ID, affected, err)
		}
		err = fmt.Errorf(errmsg)
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
	kapaserversarray, err := kapa.GetKapaServers("")
	errmsg := ""
	if err != nil {
		errmsg += fmt.Sprintf("Error getting kapacitor servers: %+s", err)
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_, _, sKapaSrvsNotOK, err = kapa.SetKapaTemplate(dev, kapaserversarray)
	}
	if err != nil && len(sKapaSrvsNotOK) > 0 {
		errmsg += " " + fmt.Sprintf("Error deploying template %s on kapacitor servers: %+v. Not deployed for kapacitor servers: %+v. Error: %s", dev.ID, dev.ServersWOLastDeployment, sKapaSrvsNotOK, err)
		log.Warningf("Error deploying template %s on kapacitor servers: %+v. Not deployed for kapacitor servers: %+v. Error: %s", dev.ID, dev.ServersWOLastDeployment, sKapaSrvsNotOK, err)
	} else if err != nil {
		errmsg += " " + fmt.Sprintf("Error deploying template %s. Error: %s", dev.ID, err)
		log.Warningf("Error deploying template %s. Error: %s", dev.ID, err)
	} else if len(sKapaSrvsNotOK) > 0 {
		errmsg += " " + fmt.Sprintf("Error deploying template %s. Not deployed for kapacitor servers: %+v.", dev.ID, sKapaSrvsNotOK)
		log.Warningf("Error deploying template %s. Not deployed for kapacitor servers: %+v.", dev.ID, sKapaSrvsNotOK)
	}
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateTemplateCfg(id, &dev)
	if err != nil || len(errmsg) > 0 {
		if err != nil {
			errmsg += " " + fmt.Sprintf("Error on update for template %s  , affected : %+v , error: %s", dev.ID, affected, err)
			log.Warningf("Error on update for template %s  , affected : %+v , error: %s", dev.ID, affected, err)
		}
		err = fmt.Errorf(errmsg)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return template data
		ctx.JSON(200, &dev)
	}
}

// DeployTemplate Deploys template into the kapacitor servers and returns the result in context
func DeployTemplate(ctx *Context, dev config.TemplateCfg) {
	sKapaSrvsNotOK, err := kapa.DeployKapaTemplate(dev)
	if err != nil {
		ctx.JSON(404, fmt.Sprintf("Error getting kapacitor servers from array: %+v. Error: %+s.", dev.ServersWOLastDeployment, err))
	} else if len(sKapaSrvsNotOK) > 0 {
		ctx.JSON(404, fmt.Sprintf("Error deploying template %s on kapacitor servers: %+v. Not updated for kapacitor servers: %+v.", dev.ID, dev.ServersWOLastDeployment, sKapaSrvsNotOK))
	} else {
		ctx.JSON(200, fmt.Sprintf("Template %s succesfully deployed on kapacitor servers: %+v.", dev.ID, dev.ServersWOLastDeployment))
	}
}

//DeleteTemplate Deletes template from resistor database and from kapacitor servers.
//First of all, a checking is done to ensure this template is not used by any resistor alert
func DeleteTemplate(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete template with id: %s.", id)
	tpl, err := agent.MainConfig.Database.GetTemplateCfgByID(id)
	if err != nil {
		log.Warningf("Error getting template %s. Error: %s", id, err)
		ctx.JSON(404, err.Error())
		return
	}
	//ensure this template is not used by any resistor alert
	idalertsarray, err := GetAlertIDCfgByTemplate(tpl.TriggerType, tpl.CritDirection, tpl.TrendType, tpl.TrendSign, tpl.FieldType, tpl.StatFunc)
	if err != nil {
		log.Warningf("Error getting alerts related to this template %s. Error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		if len(idalertsarray) > 0 {
			log.Warningf("This template, %s, can't be deleted because it's related with these alerts: %+v.", id, idalertsarray)
			ctx.JSON(404, fmt.Sprintf("This template, %s, can't be deleted because it's related with these alerts: %+v", id, idalertsarray))
		} else {
			//Get all kapacitor servers
			sKapaSrvsNotOK := make([]string, 0)
			kapaserversarray, err := kapa.GetKapaServers("")
			if err != nil {
				log.Warningf("Error getting kapacitor servers: %+s", err)
				ctx.JSON(404, err.Error())
			} else {
				//Try to delete template from all kapacitor servers
				_, _, sKapaSrvsNotOK = kapa.DeleteKapaTemplate(id, kapaserversarray)
				if len(sKapaSrvsNotOK) == 0 {
					//Try to delete template from resistor database
					affected, err := agent.MainConfig.Database.DelTemplateCfg(id)
					if err != nil {
						log.Warningf("Error on deleting for template %s  , affected : %+v , error: %s", id, affected, err)
						ctx.JSON(404, err.Error())
					} else {
						ctx.JSON(200, "deleted")
					}
				} else {
					log.Warningf("Error on deleting for template %s. Not deleted for kapacitor servers: %+v.", id, sKapaSrvsNotOK)
					ctx.JSON(404, fmt.Sprintf("Error on deleting for template %s. Not deleted for kapacitor servers: %+v", id, sKapaSrvsNotOK))
				}
			}
		}
	}
}

//GetTemplateCfgByID --pending--
func GetTemplateCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetTemplateCfgByID(id)
	if err != nil {
		log.Warningf("Error getting template with id %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		kapaserversarray, err := kapa.GetKapaServers("")
		if err != nil {
			log.Warningf("Error getting kapacitor servers: %+s", err)
		} else {
			_, _, sKapaSrvsNotOK := kapa.GetKapaTemplate(&dev, kapaserversarray)
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
		log.Warningf("Error on get object array for Templates %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
