package webui

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-macaron/binding"
	"github.com/influxdata/kapacitor/alert"
	"github.com/influxdata/kapacitor/keyvalue"
	//"github.com/influxdata/kapacitor/services/smtp"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"github.com/toni-moreno/resistor/pkg/kapa"
	//"github.com/toni-moreno/resistor/pkg/data/alertfilter"
	res_smtp "github.com/toni-moreno/resistor/pkg/data/alertsvcs/smtp"

	"gopkg.in/macaron.v1"
)

// NewAPIRtKapFilter set the runtime Kapacitor filter  API
func NewAPIRtKapFilter(m *macaron.Macaron) error {

	bind := binding.Bind
	m.Group("/api/rt/kapfilter", func() {
		m.Post("/alert/:endpoint", reqAlertSignedIn, bind(alert.Data{}), RTAlertHandler)
		m.Post("/alert/", reqAlertSignedIn, bind(alert.Data{}), RTAlertHandler)
	})
	return nil
}

//RTAlertHandler xx
func RTAlertHandler(ctx *Context, al alert.Data) {
	/**/

	rb := ctx.Req.Body()
	s, _ := rb.String()
	log.Debugf("REQ: %s", s)
	log.Debugf("ALERT: %#+v", al)
	log.Debugf("ALERT Data: %#+v", al.Data)
	log.Debugf("ALERT Series: %+v", al.Data.Series)

	for _, serie := range al.Data.Series {
		log.Debugf("ALERT Serie: %+v", serie)
	}

	//Get TaskName
	sTaskName := al.ID
	idPartsArray := strings.Split(al.ID, "|")
	if len(idPartsArray) > 0 {
		sTaskName = idPartsArray[0]
	}
	//Get AlertCfg
	alertcfg, err := agent.MainConfig.Database.GetAlertIDCfgByID(sTaskName)
	if err != nil {
		log.Warningf("Error getting alert cfg with id: %s. Error: %s", al.ID, err)
	}
	sortedtagsarray := sortTagsMap(al.Data.Series[0].Tags)
	correlationID := "[" + alertcfg.ID + "]|" + strings.Join(sortedtagsarray, ",")
	//Save current alert event and move previous alert event
	alertevent := saveAlertEvent(correlationID, al, alertcfg, sortedtagsarray)
	//makeTaskAlertInfo
	taskAlertInfo, err := makeTaskAlertInfo(al, alertcfg, correlationID, alertevent.FirstEventTime, alertevent.EventTime)
	if err != nil {
		log.Warningf("RTAlertHandler. Error making taskAlertInfo. Error: %s", err)
	}

	//Send alert event to related endpoints
	for _, endpointid := range alertcfg.Endpoint {
		endpoint, err := agent.MainConfig.Database.GetEndpointCfgByID(endpointid)
		if err != nil {
			log.Warningf("Error getting endpoint for id %s. Error: %s.", endpointid, err)
		} else {
			log.Debugf("Got endpoint: %+v", endpoint)
			err = sendData(taskAlertInfo, al, alertcfg, sortedtagsarray, endpoint)
			if err != nil {
				log.Warningf("Error sending data to endpoint with id %s. Error: %s.", endpointid, err)
			}
		}
	}

	//alertfilter.ProcessAlert(al)

	ctx.JSON(200, "DONE")
}

func sortTagsMap(tagsmap map[string]string) []string {
	var sortedtagsarray []string
	var tagkeysarray []string
	for k := range tagsmap {
		tagkeysarray = append(tagkeysarray, k)
	}
	sort.Strings(tagkeysarray)
	for _, tagkey := range tagkeysarray {
		sortedtagsarray = append(sortedtagsarray, tagkey+":"+tagsmap[tagkey])
	}
	return sortedtagsarray
}

func makeAlertEvent(correlationID string, al alert.Data, alertcfg config.AlertIDCfg, sortedtagsarray []string, prevalevtarray []*config.AlertEvent) config.AlertEvent {
	log.Debugf("makeAlertEvent. Making alert event with this CorrelationID %s", correlationID)
	alertevent := config.AlertEvent{}
	alertevent.ID = 0
	alertevent.CorrelationID = correlationID
	alertevent.AlertID = alertcfg.ID
	alertevent.Message = al.Message
	alertevent.Details = al.Details
	if len(prevalevtarray) > 0 {
		if prevalevtarray[0].Level != "OK" {
			alertevent.FirstEventTime = prevalevtarray[0].FirstEventTime
		}
	} else {
		alertevent.FirstEventTime = al.Time
	}
	alertevent.EventTime = al.Time
	alertevent.Duration = al.Duration
	alertevent.Level = al.Level.String()
	alertevent.Field = alertcfg.Field
	alertevent.ProductID = alertcfg.ProductID
	alertevent.Tags = sortedtagsarray
	alertevent.ProductTagValue = alertcfg.ProductTag + ":" + al.Data.Series[0].Tags[alertcfg.ProductTag]
	columnsarray := al.Data.Series[0].Columns
	valuesarray := al.Data.Series[0].Values[0]
	for colidx, colvalue := range columnsarray {
		if colvalue == "value" {
			alertevent.Value = valuesarray[colidx].(float64)
		} else if colvalue == "mon_exc" {
			alertevent.MonExc = fmt.Sprintf("%v", valuesarray[colidx])
		}
	}
	return alertevent
}

