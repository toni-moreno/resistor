package webui

import (
	"strconv"

	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgIfxMeasurement API for IfxMeasurement Catalog Management
func NewAPICfgIfxMeasurement(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/ifxmeasurement", func() {
		m.Get("/", reqSignedIn, GetIfxMeasurement)
		m.Post("/", reqSignedIn, bind(config.IfxMeasurementCfg{}), AddIfxMeasurement)
		m.Put("/:id", reqSignedIn, bind(config.IfxMeasurementCfg{}), UpdateIfxMeasurement)
		m.Delete("/:id", reqSignedIn, DeleteIfxMeasurement)
		m.Get("/:id", reqSignedIn, GetIfxMeasurementCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetIfxMeasurementAffectOnDel)
	})

	return nil
}

// GetIfxMeasurement Return snmpdevice list to frontend
func GetIfxMeasurement(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetIfxMeasurementCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddIfxMeasurement Insert new snmpdevice to de internal BBDD --pending--
func AddIfxMeasurement(ctx *Context, dev config.IfxMeasurementCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddIfxMeasurementCfg(dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateIfxMeasurement --pending--
func UpdateIfxMeasurement(ctx *Context, dev config.IfxMeasurementCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	nid, err := strconv.ParseInt(id, 10, 64)
	affected, err := agent.MainConfig.Database.UpdateIfxMeasurementCfg(nid, dev)
	if err != nil {
		log.Warningf("Error on update for device %d/%s  , affected : %+v , error: %s", dev.ID, dev.Name, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteIfxMeasurement delete from the catalog database
func DeleteIfxMeasurement(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	nid, err := strconv.ParseInt(id, 10, 64)
	affected, err := agent.MainConfig.Database.DelIfxMeasurementCfg(nid)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetIfxMeasurementCfgByID --pending--
func GetIfxMeasurementCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	nid, err := strconv.ParseInt(id, 10, 64)
	dev, err := agent.MainConfig.Database.GetIfxMeasurementCfgByID(nid)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetIfxMeasurementAffectOnDel --pending--
func GetIfxMeasurementAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	nid, err := strconv.ParseInt(id, 10, 64)
	obarray, err := agent.MainConfig.Database.GetIfxMeasurementCfgAffectOnDel(nid)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
