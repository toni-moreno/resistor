package webui

import (
	"time"

	kapacitorClient "github.com/influxdata/kapacitor/client/v1"
	"github.com/toni-moreno/resistor/pkg/config"
	"github.com/toni-moreno/resistor/pkg/kapa"
	"gopkg.in/macaron.v1"
)

//KapaTaskRt Structure with Kapacitor server info and Task info
type KapaTaskRt struct {
	ServerID       string                         `json:"ServerID"`
	URL            string                         `json:"URL"`
	Description    string                         `json:"Description"`
	ID             string                         `json:"ID"`
	Type           kapacitorClient.TaskType       `json:"Type"`
	DBRPs          []kapacitorClient.DBRP         `json:"DBRPs"`
	TICKscript     string                         `json:"script"`
	Vars           kapacitorClient.Vars           `json:"vars"`
	Dot            string                         `json:"Dot"`
	Status         kapacitorClient.TaskStatus     `json:"Status"`
	Executing      bool                           `json:"Executing"`
	Error          string                         `json:"Error"`
	ExecutionStats kapacitorClient.ExecutionStats `json:"stats"`
	Created        time.Time                      `json:"Created"`
	Modified       time.Time                      `json:"Modified"`
	LastEnabled    time.Time                      `json:"LastEnabled,omitempty"`
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
	kapaserversarray, err := kapa.GetKapaServers("")
	if err != nil {
		log.Warningf("Error getting kapacitor servers: %+s", err)
		ctx.JSON(404, err.Error())
	} else {
		for _, kapasrv := range kapaserversarray {
			kapaGoClient, _, _, err := kapa.GetKapaClient(*kapasrv)
			if err != nil {
				log.Warningf("Error getting kapacitor Go client for kapacitor server: %s. Error: %s", kapasrv.ID, err)
				ctx.JSON(404, err.Error())
			} else {
				kapaTasksArray, err := kapa.ListKapaTasks(kapaGoClient)
				if err != nil {
					log.Warningf("Error getting kapacitor tasks from kapacitor server: %s. Error: %s", kapasrv.ID, err)
					ctx.JSON(404, err.Error())
				} else {
					for _, kapatask := range kapaTasksArray {
						kapaTaskRt = makeKapaTaskRt(kapasrv, kapatask)
						kapaTaskRtArray = append(kapaTaskRtArray, kapaTaskRt)
					}
				}
			}
		}
	}
	log.Debugf("Got tasks list with %d tasks from kapacitor servers %+v", len(kapaTaskRtArray), &kapaTaskRtArray)
	ctx.JSON(200, &kapaTaskRtArray)
}

func makeKapaTaskRt(kapasrv *config.KapacitorCfg, kapatask kapacitorClient.Task) KapaTaskRt {
	var kapaTaskRt KapaTaskRt
	kapaTaskRt.ServerID = kapasrv.ID
	kapaTaskRt.URL = kapasrv.URL
	kapaTaskRt.Description = kapasrv.Description
	kapaTaskRt.ID = kapatask.ID
	kapaTaskRt.Type = kapatask.Type
	kapaTaskRt.DBRPs = kapatask.DBRPs
	kapaTaskRt.TICKscript = kapatask.TICKscript
	kapaTaskRt.Vars = kapatask.Vars
	kapaTaskRt.Dot = kapatask.Dot
	kapaTaskRt.Status = kapatask.Status
	kapaTaskRt.Executing = kapatask.Executing
	kapaTaskRt.Error = kapatask.Error
	kapaTaskRt.ExecutionStats = kapatask.ExecutionStats
	kapaTaskRt.Created = kapatask.Created
	kapaTaskRt.Modified = kapatask.Modified
	kapaTaskRt.LastEnabled = kapatask.LastEnabled
	return kapaTaskRt
}
