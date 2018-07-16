package webui

import (
	"fmt"
	"strings"
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
	kapaserversarray, err := GetKapaServers("")
	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		_ = GetKapaTemplates(devcfgarray, kapaserversarray)
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

// DeployTemplate Deploys template into the kapacitor servers
func DeployTemplate(ctx *Context, dev config.TemplateCfg) {
	if len(dev.ServersWOLastDeployment) > 0 {
		dev.Modified = time.Now().UTC()
		sKapaSrvsNotOK := make([]string, 0)
		kapaserversarray, err := GetKapaServersFromArray(dev.ServersWOLastDeployment)
		if err != nil {
			log.Warningf("Error getting kapacitor servers from array: %+v. Error: %+s.", dev.ServersWOLastDeployment, err)
			ctx.JSON(404, fmt.Sprintf("Error getting kapacitor servers from array: %+v. Error: %+s.", dev.ServersWOLastDeployment, err))
		} else {
			_, _, sKapaSrvsNotOK = SetKapaTemplate(dev, kapaserversarray)
			if len(sKapaSrvsNotOK) > 0 {
				log.Warningf("Error deploying template %s on kapacitor servers: %+v. Not updated for kapacitor servers: %+v.", dev.ID, dev.ServersWOLastDeployment, sKapaSrvsNotOK)
				ctx.JSON(404, fmt.Sprintf("Error deploying template %s on kapacitor servers: %+v. Not updated for kapacitor servers: %+v.", dev.ID, dev.ServersWOLastDeployment, sKapaSrvsNotOK))
			} else {
				log.Infof("Template %s succesfully deployed on kapacitor servers: %+v.", dev.ID, dev.ServersWOLastDeployment)
				ctx.JSON(200, fmt.Sprintf("Template %s succesfully deployed on kapacitor servers: %+v.", dev.ID, dev.ServersWOLastDeployment))
			}
		}
	} else {
		log.Debugf("Template %s is deployed with the last version on all kapacitor servers.", dev.ID)
		ctx.JSON(200, fmt.Sprintf("Template %s is deployed with the last version on all kapacitor servers.", dev.ID))
	}
}

//DeleteTemplate Deletes template from resistor database and from kapacitor servers.
//First of all, a checking is done to ensure this template is not used by any resistor alert
func DeleteTemplate(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete template with id: %s.", id)
	//ensure this template is not used by any resistor alert
	sTriggerType, sCritDirection, sThresholdType, sTrendSign, sStatFunc := getTemplateIDParts(id)
	idalertsarray, err := GetAlertIDCfgByTemplate(sTriggerType, sCritDirection, sThresholdType, sTrendSign, sStatFunc)
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
			kapaserversarray, err := GetKapaServers("")
			if err != nil {
				log.Warningf("Error getting kapacitor servers: %+s", err)
				ctx.JSON(404, err.Error())
			} else {
				//Try to delete template from all kapacitor servers
				_, _, sKapaSrvsNotOK = DeleteKapaTemplate(id, kapaserversarray)
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

//GetResTemplateCfgByID Gets the TemplateCfg information stored on resistor database
//including the kapacitor servers without last deployment
func GetResTemplateCfgByID(id string) (config.TemplateCfg, error) {
	log.Debugf("GetResTemplateCfgByID. Trying to get template with id: %s.", id)
	dev, err := agent.MainConfig.Database.GetTemplateCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
	} else {
		kapaserversarray, err := GetKapaServers("")
		if err != nil {
			log.Warningf("Error getting kapacitor servers: %+s", err)
		} else {
			_, _, sKapaSrvsNotOK := GetKapaTemplate(&dev, kapaserversarray)
			dev.ServersWOLastDeployment = sKapaSrvsNotOK
			log.Debugf("GetResTemplateCfgByID. Template with id: %s has not the last version deployed on: %+v.", id, sKapaSrvsNotOK)
		}
	}
	return dev, err
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

// getTemplateIDParts Gets TemplateID parts from TemplateID
// example: from: "THRESHOLD_2EX_AC_TAP_FMOVAVG" --> result: "THRESHOLD", "AC", "absolute", "positive", "MOVAVG"
// TrigerType + CritDirection + ThresholdTypeTranslated + StatFunc
func getTemplateIDParts(sTemplateID string) (string, string, string, string, string) {
	sTriggerType, sCritDirection, sThresholdType, sTrendSign, sStatFunc := "DEADMAN", "", "", "", ""
	if sTemplateID != "DEADMAN" {
		partsarray := strings.Split(sTemplateID, "_")
		if len(partsarray) == 5 {
			sTriggerType = partsarray[0]
			sCritDirection = partsarray[2]
			sThresholdType, sTrendSign = getTrendDetails(partsarray[3][1:])
			sStatFunc = partsarray[4][1:]
		}
	}
	log.Debugf("getTemplateIDParts. %s, %s, %s, %s, %s, %s.", sTemplateID, sTriggerType, sCritDirection, sThresholdType, sTrendSign, sStatFunc)
	return sTriggerType, sCritDirection, sThresholdType, sTrendSign, sStatFunc
}

// getTrendDetails Gets trend details
// A to absolute
// R to relative
// P to positive
// N to negative
func getTrendDetails(sInput string) (string, string) {
	sAbsRel := "absolute"
	sTrendSign := ""
	if strings.Index(sInput, "R") > -1 {
		sAbsRel = "relative"
	}
	if strings.Index(sInput, "P") > -1 {
		sTrendSign = "positive"
	} else if strings.Index(sInput, "N") > -1 {
		sTrendSign = "negative"
	}
	return sAbsRel, sTrendSign
}
