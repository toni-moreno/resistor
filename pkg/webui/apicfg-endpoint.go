package webui

import (
	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgEndpoint create new API
func NewAPICfgEndpoint(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/endpoint", func() {
		m.Get("/", reqSignedIn, GetEndpoint)
		m.Post("/", reqSignedIn, bind(config.EndpointCfg{}), AddEndpoint)
		m.Put("/:id", reqSignedIn, bind(config.EndpointCfg{}), UpdateEndpoint)
		m.Delete("/:id", reqSignedIn, DeleteEndpoint)
		m.Get("/:id", reqSignedIn, GetEndpointCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetEndpointAffectOnDel)
	})

	return nil
}

// GetEndpoint Return snmpdevice list to frontend
func GetEndpoint(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetEndpointCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddEndpoint Insert new snmpdevice to de internal BBDD --pending--
func AddEndpoint(ctx *Context, dev config.EndpointCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddEndpointCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateEndpoint --pending--
func UpdateEndpoint(ctx *Context, dev config.EndpointCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateEndpointCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteEndpoint removes and output backend
func DeleteEndpoint(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelEndpointCfg(id)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetEndpointCfgByID --pending--
func GetEndpointCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetEndpointCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetEndpointAffectOnDel --pending--
func GetEndpointAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetEndpointCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