func saveAlertEvent(correlationID string, al alert.Data, alertcfg config.AlertIDCfg, sortedtagsarray []string) config.AlertEvent {
	//Move previous alert events with this correlationid from alert_event to alert_event_hist
	filter := "correlationid = '" + correlationID + "'"
	prevalevtarray := MoveAlertEvents(filter)
	//Make current alert event
	alertevent := makeAlertEvent(correlationID, al, alertcfg, sortedtagsarray, prevalevtarray)
	//Insert current alert event into alert_event table
	err := addAlertEvent(alertevent)
	if err != nil {
		log.Warningf("saveAlertEvent. Error inserting current alert event with this correlationid %s into alert_event. Error: %s", correlationID, err)
	}
	return alertevent
}

//MoveAlertEvents Moves previous alert events with this filter from alert_event to alert_event_hist
func MoveAlertEvents(filter string) []*config.AlertEvent {
	log.Debugf("MoveAlertEvents. Moving previous alert event with this filter %s from alert_event to alert_event_hist.", filter)
	//Get previous alert events with this filter from alert_event
	prevalevtarray, err := getAlertEventArray(filter)
	if err != nil {
		log.Warningf("MoveAlertEvents. Error getting previous alert events with this filter %s from alert_event. Error: %s", filter, err)
	}
	if len(prevalevtarray) > 0 {
		log.Debugf("MoveAlertEvents. Got previous alert events with this filter %s from alert_event.", filter)
		for _, prevalevt := range prevalevtarray {
			//Insert previous alert event with this filter into alert_event_hist
			err = addAlertEventHist(config.AlertEventHist(*prevalevt))
			if err != nil {
				log.Warningf("MoveAlertEvents. Error inserting previous alert event with id %d into alert_event_hist. Error: %s", prevalevt.ID, err)
			} else {
				//Delete previous alert event with this filter from alert_event
				err = deleteAlertEvent(*prevalevt)
				if err != nil {
					log.Warningf("MoveAlertEvents. Error deleting previous alert event with id %d from alert_event. Error: %s", prevalevt.ID, err)
				}
			}
		}
	}
	return prevalevtarray
}

// getAlertEventArray Gets alert events with filter from alert_event
func getAlertEventArray(filter string) ([]*config.AlertEvent, error) {
	log.Debugf("Getting alert events with filter %s", filter)
	alevtarray, err := agent.MainConfig.Database.GetAlertEventArray(filter)
	if err != nil {
		log.Warningf("Error getting alert events with filter %s. Error: %s", filter, err)
	}
	return alevtarray, err
}

//addAlertEventHist Inserts alert event hist into alert_event_hist
func addAlertEventHist(dev config.AlertEventHist) error {
	log.Debugf("ADDING alert event hist %+v", dev)
	affected, err := agent.MainConfig.Database.AddAlertEventHist(&dev)
	if err != nil {
		log.Warningf("Error on insert for alert event hist %d , affected : %+v , error: %s", dev.ID, affected, err)
	}
	return err
}

//deleteAlertEvent Deletes previous alert event from alert_event
func deleteAlertEvent(alevt config.AlertEvent) error {
	log.Debugf("Deleting alert event with id %v", alevt.ID)
	_, err := agent.MainConfig.Database.DelAlertEvent(fmt.Sprintf("%v", alevt.ID))
	if err != nil {
		log.Warningf("Error deleting alert event with id %v. Error: %s", alevt.ID, err)
	}
	return err
}

//addAlertEvent Inserts current alert event into alert_event
func addAlertEvent(dev config.AlertEvent) error {
	log.Debugf("ADDING alert event %+v", dev)
	affected, err := agent.MainConfig.Database.AddAlertEvent(&dev)
	if err != nil {
		log.Warningf("Error on insert for alert event %d , affected : %+v , error: %s", dev.ID, affected, err)
	}
	return err
}

