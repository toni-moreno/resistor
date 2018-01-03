package webui

import (
	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgProductGroup API for Product Catalog Management
func NewAPICfgProductGroup(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/productgroup", func() {
		m.Get("/", reqSignedIn, GetProductGroup)
		m.Post("/", reqSignedIn, bind(config.ProductGroupCfg{}), AddProductGroup)
		m.Put("/:id", reqSignedIn, bind(config.ProductGroupCfg{}), UpdateProductGroup)
		m.Delete("/:id", reqSignedIn, DeleteProductGroup)
		m.Get("/:id", reqSignedIn, GetProductGroupCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetProductGroupAffectOnDel)
	})

	return nil
}

// GetProductGroup Return snmpdevice list to frontend
func GetProductGroup(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetProductGroupCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddProductGroup Insert new snmpdevice to de internal BBDD --pending--
func AddProductGroup(ctx *Context, dev config.ProductGroupCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddProductGroupCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateProductGroup --pending--
func UpdateProductGroup(ctx *Context, dev config.ProductGroupCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateProductGroupCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteProductGroup delete from the catalog database
func DeleteProductGroup(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelProductGroupCfg(id)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetProductGroupCfgByID --pending--
func GetProductGroupCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetProductGroupCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetProductGroupAffectOnDel --pending--
func GetProductGroupAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetProductGroupCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
