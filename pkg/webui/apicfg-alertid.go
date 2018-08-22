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

// NewAPICfgAlertID config API for alerts
func NewAPICfgAlertID(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/alertid", func() {
		m.Get("/", reqSignedIn, GetAlertID)
		m.Post("/", reqSignedIn, bind(config.AlertIDCfg{}), AddAlertID)
		m.Put("/:id", reqSignedIn, bind(config.AlertIDCfg{}), UpdateAlertID)
		m.Post("/deploy", reqSignedIn, bind(config.AlertIDCfg{}), DeployAlertID)
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
	_ = kapa.GetKapaTasks(devcfgarray)
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddAlertID Inserts new alert into the internal DB and into the kapacitor servers
func AddAlertID(ctx *Context, dev config.AlertIDCfg) {
	dev.Modified = time.Now().UTC()
	kapa.DeployKapaTask(dev)
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
	kapa.DeployKapaTask(dev)
	id := ctx.Params(":id") //oldID from form
	log.Debugf("Trying to update alert with id: %s and info: %+v", id, dev)
	affected, err := agent.MainConfig.Database.UpdateAlertIDCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for alert %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		if id != dev.ID {
			//If the name of the alert has been changed,
			//a new task has been created on kapacitor servers with the new name,
			//the kapacitor task with the old name must be deleted.
			_, _, sKapaSrvsNotOK := DeleteKapaTask(id)
			if len(sKapaSrvsNotOK) > 0 {
				log.Warningf("Error deleting task %s from kapacitor servers: %s", id, sKapaSrvsNotOK)
				ctx.JSON(404, err.Error())
			}
		}
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

// DeployAlertID Deploys the task related to this alert into the kapacitor server and returns the result in context
func DeployAlertID(ctx *Context, dev config.AlertIDCfg) {
	if len(dev.ServersWOLastDeployment) > 0 {
		sKapaSrvsNotOK, err := kapa.DeployKapaTask(dev)
		if err != nil {
			ctx.JSON(404, fmt.Sprintf("Error getting kapacitor servers: %+s", err))
		} else if len(sKapaSrvsNotOK) > 0 {
			ctx.JSON(404, fmt.Sprintf("Error deploying task %s. Not deployed for kapacitor server %s.", dev.ID, dev.KapacitorID))
		} else {
			ctx.JSON(200, fmt.Sprintf("Task %s deployed for kapacitor server %s.", dev.ID, dev.KapacitorID))
		}
	} else {
		log.Debugf("Task %s is deployed with the last version on the kapacitor server: %s.", dev.ID, dev.KapacitorID)
		ctx.JSON(200, fmt.Sprintf("Task %s is deployed with the last version on the kapacitor server: %s.", dev.ID, dev.KapacitorID))
	}
}

//DeleteKapaTask Deletes task from kapacitor servers
func DeleteKapaTask(id string) (int, int, []string) {
	kapaserversarray, err := kapa.GetKapaServers("")
	iNumKapaServers := len(kapaserversarray)
	iNumDeleted := 0
	sKapaSrvsNotOK := make([]string, 0)
	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
	} else {
		iNumKapaServers, iNumDeleted, sKapaSrvsNotOK = kapa.DeleteKapaTask(id, kapaserversarray)
	}
	return iNumKapaServers, iNumDeleted, sKapaSrvsNotOK
}

//DeleteAlertID removes alert from resistor database
func DeleteAlertID(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	_, _, sKapaSrvsNotOK := DeleteKapaTask(id)
	if len(sKapaSrvsNotOK) == 0 {
		affected, err := agent.MainConfig.Database.DelAlertIDCfg(id)
		if err != nil {
			log.Warningf("Error deleting alert %s, affected: %+v, error: %s", id, affected, err)
			ctx.JSON(404, err.Error())
		} else {
			ctx.JSON(200, "deleted")
		}
	} else {
		log.Warningf("Error deleting alert %s. It can't be deleted for kapacitor servers: %+v", id, sKapaSrvsNotOK)
		ctx.JSON(404, fmt.Sprintf("Error deleting alert %s. It can't be deleted for kapacitor servers: %+v", id, sKapaSrvsNotOK))
	}
}

//GetAlertIDCfgByID Gets AlertIDCfg By ID from resistor database
//and checks if the corresponding kapacitor task is deployed on all kapacitor servers.
//Returns the information of the process with a JSON in context
func GetAlertIDCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetAlertIDCfgByID(id)
	if err != nil {
		log.Warningf("Error getting alert with id: %s. Error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		_, _, sKapaSrvsNotOK := kapa.GetKapaTask(&dev)
		dev.ServersWOLastDeployment = sKapaSrvsNotOK
		ctx.JSON(200, &dev)
	}
}

//GetAlertIDCfgByTemplate Gets an array of strings with the IDs of the Alerts where this template is used.
//The input parameters are the 5 fields needed to define a template.
func GetAlertIDCfgByTemplate(sTriggerType string, sCritDirection string, sThresholdType string, sTrendSign string, sStatFunc string) ([]string, error) {
	filter := fmt.Sprintf("trigertype = '%s'", sTriggerType)
	if sTriggerType != "DEADMAN" {
		if len(sCritDirection) > 0 {
			filter += fmt.Sprintf(" and critdirection = '%s'", sCritDirection)
		}
		if len(sThresholdType) > 0 {
			filter += fmt.Sprintf(" and thresholdtype = '%s'", sThresholdType)
		}
		if len(sStatFunc) > 0 {
			filter += fmt.Sprintf(" and statfunc = '%s'", sStatFunc)
		}
		if sTriggerType == "TREND" {
			if len(sTrendSign) > 0 {
				filter += fmt.Sprintf(" and trendsign = '%s'", sTrendSign)
			}
		}
	}
	log.Debugf("GetAlertIDCfgByTemplate. Getting alerts with filter: %s.", filter)
	devcfgarray, err := agent.MainConfig.Database.GetAlertIDCfgArray(filter)
	idarray := make([]string, 0)
	if err != nil {
		log.Errorf("GetAlertIDCfgByTemplate. Error getting alerts with filter: %s. Error: %+s.", filter, err)
	} else {
		idarray = getAlertCfgIDArray(devcfgarray)
		log.Debugf("GetAlertIDCfgByTemplate. Filter: %s. Alerts: %+v.", filter, idarray)
	}
	return idarray, err
}

func getAlertCfgIDArray(devcfgarray []*config.AlertIDCfg) []string {
	idarray := make([]string, 0)
	for i := 0; i < len(devcfgarray); i++ {
		cfg := devcfgarray[i]
		idarray = append(idarray, cfg.ID)
	}
	return idarray
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