func makeTaskAlertInfo(alertkapa alert.Data, alertcfg config.AlertIDCfg, correlationID string, firsteventtime time.Time, eventtime time.Time) (TaskAlertInfo, error) {
	var taskAlertInfo = TaskAlertInfo{}
	var err error

	//from alertkapa to taskalertinfo
	jsonArByt, err := json.Marshal(alertkapa)
	if err != nil {
		log.Warningf("makeTaskAlertInfo. Error Marshalling alertkapa. Error: %s", err)
		return taskAlertInfo, err
	}
	log.Debugf("makeTaskAlertInfo. alertkapa to jsonArByt: %v", string(jsonArByt))
	err = json.Unmarshal(jsonArByt, &taskAlertInfo)
	if err != nil {
		log.Warningf("makeTaskAlertInfo. Error Unmarshalling alertkapa. Error: %s", err)
		return taskAlertInfo, err
	}

	//from alert to taskalertinfo
	jsonArByt, err = json.Marshal(alertcfg)
	if err != nil {
		log.Warningf("makeTaskAlertInfo. Error Marshalling alertcfg. Error: %s", err)
		return taskAlertInfo, err
	}
	log.Debugf("makeTaskAlertInfo. alertcfg to jsonArByt: %v", string(jsonArByt))
	err = json.Unmarshal(jsonArByt, &taskAlertInfo.ResistorAlertInfo)
	if err != nil {
		log.Warningf("makeTaskAlertInfo. Error Unmarshalling alertcfg. Error: %s", err)
		return taskAlertInfo, err
	}

	//calculated fields
	taskAlertInfo.ID = alertcfg.ID
	taskAlertInfo.ResistorAlertInfo.ID = alertcfg.ID
	taskAlertInfo.ResistorAlertInfo.InfluxDBName = kapa.GetIfxDBNameByID(alertcfg.InfluxDB)
	taskAlertInfo.Origin = "resistor"
	taskAlertInfo.CorrelationID = correlationID
	sProductTagName := alertcfg.ProductTag
	taskAlertInfo.ResistorProductTagName = sProductTagName
	taskAlertInfo.ResistorProductTagValue = alertkapa.Data.Series[0].Tags[sProductTagName]
	sIDTagName := alertcfg.IDTag
	if len(sIDTagName) == 0 {
		sIDTagName = alertcfg.ProductTag
	}
	taskAlertInfo.ResistorIDTagName = sIDTagName
	taskAlertInfo.ResistorIDTagValue = alertkapa.Data.Series[0].Tags[sIDTagName]
	taskAlertInfo.ResistorAlertTags = alertkapa.Data.Series[0].Tags
	taskAlertInfo.ResistorAlertFields = makeResistorAlertFields(alertkapa)
	if len(alertcfg.FieldDesc) > 0 {
		taskAlertInfo.ResistorAlertTriggered = fmt.Sprintf("%s : %s = ", alertcfg.InfluxMeasurement, alertcfg.FieldDesc)
	} else {
		taskAlertInfo.ResistorAlertTriggered = fmt.Sprintf("%s : %s = ", alertcfg.InfluxMeasurement, alertcfg.Field)
	}
	resistorAlertFieldValue := taskAlertInfo.ResistorAlertFields["value"]
	if resistorAlertFieldValue != nil {
		taskAlertInfo.ResistorAlertTriggered = taskAlertInfo.ResistorAlertTriggered + fmt.Sprintf("%.2f", resistorAlertFieldValue)
	}
	monExc := fmt.Sprintf("%v", taskAlertInfo.ResistorAlertFields["mon_exc"])
	taskAlertInfo.ResistorAlertInfo.ThCrit = getResistorAlertTh("crit", monExc, alertcfg)
	taskAlertInfo.ResistorAlertInfo.ThWarn = getResistorAlertTh("warn", monExc, alertcfg)
	taskAlertInfo.ResistorAlertInfo.ThInfo = getResistorAlertTh("info", monExc, alertcfg)
	taskAlertInfo.ResistorAlertInfo.ProductGroup = getResistorAlertProdGrp(alertcfg.ProductID)
	taskAlertInfo.ResistorOperationID = alertcfg.OperationID
	taskAlertInfo.ResistorOperationURL = getResistorAlertOperationURL(alertcfg.OperationID)
	taskAlertInfo.ResistorDashboardURL = makeDashboardURL(taskAlertInfo.ResistorAlertTags, alertkapa, alertcfg, firsteventtime, eventtime)

	//log json
	jsonArByt, err = json.Marshal(taskAlertInfo)
	if err != nil {
		log.Warningf("makeTaskAlertInfo. Error Marshalling taskAlertInfo. Error: %s", err)
		return taskAlertInfo, err
	}
	log.Debugf("makeTaskAlertInfo. taskAlertInfo to jsonArByt: %v", string(jsonArByt))

	return taskAlertInfo, err
}

