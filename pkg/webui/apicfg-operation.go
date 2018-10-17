package webui

import (
	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgOperation Operation instructions
func NewAPICfgOperation(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/operation", func() {
		m.Get("/", reqSignedIn, GetOperation)
		m.Post("/", reqSignedIn, bind(config.OperationCfg{}), AddOperation)
		m.Put("/:id", reqSignedIn, bind(config.OperationCfg{}), UpdateOperation)
		m.Delete("/:id", reqSignedIn, DeleteOperation)
		m.Get("/:id", reqSignedIn, GetOperationCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetOperationAffectOnDel)
	})

	return nil
}

// GetOperation Return operation list to frontend
func GetOperation(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetOperationCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Operation :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting Operation %+v", &devcfgarray)
}

// AddOperation Insert new Operation to the internal BBDD
func AddOperation(ctx *Context, dev config.OperationCfg) {
	log.Printf("ADDING Operation %+v", dev)
	affected, err := agent.MainConfig.Database.AddOperationCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for Operation %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateOperation Updates Operation
func UpdateOperation(ctx *Context, dev config.OperationCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateOperationCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for Operation %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteOperation deletes Operation
func DeleteOperation(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelOperationCfg(id)
	if err != nil {
		log.Warningf("Error on delete for Operation %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetOperationCfgByID Gets OperationCfg By ID
func GetOperationCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetOperationCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Operation with id %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetOperationAffectOnDel Gets array of alerts affected On Deletion
func GetOperationAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetOperationCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for Operations %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
