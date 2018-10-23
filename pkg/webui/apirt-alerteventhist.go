package webui

import (
	"strconv"
	"strings"

	"github.com/toni-moreno/resistor/pkg/agent"
	"gopkg.in/macaron.v1"
)

// NewAPIRtAlertEventHist creates an Alert Event History management API
func NewAPIRtAlertEventHist(m *macaron.Macaron) error {

	// Data sources
	m.Group("/api/rt/alerteventhist", func() {
		m.Get("/", GetAlertEventHist)
		m.Get("/:id", GetAlertEventHistByID)
		m.Delete("/:id", reqSignedIn, DeleteAlertEventHist)
		m.Get("/checkondel/:id", reqSignedIn, GetAlertEventHistAffectOnDel)
		m.Get("/withparams/:params", GetAlertEventHistWithParams)
		m.Get("/groupbylevel/", GetAlertEventsHistByLevel)
	})

	return nil
}

// GetAlertEventHistWithParams Returns Alert Events list to frontend
func GetAlertEventHistWithParams(ctx *Context) {
	params := ctx.Params(":params")
	log.Debugf("GetAlertEventHistWithParams. params:%s", params)
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
		log.Errorf("Error getting AlertEventHist:%+s", err)
		return
	}
	ctx.JSON(200, &alevtarray)
}

// GetAlertEventHist Returns Alert Events list to frontend
func GetAlertEventHist(ctx *Context) {
	alevtarray, err := agent.MainConfig.Database.GetAlertEventHistArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error getting AlertEventHist:%+s", err)
		return
	}
	ctx.JSON(200, &alevtarray)
}

// GetAlertEventsHistByLevel Returns Alert Events History grouped by level to frontend
func GetAlertEventsHistByLevel(ctx *Context) {
	alevthistsummarray, err := agent.MainConfig.Database.GetAlertEventsHistByLevelArray("")
	if err != nil {
		ctx.JSON(404, err.Error())
		log.Errorf("Error getting AlertEventsHistByLevel:%+s", err)
		return
	}
	ctx.JSON(200, &alevthistsummarray)
}

//GetAlertEventHistByID Returns Alert Event to frontend
func GetAlertEventHistByID(ctx *Context) {
	id, _ := strconv.ParseInt(ctx.Params(":id"), 10, 64)
	dev, err := agent.MainConfig.Database.GetAlertEventHistByID(id)
	if err != nil {
		log.Warningf("Error getting AlertEventHist with id %d , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &dev)
	}
}

//GetAlertEventHistAffectOnDel Returns array of objects affected when deleting an alert event (empty array in this case)
func GetAlertEventHistAffectOnDel(ctx *Context) {
	id := ctx.Params(":id")
	obarray, err := agent.MainConfig.Database.GetAlertEventHistAffectOnDel(id)
	if err != nil {
		log.Warningf("Error on get object array for Alert Event %s  , error: %s", id, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, &obarray)
	}
}

//DeleteAlertEventHist removes alert event from resistor database
func DeleteAlertEventHist(ctx *Context) {
	id := ctx.Params(":id")
	log.Debugf("DeleteAlertEventHist. Trying to delete: %+v", id)
	affected, err := agent.MainConfig.Database.DelAlertEventHist(id)
	if err != nil {
		log.Warningf("Error deleting alert event %s, affected: %+v, error: %s", id, affected, err)
		ctx.JSON(404, err.Error())
	} else {
		ctx.JSON(200, "deleted")
	}
}
