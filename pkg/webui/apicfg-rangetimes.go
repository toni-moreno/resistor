package webui

import (
	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgRangeTime create a Range Time  management API
func NewAPICfgRangeTime(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/rangetimes", func() {
		m.Get("/", reqSignedIn, GetRangeTime)
		m.Post("/", reqSignedIn, bind(config.RangeTimeCfg{}), AddRangeTime)
		m.Put("/:id", reqSignedIn, bind(config.RangeTimeCfg{}), UpdateRangeTime)
		m.Delete("/:id", reqSignedIn, DeleteRangeTime)
		m.Get("/:id", reqSignedIn, GetRangeTimeCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetRangeTimeAffectOnDel)
	})

	return nil
}

// GetRangeTime Return snmpdevice list to frontend
func GetRangeTime(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetRangeTimeCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddRangeTime Insert new snmpdevice to de internal BBDD --pending--
func AddRangeTime(ctx *Context, dev config.RangeTimeCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddRangeTimeCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateRangeTime --pending--
func UpdateRangeTime(ctx *Context, dev config.RangeTimeCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	affected, err := agent.MainConfig.Database.UpdateRangeTimeCfg(id, &dev)
	if err != nil {
		log.Warningf("Error on update for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

//DeleteRangeTime delete range time
func DeleteRangeTime(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelRangeTimeCfg(id)
	if err != nil {
		log.Warningf("Error deleting range time %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetRangeTimeCfgByID Gets RangeTimeCfg By ID and returns it on context
func GetRangeTimeCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetRangeTimeCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetRangeTimeAffectOnDel Check if there are any AlertIDCfg affected when deleting RangeTimeCfg
func GetRangeTimeAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetRangeTimeCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error getting object array for range time %s , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
