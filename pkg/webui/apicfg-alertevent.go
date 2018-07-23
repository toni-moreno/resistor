package webui

import (
	"strconv"

	"github.com/toni-moreno/resistor/pkg/agent"
	"gopkg.in/macaron.v1"
)

// NewAPICfgAlertEvent create a Range Time  management API
func NewAPICfgAlertEvent(m *macaron.Macaron) error {

	// Data sources
	m.Group("/api/cfg/alertevent", func() {
		m.Get("/", reqSignedIn, GetAlertEvent)
		m.Get("/:uid", reqSignedIn, GetAlertEventCfgByID)
	})

	return nil
}

// GetAlertEvent Return snmpdevice list to frontend
func GetAlertEvent(ctx *Context) {
	devcfgarray, err := agent.MainConfig.Database.GetAlertEventCfgArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error on get Devices :%+s", err)
		return
	}
	ctx.JSON(200, &devcfgarray)
	log.Debugf("Getting DEVICEs %+v", &devcfgarray)
}

//GetAlertEventCfgByID --pending--
func GetAlertEventCfgByID(ctx *Context) {
	uid, _ := strconv.ParseInt(ctx.Params(":uid"), 10, 64)
	dev, err := agent.MainConfig.Database.GetAlertEventCfgByID(uid)
	if err != nil {
		log.Warningf("Error on get Device for device %d , error: %s", uid, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}
