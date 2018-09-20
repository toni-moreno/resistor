package webui

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	kapacitorClient "github.com/influxdata/kapacitor/client/v1"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"github.com/toni-moreno/resistor/pkg/kapa"
	"gopkg.in/macaron.v1"
)

//KapaTaskRt Structure with Kapacitor server info and Task info
type KapaTaskRt struct {
	ID             string                     `json:"ID"`
	ServerID       string                     `json:"ServerID"`
	URL            string                     `json:"URL,omitempty"`
	Type           kapacitorClient.TaskType   `json:"Type,omitempty"`
	DBRPs          string                     `json:"DBRPs,omitempty"`
	TICKscript     string                     `json:"script,omitempty"`
	Vars           string                     `json:"vars,omitempty"`
	Dot            string                     `json:"Dot,omitempty"`
	Status         kapacitorClient.TaskStatus `json:"Status,omitempty"`
	Executing      bool                       `json:"Executing,omitempty"`
	Error          string                     `json:"Error,omitempty"`
	NumErrors      int64                      `json:"NumErrors,omitempty"`
	ExecutionStats string                     `json:"stats,omitempty"`
	Created        time.Time                  `json:"Created,omitempty"`
	Modified       time.Time                  `json:"Modified,omitempty"`
	LastEnabled    time.Time                  `json:"LastEnabled,omitempty"`
	AlertModified  time.Time                  `json:"AlertModified,omitempty"`
}

// NewAPIRtKapacitor Kapacitor ouput
func NewAPIRtKapacitor(m *macaron.Macaron) error {

	//bind := binding.Bind

	// Data sources
	m.Group("/api/rt/kapacitor", func() {
		m.Get("/tasks/", reqSignedIn, GetKapacitorRtTasks)
		/*
			m.Get("/tasks/", reqSignedIn, GetKapacitorRtTasks)
			m.Get("/tasks/:id", reqSignedIn, GetKapacitorRtTasksByID)
			m.Post("/tasks/enable/:id", reqSignedIn, EnableKapacitorRtTasksByID)
			m.Post("/tasks/disable/:id", reqSignedIn, DisableKapacitorRtTasksByID)
		*/
	})

	return nil
}

// GetKapacitorRtTasks Return tasks list from kapacitor servers
func GetKapacitorRtTasks(ctx *Context) {
	var kapaTaskRtArray []KapaTaskRt
	var kapaTaskRt KapaTaskRt
	var kapaError = ""
	kapaTaskRtMap := make(map[string]KapaTaskRt)
	alertcfgarray, _ := agent.MainConfig.Database.GetAlertIDCfgArray("")
	if len(alertcfgarray) > 0 {
		kapaserversmap, err := agent.MainConfig.Database.GetKapacitorCfgMap("")
		if err != nil {
			kapaError = fmt.Sprintf("Error getting kapacitor servers: %+s", err)
			log.Warningf(kapaError)
		} else {
			for _, kapasrv := range kapaserversmap {
				kapaGoClient, _, _, err := kapa.GetKapaClient(*kapasrv)
				if err != nil {
					kapaError = fmt.Sprintf("Error getting kapacitor Go client for kapacitor server: %s. Error: %s", kapasrv.ID, err)
					log.Warningf(kapaError)
				} else {
					kapaTasksArray, err := kapa.ListKapaTasks(kapaGoClient)
					if err != nil {
						kapaError = fmt.Sprintf("Error getting kapacitor tasks from kapacitor server: %s. Error: %s", kapasrv.ID, err)
						log.Warningf(kapaError)
					} else {
						for _, kapatask := range kapaTasksArray {
							kapaTaskRt = makeKapaTaskRt(kapasrv, kapatask)
							kapaTaskRtMap[kapaTaskRt.ID] = kapaTaskRt
						}
					}
				}
			}
		}
		kapaTaskRtArray = makeKapaTaskRtArray(alertcfgarray, kapaserversmap, kapaTaskRtMap, kapaError)
		log.Debugf("Got tasks list with %d tasks from kapacitor servers %+v", len(kapaTaskRtArray), &kapaTaskRtArray)
	}
	ctx.JSON(200, &kapaTaskRtArray)
}

