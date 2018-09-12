package webui

import (
	"strconv"
	"strings"

	"github.com/toni-moreno/resistor/pkg/agent"
	"gopkg.in/macaron.v1"
)

// NewAPIRtAlertEvent creates an Alert Events management API
func NewAPIRtAlertEvent(m *macaron.Macaron) error {

	// Data sources
	m.Group("/api/rt/alertevent", func() {
		m.Get("/", reqSignedIn, GetAlertEvent)
		m.Get("/:id", reqSignedIn, GetAlertEventHistByID)
		m.Delete("/:id", reqSignedIn, DeleteAlertEvent)
		m.Get("/checkondel/:id", reqSignedIn, GetAlertEventAffectOnDel)
		m.Get("/withparams/:params", reqSignedIn, GetAlertEventWithParams)
	})

	return nil
}

// GetAlertEventWithParams Returns Alert Events list to frontend
func GetAlertEventWithParams(ctx *Context) {
	params := ctx.Params(":params")
	log.Debugf("GetAlertEventWithParams. params:%s", params)
	paramsarray := strings.Split(params, "&")
	var page int64
	var itemsPerPage int64
	var maxSize int64
	filter := ""
	sortColumn := ""
	sortDir := ""
	for _, paramkv := range paramsarray {
		paramkvarray := strings.Split(paramkv, "=")
		switch paramkvarray[0] {
		case "page":
			page, _ = strconv.ParseInt(paramkvarray[1], 10, 64)
		case "itemsPerPage":
			itemsPerPage, _ = strconv.ParseInt(paramkvarray[1], 10, 64)
		case "maxSize":
			maxSize, _ = strconv.ParseInt(paramkvarray[1], 10, 64)
		case "sortColumn":
			sortColumn = paramkvarray[1]
		case "sortDir":
			sortDir = paramkvarray[1]
		}
	}
	alevtarray, err := agent.MainConfig.Database.GetAlertEventHistArrayWithParams(filter, page, itemsPerPage, maxSize, sortColumn, sortDir)
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error getting AlertEvents:%+s", err)
		return
	}
	ctx.JSON(200, &alevtarray)
}

// GetAlertEvent Returns Alert Events list to frontend
func GetAlertEvent(ctx *Context) {
	alevtarray, err := agent.MainConfig.Database.GetAlertEventHistArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error getting AlertEvents:%+s", err)
		return
	}
	ctx.JSON(200, &alevtarray)
}

//GetAlertEventHistByID Returns Alert Event to frontend
func GetAlertEventHistByID(ctx *Context) {
	id, _ := strconv.ParseInt(ctx.Params(":id"), 10, 64)
	dev, err := agent.MainConfig.Database.GetAlertEventHistByID(id)
	if err != nil {
		log.Warningf("Error getting AlertEvent with id %d , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetAlertEventAffectOnDel Returns array of objects affected when deleting an alert event (empty array in this case)
func GetAlertEventAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetAlertEventAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for Alert Event %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}

//DeleteAlertEvent removes alert event from resistor database
func DeleteAlertEvent(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("DeleteAlertEvent. Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelAlertEventHist(id)
	if err != nil {
		log.Warningf("Error deleting alert event %s, affected: %+v, error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}
