package webui

import (
	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgIfxDB API for IfxDB Catalog Management
func NewAPICfgIfxDB(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/ifxdb", func() {
		m.Get("/", reqSignedIn, GetIfxDB)
		m.Post("/", reqSignedIn, bind(config.IfxDBCfg{}), AddIfxDB)
		m.Put("/:id", reqSignedIn, bind(config.IfxDBCfg{}), UpdateIfxDB)
		m.Delete("/:id", reqSignedIn, DeleteIfxDB)
		m.Get("/:id", reqSignedIn, GetIfxDBCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetIfxDBAffectOnDel)
	})

	return nil
}

// GetIfxDB Return snmpdevice list to frontend
func GetIfxDB(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetIfxDBCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddIfxDB Insert new snmpdevice to de internal BBDD --pending--
func AddIfxDB(ctx *Context, dev config.IfxDBCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddIfxDBCfg(dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateIfxDB --pending--
func UpdateIfxDB(ctx *Context, dev config.IfxDBCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateIfxDBCfg(id, dev)
	if err != nil {
		log.Warningf("Error on update for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteIfxDB delete from the catalog database
func DeleteIfxDB(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelIfxDBCfg(id)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetIfxDBCfgByID --pending--
func GetIfxDBCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetIfxDBCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetIfxDBAffectOnDel --pending--
func GetIfxDBAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetIfxDBCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
