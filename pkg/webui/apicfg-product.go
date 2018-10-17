package webui

import (
	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgProduct API for Product Catalog Management
func NewAPICfgProduct(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/product", func() {
		m.Get("/", reqSignedIn, GetProduct)
		m.Post("/", reqSignedIn, bind(config.ProductCfg{}), AddProduct)
		m.Put("/:id", reqSignedIn, bind(config.ProductCfg{}), UpdateProduct)
		m.Delete("/:id", reqSignedIn, DeleteProduct)
		m.Get("/:id", reqSignedIn, GetProductCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetProductAffectOnDel)
	})

	return nil
}

// GetProduct Return Product list to frontend
func GetProduct(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetProductCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Product :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting Product %+v", &devcfgarray)
}

// AddProduct Insert new Product to the internal BBDD --pending--
func AddProduct(ctx *Context, dev config.ProductCfg) {
	log.Printf("ADDING Product %+v", dev)
	affected, err := agent.MainConfig.Database.AddProductCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for Product %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateProduct --pending--
func UpdateProduct(ctx *Context, dev config.ProductCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateProductCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for Product %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteProduct delete from the catalog database
func DeleteProduct(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelProductCfg(id)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetProductCfgByID --pending--
func GetProductCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetProductCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetProductAffectOnDel --pending--
func GetProductAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetProductCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for Products %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