func makeKapaTaskRt(kapasrv *config.KapacitorCfg, kapatask kapacitorClient.Task) KapaTaskRt {
	var kapaTaskRt KapaTaskRt
	kapaTaskRt.ID = kapatask.ID
	kapaTaskRt.Type = kapatask.Type
	kapaTaskRt.ServerID = kapasrv.ID
	kapaTaskRt.URL = kapasrv.URL
	kapaTaskRt.DBRPs = kapatask.DBRPs[0].Database + "." + kapatask.DBRPs[0].RetentionPolicy
	kapaTaskRt.TICKscript = kapatask.TICKscript
	jsonArByt, err := json.Marshal(kapatask.Vars)
	if err != nil {
		log.Warningf("makeKapaTaskRt. Error Marshalling kapatask.Vars. Error: %s", err)
	}
	kapaTaskRt.Vars = string(jsonArByt)
	kapaTaskRt.Dot = kapatask.Dot
	kapaTaskRt.Status = kapatask.Status
	kapaTaskRt.Executing = kapatask.Executing
	kapaTaskRt.Error = kapatask.Error
	jsonArByt, err = json.Marshal(kapatask.ExecutionStats)
	if err != nil {
		log.Warningf("makeKapaTaskRt. Error Marshalling kapatask.ExecutionStats. Error: %s", err)
	}
	kapaTaskRt.ExecutionStats = string(jsonArByt)
	kapaTaskRt.Created = kapatask.Created
	kapaTaskRt.Modified = kapatask.Modified
	kapaTaskRt.LastEnabled = kapatask.LastEnabled
	var numErrors int64
	for _, nodestats := range kapatask.ExecutionStats.NodeStats {
		nodeErrors, err := strconv.ParseInt(fmt.Sprintf("%v", nodestats["errors"]), 10, 64)
		if err == nil {
			numErrors = numErrors + nodeErrors
		}
	}
	kapaTaskRt.NumErrors = numErrors
	return kapaTaskRt
}

func makeKapaTaskRtArray(alertcfgarray []*config.AlertIDCfg, kapaserversmap map[string]*config.KapacitorCfg, kapaTaskRtMap map[string]KapaTaskRt, kapaError string) []KapaTaskRt {
	var kapaTaskRtArray []KapaTaskRt
	for _, alertcfg := range alertcfgarray {
		kapaTaskRt, found := kapaTaskRtMap[alertcfg.ID]
		if !found {
			kapaTaskRt = newKapaTaskRt(alertcfg, kapaserversmap)
			if len(kapaError) > 0 {
				kapaTaskRt.Error = kapaError
			} else {
				kapaTaskRt.Error = "Error when deploying task on kapacitor server"
			}
		}
		if len(kapaError) > 0 {
			kapaTaskRt.Error = kapaError
		}
		kapaTaskRt.AlertModified = alertcfg.Modified
		kapaTaskRtArray = append(kapaTaskRtArray, kapaTaskRt)
	}
	return kapaTaskRtArray
}

func newKapaTaskRt(alertcfg *config.AlertIDCfg, kapaserversmap map[string]*config.KapacitorCfg) KapaTaskRt {
	var kapaTaskRt KapaTaskRt
	kapaTaskRt.ID = alertcfg.ID
	kapaTaskRt.ServerID = alertcfg.KapacitorID
	kapaserverCfg, found := kapaserversmap[alertcfg.KapacitorID]
	if found {
		kapaTaskRt.URL = kapaserverCfg.URL
	}
	kapaTaskRt.DBRPs = kapa.GetIfxDBNameByID(alertcfg.InfluxDB) + "." + alertcfg.InfluxRP
	return kapaTaskRt
}