func getResistorAlertOperationURL(operationid string) string {
	cfg, err := agent.MainConfig.Database.GetOperationCfgByID(operationid)
	if err != nil {
		log.Warningf("getResistorAlertOperationURL. Error getting OperationCfg By ID. Error: %s", err)
		return ""
	}
	return cfg.URL
}

func getResistorAlertProdGrp(productid string) string {
	productgroup := ""
	filter := "products LIKE '%\"" + productid + "\"%'"
	cfgarray, err := agent.MainConfig.Database.GetProductGroupCfgArray(filter)
	if err != nil {
		log.Warningf("getResistorAlertProdGrp. Error getting ProductGroupCfgArray. Error: %s", err)
	} else {
		for _, pg := range cfgarray {
			productgroup = productgroup + pg.ID + ","
		}
		if len(productgroup) > 0 {
			productgroup = productgroup[:len(productgroup)-1]
		}
	}
	return productgroup
}

func getResistorAlertTh(level string, monExc string, alertcfg config.AlertIDCfg) float64 {
	var value float64
	if level == "crit" {
		if monExc == "0" {
			value = alertcfg.ThCritDef
		} else if monExc == "1" {
			value = alertcfg.ThCritEx1
		} else if monExc == "2" {
			value = alertcfg.ThCritEx2
		}
	} else if level == "warn" {
		if monExc == "0" {
			value = alertcfg.ThWarnDef
		} else if monExc == "1" {
			value = alertcfg.ThWarnEx1
		} else if monExc == "2" {
			value = alertcfg.ThWarnEx2
		}
	} else if level == "info" {
		if monExc == "0" {
			value = alertcfg.ThInfoDef
		} else if monExc == "1" {
			value = alertcfg.ThInfoEx1
		} else if monExc == "2" {
			value = alertcfg.ThInfoEx2
		}
	}
	return value
}

func makeResistorAlertFields(alertkapa alert.Data) map[string]interface{} {
	raf := make(map[string]interface{})
	for idx, fieldname := range alertkapa.Data.Series[0].Columns {
		raf[fieldname] = alertkapa.Data.Series[0].Values[0][idx]
	}
	return raf
}

func makeDashboardURL(rat map[string]string, alertkapa alert.Data, alertcfg config.AlertIDCfg, firsteventtime time.Time, eventtime time.Time) string {
	var uDashboardURL *url.URL
	var err error
	sDashboardURL := ""
	log.Debugf("Entering makeDashboardURL.")

	if len(alertcfg.GrafanaServer) > 0 {
		sDashboardURL += alertcfg.GrafanaServer
	} else {
		log.Warningf("makeDashboardURL. GrafanaServer NOT informed. Empty URL will be assigned.")
		return ""
	}
	sDashboardURL += "/dashboard/db/"
	if len(alertcfg.GrafanaDashLabel) > 0 {
		sDashboardURL += alertcfg.GrafanaDashLabel
	} else {
		log.Warningf("makeDashboardURL. GrafanaDashLabel NOT informed. Empty URL will be assigned.")
		return ""
	}

	//Add time parameters
	timefrom := strconv.FormatInt(firsteventtime.Add(-2*time.Hour).Unix()*1000, 10) // firsteventtime - 2h (in ms)
	timeto := strconv.FormatInt(eventtime.Add(15*time.Minute).Unix()*1000, 10)      // eventtime + 15m (in ms)
	sDashboardURL += "?var-time_interval=$__auto_interval&from=" + timefrom + "&to=" + timeto
	//Add panelId
	if len(alertcfg.GrafanaDashPanelID) > 0 {
		sDashboardURL += "&fullscreen&panelId=" + alertcfg.GrafanaDashPanelID
	} else {
		log.Debugf("makeDashboardURL. GrafanaDashPanelID NOT informed.")
	}
	//Add all tags from kapacitor
	for tagkey, tagvalue := range rat {
		sDashboardURL += "&var-" + tagkey + "=" + tagvalue
	}
	//Add special tags from resistor
	if len(alertcfg.DeviceIDLabel) > 0 && len(alertcfg.ProductTag) > 0 {
		//Don't add duplicated vars
		_, exists := rat[alertcfg.DeviceIDLabel]
		if !exists {
			sDashboardURL += "&var-" + alertcfg.DeviceIDLabel + "=" + rat[alertcfg.ProductTag]
		}
	} else {
		log.Debugf("makeDashboardURL. DeviceIDLabel or ProductTag NOT informed.")
	}
	if len(alertcfg.ExtraLabel) > 0 && len(alertcfg.ExtraTag) > 0 {
		//Don't add duplicated vars
		_, exists := rat[alertcfg.ExtraLabel]
		if !exists {
			sDashboardURL += "&var-" + alertcfg.ExtraLabel + "=" + rat[alertcfg.ExtraTag]
		}
	} else {
		log.Debugf("makeDashboardURL. ExtraLabel or ExtraTag NOT informed.")
	}
	//Encode URL
	uDashboardURL, err = url.Parse(sDashboardURL)
	if err != nil {
		log.Warningf("makeDashboardURL. Error parsing Grafana URL. Empty URL will be assigned. Error: %s", err)
		return ""
	}
	log.Debugf("makeDashboardURL. Encoded uDashboardURL: %s", uDashboardURL.String())

	return uDashboardURL.String()
}

