package webui

import (
	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgOutHTTP create new API
func NewAPICfgOutHTTP(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/outhttp", func() {
		m.Get("/", reqSignedIn, GetOutHTTP)
		m.Post("/", reqSignedIn, bind(config.OutHTTPCfg{}), AddOutHTTP)
		m.Put("/:id", reqSignedIn, bind(config.OutHTTPCfg{}), UpdateOutHTTP)
		m.Delete("/:id", reqSignedIn, DeleteOutHTTP)
		m.Get("/:id", reqSignedIn, GetOutHTTPCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetOutHTTPAffectOnDel)
	})

	return nil
}

// GetOutHTTP Return snmpdevice list to frontend
func GetOutHTTP(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetOutHTTPCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddOutHTTP Insert new snmpdevice to de internal BBDD --pending--
func AddOutHTTP(ctx *Context, dev config.OutHTTPCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddOutHTTPCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateOutHTTP --pending--
func UpdateOutHTTP(ctx *Context, dev config.OutHTTPCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateOutHTTPCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteOutHTTP removes and output backend
func DeleteOutHTTP(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelOutHTTPCfg(id)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetOutHTTPCfgByID --pending--
func GetOutHTTPCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetOutHTTPCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetOutHTTPAffectOnDel --pending--
func GetOutHTTPAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetOutHTTPCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
