package webui

import (
	"strconv"

	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"gopkg.in/macaron.v1"
)

// NewAPICfgDeviceStat get DeviceStat API
func NewAPICfgDeviceStat(m *macaron.Macaron) error {

	bind := binding.Bind

	// Data sources
	m.Group("/api/cfg/devicestat", func() {
		m.Get("/", reqSignedIn, GetDeviceStat)
		m.Post("/", reqSignedIn, bind(config.DeviceStatCfg{}), AddDeviceStat)
		m.Put("/:id", reqSignedIn, bind(config.DeviceStatCfg{}), UpdateDeviceStat)
		m.Delete("/:id", reqSignedIn, DeleteDeviceStat)
		m.Get("/:id", reqSignedIn, GetDeviceStatCfgByID)
		m.Get("/checkondel/:id", reqSignedIn, GetDeviceStatAffectOnDel)
	})

	return nil
}

// GetDeviceStat Return snmpdevice list to frontend
func GetDeviceStat(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetDeviceStatCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

// AddDeviceStat Insert new snmpdevice to de internal BBDD --pending--
func AddDeviceStat(ctx *Context, dev config.DeviceStatCfg) {
	log.Printf("ADDING DEVICE %+v", dev)
	affected, err := agent.MainConfig.Database.AddDeviceStatCfg(dev)
	if err != nil {
		log.Warningf("Error on insert for device %s  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return data  or affected
		ctx.JSON(200, &dev)
	}
}

// UpdateDeviceStat --pending--
func UpdateDeviceStat(ctx *Context, dev config.DeviceStatCfg) {
	id := ctx.Params(":id")
	log.Debugf("Trying to update: %+v", dev)
	nid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Warningf("Error on get ID for UpdateDeviceStats device %d , error: %s", dev.ID, err)
		ctx.JSON(404, err.Error())
	}
	affected, err := agent.MainConfig.Database.UpdateDeviceStatCfg(nid, dev)
	if err != nil {
		log.Warningf("Error on update for device %d  , affected : %+v , error: %s", dev.ID, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		//TODO: review if needed return device data
		ctx.JSON(200, &dev)
	}
}

// DeleteDeviceStat delete device stats
func DeleteDeviceStat(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("Trying to delete: %+v", id)
	nid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Warningf("Error on get ID for UpdateDeviceStats device %d  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	}
	affected, err := agent.MainConfig.Database.DelDeviceStatCfg(nid)
	if err != nil {
		log.Warningf("Error on delete1 for device %s  , affected : %+v , error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}

//GetDeviceStatCfgByID --pending--
func GetDeviceStatCfgByID(ctx *Context) {
	id := ctx.Params(":id")
	dev, err := agent.MainConfig.Database.GetDeviceStatCfgByID(id)
	if err != nil {
		log.Warningf("Error on get Device  for device %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetDeviceStatAffectOnDel --pending--
func GetDeviceStatAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetDeviceStatCfgAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for SNMP metrics %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}