func sendData(taskAlertInfo TaskAlertInfo, al alert.Data, alertcfg config.AlertIDCfg, sortedtagsarray []string, endpoint config.EndpointCfg) error {
	var err error
	log.Debugf("sendData. endpoint.Type: %s, endpoint.Enabled: %v", endpoint.Type, endpoint.Enabled)
	if endpoint.Enabled {
		if endpoint.Type == "logging" {
			err = sendDataToLog(taskAlertInfo, endpoint)
		} else if endpoint.Type == "httppost" {
			err = sendDataToHTTPPost(taskAlertInfo, endpoint)
		} else if endpoint.Type == "slack" {
			err = sendDataToSlack(taskAlertInfo, endpoint)
		} else if endpoint.Type == "email" {
			err = sendDataToEmail(taskAlertInfo, endpoint)
		}
	}
	return err
}

func sendDataToEmail(al TaskAlertInfo, endpoint config.EndpointCfg) error {
	log.Debugf("sendDataToEmail. endpoint.Type: %s, endpoint.Enabled: %v", endpoint.Type, endpoint.Enabled)
	config := getResSMTPConfig(endpoint)
	err := config.Validate()
	if err != nil {
		log.Warningf("sendDataToEmail. Error validating config data: %s", err)
		return err
	}
	// Send 1 email for each To address
	for _, itemto := range endpoint.To {
		config.To = []string{itemto}
		msg := al.Message + " - Triggered by: " + al.ResistorAlertTriggered
		err = res_smtp.SendEmail(config, msg, al.Details)
		if err != nil {
			log.Warningf("sendDataToEmail. Error sending email to %s: %s", config.To, err)
		} else {
			log.Debugf("sendDataToEmail. Email successfully sent to %s", config.To)
		}
	}
	return err
}

func getResSMTPConfig(endpoint config.EndpointCfg) res_smtp.Config {
	var config res_smtp.Config
	config.Enabled = endpoint.Enabled
	config.Host = endpoint.Host
	config.Port = endpoint.Port
	config.Username = endpoint.Username
	config.Password = endpoint.Password
	config.From = endpoint.From
	config.InsecureSkipVerify = endpoint.InsecureSkipVerify
	return config
}

func sendDataToHTTPPost(al TaskAlertInfo, endpoint config.EndpointCfg) error {
	log.Debugf("sendDataToHTTPPost. endpoint.ID: %+v, endpoint.URL: %+v", endpoint.ID, endpoint.URL)

	jsonStr, err := json.Marshal(al)
	if err != nil {
		log.Errorf("sendDataToHTTPPost. Error Marshalling TaskAlertInfo as JSON. Error: %+v", err)
		return err
	}
	log.Debugf("sendDataToHTTPPost. Sending jsonStr: %v", string(jsonStr))
	req, err := http.NewRequest("POST", endpoint.URL, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Errorf("sendDataToHTTPPost. Error creating POST request: %+v", err)
		return err
	}
	//Set headers
	for _, hkv := range endpoint.Headers {
		kv := strings.Split(hkv, "=")
		if len(kv) > 0 && len(kv[0]) > 0 {
			headervalue := ""
			if len(kv) > 1 {
				headervalue = kv[1]
			}
			log.Debugf("sendDataToHTTPPost. Setting header(name:value) - (%s:%s)", kv[0], headervalue)
			req.Header.Set(kv[0], headervalue)
		}
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Type", "text/plain")

	//Set basic auth
	if len(endpoint.BasicAuthUsername) > 0 && len(endpoint.BasicAuthPassword) > 0 {
		log.Debugf("sendDataToHTTPPost. Setting BasicAuth with Username: %s and pwd: *****", endpoint.BasicAuthUsername)
		req.SetBasicAuth(endpoint.BasicAuthUsername, endpoint.BasicAuthPassword)
	}

	//Get HTTP Client
	client := &http.Client{Transport: &http.Transport{Proxy: getProxyURLFunc()}}
	//Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("sendDataToHTTPPost. Error sending request: %+v", err)
		return err
	}
	defer resp.Body.Close()

	log.Debugf("sendDataToHTTPPost. response Status:%+v", resp.Status)
	log.Debugf("sendDataToHTTPPost. response Headers:%+v", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Debugf("sendDataToHTTPPost. response Body:%+v", string(body))

	return err
}

func sendDataToLog(al TaskAlertInfo, endpoint config.EndpointCfg) error {

	var err error
	log.Debugf("sendDataToLog. endpoint.LogLevel: %+v, endpoint.LogFile: %+v", endpoint.LogLevel, endpoint.LogFile)
	// New log
	logout := logrus.New()
	//Log format
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logout.Formatter = customFormatter
	customFormatter.FullTimestamp = true
	//Log level
	l, _ := logrus.ParseLevel(endpoint.LogLevel)
	logout.Level = l
	//Log file
	if len(endpoint.LogFile) > 0 {
		logConfDir, _ := filepath.Split(endpoint.LogFile)
		err = os.MkdirAll(logConfDir, 0755)
		if err != nil {
			log.Warningf("sendDataToLog. Error creating logConfDir: %s. Error: %s", logConfDir, err)
			return err
		}
		//Log output
		f, err := os.OpenFile(endpoint.LogFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			log.Warningf("sendDataToLog. Error opening logfile: %s", err)
		} else {
			logout.Out = f
			//Log message
			logout.Debugf("Alert received from kapacitor:%+v", al)
		}
	}
	return err
}

// TaskAlertInfo represents the info of a kapacitor alert event completed with the info of the related resistor alert
type TaskAlertInfo struct {
	ID                      string                 `json:"id"`
	Message                 string                 `json:"message"`
	Details                 string                 `json:"details"`
	Time                    time.Time              `json:"time"`
	Level                   alert.Level            `json:"level"`
	CorrelationID           string                 `json:"correlationid"`
	Origin                  string                 `json:"origin"`
	ResistorOperationID     string                 `json:"resistor-operationid"`
	ResistorOperationURL    string                 `json:"resistor-operationurl"`
	ResistorDashboardURL    string                 `json:"resistor-dashboardurl"`
	ResistorIDTagName       string                 `json:"resistor-id-tag-name"`
	ResistorIDTagValue      string                 `json:"resistor-id-tag-value"`
	ResistorProductTagName  string                 `json:"resistor-product-tag-name"`
	ResistorProductTagValue string                 `json:"resistor-product-tag-value"`
	ResistorAlertTriggered  string                 `json:"resistor-alert-triggered"`
	ResistorAlertTags       map[string]string      `json:"resistor-alert-tags,omitempty"`
	ResistorAlertFields     map[string]interface{} `json:"resistor-alert-fields,omitempty"`
	ResistorAlertInfo       config.AlertIDCfgJSON  `json:"resistor-alert-info,omitempty"`
}

//SlackConfig data for Slack
type SlackConfig struct {
	// Whether Slack integration is enabled.
	Enabled bool `json:"enabled" override:"enabled"`
	// The Slack webhook URL, can be obtained by adding Incoming Webhook integration.
	URL string `json:"url" override:"url,redact"`
	// The default channel, can be overridden per alert.
	Channel string `json:"channel" override:"channel"`
	// The username of the Slack bot.
	// Default: kapacitor
	Username string `json:"username" override:"username"`
	// IconEmoji uses an emoji instead of the normal icon for the message.
	// The contents should be the name of an emoji surrounded with ':', i.e. ':chart_with_upwards_trend:'
	IconEmoji string `json:"icon-emoji" override:"icon-emoji"`
	// Whether all alerts should automatically post to slack
	Global bool `json:"global" override:"global"`
	// Whether all alerts should automatically use stateChangesOnly mode.
	// Only applies if global is also set.
	StateChangesOnly bool `json:"state-changes-only" override:"state-changes-only"`

	// Path to CA file
	SSLCA string `json:"ssl-ca" override:"ssl-ca"`
	// Path to host cert file
	SSLCert string `json:"ssl-cert" override:"ssl-cert"`
	// Path to cert key file
	SSLKey string `json:"ssl-key" override:"ssl-key"`
	// Use SSL but skip chain & host verification
	InsecureSkipVerify bool `json:"insecure-skip-verify" override:"insecure-skip-verify"`
}

//Diagnostic data for Slack
type Diagnostic interface {
	WithContext(ctx ...keyvalue.T) Diagnostic

	InsecureSkipVerify()

	Error(msg string, err error)
}

//Service data for Slack
type Service struct {
	configValue atomic.Value
	clientValue atomic.Value
	diag        Diagnostic
	client      *http.Client
}

func sendDataToSlack(al TaskAlertInfo, endpoint config.EndpointCfg) error {

	slackConfig := SlackConfig{}
	slackConfig.Enabled = endpoint.Enabled
	slackConfig.URL = endpoint.URL
	slackConfig.Channel = endpoint.Channel
	slackConfig.Username = endpoint.SlackUsername
	slackConfig.IconEmoji = endpoint.IconEmoji
	slackConfig.SSLCA = endpoint.SslCa
	slackConfig.SSLCert = endpoint.SslCert
	slackConfig.SSLKey = endpoint.SslKey
	slackConfig.InsecureSkipVerify = endpoint.InsecureSkipVerify
	log.Debugf("slackConfig: %+v", slackConfig)
	var diag Diagnostic
	s, err := NewService(slackConfig, diag)
	if err != nil {
		log.Warningf("sendDataToSlack. Error creating slack service. Error: %v", err)
		return err
	}
	if slackConfig.Enabled {
		msg := al.Message + " - Triggered by: " + al.ResistorAlertTriggered
		s.Alert(al, slackConfig.Channel, msg, slackConfig.Username, slackConfig.IconEmoji, al.Level)
	}
	return err
}

func getProxyURLFunc() func(*http.Request) (*url.URL, error) {
	proxyURLFunc := http.ProxyFromEnvironment
	proxyURLStr := agent.MainConfig.Endpoints.ProxyURL
	log.Debugf("getProxyURLFunc. proxyURLStr: %s", proxyURLStr)
	if len(proxyURLStr) > 0 {
		proxyURL, err := url.Parse(proxyURLStr)
		if err != nil {
			log.Warningf("getProxyURLFunc. Error parsing proxyURLStr: %s. Error: %v", proxyURLStr, err)
			return proxyURLFunc
		}
		proxyURLFunc = http.ProxyURL(proxyURL)
	}
	return proxyURLFunc
}

//NewService function for Slack
func NewService(c SlackConfig, d Diagnostic) (*Service, error) {
	tlsConfig, err := Create(c.SSLCA, c.SSLCert, c.SSLKey, c.InsecureSkipVerify)
	if err != nil {
		return nil, err
	}
	if (len(c.SSLCA) > 0 || len(c.SSLCert) > 0) && c.InsecureSkipVerify {
		log.Debugf("NewService. Calling d.InsecureSkipVerify()")
		d.InsecureSkipVerify()
	}
	log.Debugf("NewService. Diagnostic: %v", d)
	s := &Service{
		diag: d,
	}
	s.configValue.Store(c)
	s.clientValue.Store(&http.Client{
		Transport: &http.Transport{
			Proxy:           getProxyURLFunc(),
			TLSClientConfig: tlsConfig,
		},
	})
	log.Debugf("NewService. Service: %+v", s)
	return s, nil
}

// Create creates a new tls.Config object from the given certs, key, and CA files.
func Create(
	SSLCA, SSLCert, SSLKey string,
	InsecureSkipVerify bool,
) (*tls.Config, error) {
	t := &tls.Config{
		InsecureSkipVerify: InsecureSkipVerify,
	}
	if SSLCert != "" && SSLKey != "" {
		cert, err := tls.LoadX509KeyPair(SSLCert, SSLKey)
		if err != nil {
			return nil, fmt.Errorf(
				"Could not load TLS client key/certificate: %s",
				err)
		}
		t.Certificates = []tls.Certificate{cert}
	} else if SSLCert != "" {
		return nil, errors.New("Must provide both key and cert files: only cert file provided")
	} else if SSLKey != "" {
		return nil, errors.New("Must provide both key and cert files: only key file provided")
	}

	if SSLCA != "" {
		caCert, err := ioutil.ReadFile(SSLCA)
		if err != nil {
			return nil, fmt.Errorf("Could not load TLS CA: %s",
				err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		t.RootCAs = caCertPool
	}
	return t, nil
}

//Alert function for Slack
func (s *Service) Alert(al TaskAlertInfo, channel, message, username, iconEmoji string, level alert.Level) error {
	url, post, err := s.preparePost(al, channel, message, username, iconEmoji, level)
	if err != nil {
		log.Warningf("Slack Alert. Error preparing post to slack. Error: %v", err)
		return err
	}
	client := s.clientValue.Load().(*http.Client)
	resp, err := client.Post(url, "application/json", post)
	if err != nil {
		log.Warningf("Slack Alert. Error sending post to slack. Error: %v", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warningf("Slack Alert. Error reading response body from slack. Error: %v", err)
			return err
		}
		type response struct {
			Error string `json:"error"`
		}
		r := &response{Error: fmt.Sprintf("failed to understand Slack response. code: %d content: %s", resp.StatusCode, string(body))}
		b := bytes.NewReader(body)
		dec := json.NewDecoder(b)
		dec.Decode(r)
		return errors.New(r.Error)
	}
	return nil
}

func (s *Service) preparePost(al TaskAlertInfo, channel, message, username, iconEmoji string, level alert.Level) (string, io.Reader, error) {
	c := s.config()

	if !c.Enabled {
		return "", nil, errors.New("service is not enabled")
	}
	if channel == "" {
		channel = c.Channel
	}
	var color string
	switch level {
	case alert.Warning:
		color = "warning"
	case alert.Critical:
		color = "danger"
	case alert.Info:
		color = "#439FE0"
	default:
		color = "good"
	}

	postData := make(map[string]interface{})
	postData["as_user"] = false
	postData["channel"] = channel
	postData["text"] = "*" + al.ID + "*"

	//summaryAttachment
	summaryAttachment := attachment{}
	summaryAttachment.Color = color
	summaryAttachment.Text = ""

	fieldAttachLevel := attachfield{}
	fieldAttachLevel.Title = level.String()
	fieldAttachLevel.Value = ""
	fieldAttachLevel.Short = false
	summaryAttachment.Fields = append(summaryAttachment.Fields, fieldAttachLevel)

	fieldAttachInfo := attachfield{}
	fieldAttachInfo.Title = al.ResistorAlertTriggered
	fieldAttachInfo.Value = "<" + al.ResistorDashboardURL + "|DASHBOARD LINK>"
	fieldAttachInfo.Short = false
	summaryAttachment.Fields = append(summaryAttachment.Fields, fieldAttachInfo)

	//tagsAttachment
	tagsAttachment := attachment{}
	tagsAttachment.Color = color
	tagsAttachment.Text = ""

	fieldAttachTags := attachfield{}
	fieldAttachTags.Title = "TAGS"
	fieldAttachTags.Value = ""
	fieldAttachTags.Short = false
	tagsAttachment.Fields = append(tagsAttachment.Fields, fieldAttachTags)

	for tagname, tagvalue := range al.ResistorAlertTags {
		fieldAttachTag := attachfield{}
		fieldAttachTag.Title = tagname
		fieldAttachTag.Value = tagvalue
		fieldAttachTag.Short = true
		tagsAttachment.Fields = append(tagsAttachment.Fields, fieldAttachTag)
	}

	//fieldsAttachment
	fieldsAttachment := attachment{}
	fieldsAttachment.Color = color
	fieldsAttachment.Text = ""

	fieldAttachFields := attachfield{}
	fieldAttachFields.Title = "FIELDS"
	fieldAttachFields.Value = ""
	fieldAttachFields.Short = false
	fieldsAttachment.Fields = append(fieldsAttachment.Fields, fieldAttachFields)

	for fieldname, fieldvalue := range al.ResistorAlertFields {
		fieldAttachField := attachfield{}
		fieldAttachField.Title = fieldname
		if fieldname == "value" {
			fieldAttachField.Value = fmt.Sprintf("%.2f", fieldvalue)
		} else {
			fieldAttachField.Value = fmt.Sprintf("%v", fieldvalue)
		}
		fieldAttachField.Short = true
		fieldsAttachment.Fields = append(fieldsAttachment.Fields, fieldAttachField)
	}

	//attachmentsarray
	attachmentsarray := []attachment{}
	attachmentsarray = append(attachmentsarray, summaryAttachment)
	attachmentsarray = append(attachmentsarray, tagsAttachment)
	attachmentsarray = append(attachmentsarray, fieldsAttachment)
	postData["attachments"] = attachmentsarray

	if username == "" {
		username = c.Username
	}
	postData["username"] = username

	if iconEmoji == "" {
		iconEmoji = c.IconEmoji
	}
	postData["icon_emoji"] = iconEmoji

	var post bytes.Buffer
	enc := json.NewEncoder(&post)
	err := enc.Encode(postData)
	if err != nil {
		log.Warningf("Slack preparePost. Error encoding post data for slack. Error: %v", err)
		return "", nil, err
	}

	return c.URL, &post, nil
}

func (s *Service) config() SlackConfig {
	return s.configValue.Load().(SlackConfig)
}

// slack attachment info
type attachment struct {
	Color  string        `json:"color"`
	Text   string        `json:"text"`
	Fields []attachfield `json:"fields"`
}

type attachfield struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}
