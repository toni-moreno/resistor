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
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-macaron/binding"
	"github.com/influxdata/kapacitor/alert"
	"github.com/influxdata/kapacitor/keyvalue"
	//kapaPost "github.com/influxdata/kapacitor/services/httppost"
	//kapaSlack "github.com/influxdata/kapacitor/services/slack"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	"github.com/toni-moreno/resistor/pkg/kapa"
	//"github.com/toni-moreno/resistor/pkg/data/alertfilter"

	"gopkg.in/macaron.v1"
)

// NewAPIRtKapFilter set the runtime Kapacitor filter  API
func NewAPIRtKapFilter(m *macaron.Macaron) error {

	bind := binding.Bind
	m.Group("/api/rt/kapfilter", func() {
		m.Post("/alert/:endpoint", reqAlertSignedIn, bind(alert.Data{}), RTAlertHandler)
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

	//Get AlertCfg
	alertcfg, err := agent.MainConfig.Database.GetAlertIDCfgByID(al.ID)
	if err != nil {
		log.Warningf("Error getting alert cfg with id: %s. Error: %s", al.ID, err)
	}
	sortedtagsarray := sortTagsMap(al.Data.Series[0].Tags)
	//Add alert event to list of alert events
	alertevent := makeAlertEvent(al, alertcfg, sortedtagsarray)
	AddAlertEvent(alertevent)

	//Send alert event to related endpoints
	arendpoint := alertcfg.Endpoint
	for _, endpointid := range arendpoint {
		endpoint, err := agent.MainConfig.Database.GetEndpointCfgByID(endpointid)
		if err != nil {
			log.Warningf("Error getting endpoint for id %s. Error: %s.", endpointid, err)
		} else {
			log.Debugf("Got endpoint: %+v", endpoint)
			err = sendData(al, alertcfg, sortedtagsarray, endpoint)
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

func makeAlertEvent(al alert.Data, alertcfg config.AlertIDCfg, sortedtagsarray []string) (dev config.AlertEventHist) {
	alertevent := config.AlertEventHist{}
	alertevent.ID = 0
	alertevent.AlertID = al.ID
	alertevent.Message = al.Message
	alertevent.Details = al.Details
	alertevent.Time = al.Time
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
			break
		}
	}
	return alertevent
}

// AddAlertEvent Inserts new alert event into the internal DB
func AddAlertEvent(dev config.AlertEventHist) {
	log.Debugf("ADDING alert event %+v", dev)
	affected, err := agent.MainConfig.Database.AddAlertEventHist(&dev)
	if err != nil {
		log.Warningf("Error on insert for alert event %d , affected : %+v , error: %s", dev.ID, affected, err)
	} else {
		log.Infof("Alert event %d successfully inserted", dev.ID)
	}
}

func makeTaskAlertInfo(alertkapa alert.Data, alertcfg config.AlertIDCfg, sortedtagsarray []string) (TaskAlertInfo, error) {
	var taskAlertInfo = TaskAlertInfo{}
	var err error

	//from alertkapa to taskalert
	jsonArByt, err := json.Marshal(alertkapa)
	if err != nil {
		log.Warningf("makeTaskAlertInfo. Error Marshalling alertkapa. Error: %s", err)
	}
	log.Debugf("makeTaskAlertInfo. alertkapa to jsonArByt: %v", string(jsonArByt))
	err = json.Unmarshal(jsonArByt, &taskAlertInfo)
	if err != nil {
		log.Warningf("makeTaskAlertInfo. Error Unmarshalling alertkapa. Error: %s", err)
	}

	//from alert to taskalert
	jsonArByt, err = json.Marshal(alertcfg)
	if err != nil {
		log.Warningf("makeTaskAlertInfo. Error Marshalling alertcfg. Error: %s", err)
	}
	log.Debugf("makeTaskAlertInfo. alertcfg to jsonArByt: %v", string(jsonArByt))
	err = json.Unmarshal(jsonArByt, &taskAlertInfo.ResistorAlertInfo)
	if err != nil {
		log.Warningf("makeTaskAlertInfo. Error Unmarshalling alertcfg. Error: %s", err)
	}

	//calculated fields
	taskAlertInfo.ResistorAlertInfo.ID = alertcfg.ID
	taskAlertInfo.ResistorAlertInfo.InfluxDBName = kapa.GetIfxDBNameByID(alertcfg.InfluxDB)
	taskAlertInfo.Origin = "resistor"
	taskAlertInfo.CorrelationID = "[" + alertcfg.ID + "]|" + strings.Join(sortedtagsarray, ",")
	taskAlertInfo.ResistorProductTagName = alertcfg.ProductTag
	taskAlertInfo.ResistorProductTagValue = alertkapa.Data.Series[0].Tags[alertcfg.ProductTag]
	taskAlertInfo.ResistorAlertTags = alertkapa.Data.Series[0].Tags
	taskAlertInfo.ResistorAlertFields = makeResistorAlertFields(alertkapa)
	taskAlertInfo.ResistorAlertTriggered = fmt.Sprintf("%s : %s = %.2f", alertcfg.InfluxMeasurement, alertcfg.Field, taskAlertInfo.ResistorAlertFields["value"])
	monExc := fmt.Sprintf("%v", taskAlertInfo.ResistorAlertFields["mon_exc"])
	taskAlertInfo.ResistorAlertInfo.ThCrit = getResistorAlertTh("crit", monExc, alertcfg)
	taskAlertInfo.ResistorAlertInfo.ThWarn = getResistorAlertTh("warn", monExc, alertcfg)
	taskAlertInfo.ResistorAlertInfo.ThInfo = getResistorAlertTh("info", monExc, alertcfg)
	taskAlertInfo.ResistorAlertInfo.ProductGroup = getResistorAlertProdGrp(alertcfg.ProductID)

	//log json
	jsonArByt, err = json.Marshal(taskAlertInfo)
	if err != nil {
		log.Warningf("makeTaskAlertInfo. Error Marshalling taskAlertInfo. Error: %s", err)
	}
	log.Debugf("makeTaskAlertInfo. taskAlertInfo to jsonArByt: %v", string(jsonArByt))

	return taskAlertInfo, err
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

func sendData(al alert.Data, alertcfg config.AlertIDCfg, sortedtagsarray []string, endpoint config.EndpointCfg) error {
	var err error
	strouttype := endpoint.Type
	log.Debugf("strouttype: %s", strouttype)
	//makeTaskAlertInfo
	taskAlertInfo, err := makeTaskAlertInfo(al, alertcfg, sortedtagsarray)
	if err != nil {
		log.Warningf("sendData. Error making taskAlertInfo. Error: %s", err)
	}

	if strouttype == "logging" {
		err = sendDataToLog(taskAlertInfo, endpoint)
	} else if strouttype == "httppost" {
		err = sendDataToHTTPPost(taskAlertInfo, endpoint)
	} else if strouttype == "slack" {
		err = sendDataToSlack(al, endpoint)
	}
	return err
}

func sendDataToHTTPPost(al TaskAlertInfo, endpoint config.EndpointCfg) error {
	log.Debugf("sendDataToHTTPPost. endpoint.ID: %+v, endpoint.URL: %+v", endpoint.ID, endpoint.URL)

	jsonStr, err := json.Marshal(al)
	log.Debugf("sendDataToHTTPPost. Sending jsonStr: %v", string(jsonStr))

	req, err := http.NewRequest("POST", endpoint.URL, bytes.NewBuffer(jsonStr))

	//Set headers
	for _, hkv := range endpoint.Headers {
		kv := strings.Split(hkv, "=")
		if len(kv) > 0 {
			headervalue := ""
			if len(kv) > 1 {
				headervalue = kv[1]
			}
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("sendDataToHTTPPost. Error:%+v", err)
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
	ResistorProductTagName  string                 `json:"resistor-product-tag-name"`
	ResistorProductTagValue string                 `json:"resistor-product-tag-value"`
	ResistorAlertTriggered  string                 `json:"resistor-alert-triggered"`
	ResistorAlertTags       map[string]string      `json:"resistor-alert-tags,omitempty"`
	ResistorAlertFields     map[string]interface{} `json:"resistor-alert-fields,omitempty"`
	ResistorAlertInfo       config.AlertIDCfgJSON  `json:"resistor-alert-info,omitempty"`
}

//Config data for Slack
type Config struct {
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

func sendDataToSlack(al alert.Data, endpoint config.EndpointCfg) error {

	slackConfig := Config{}
	slackConfig.Enabled = endpoint.SlackEnabled
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
	if slackConfig.Enabled {
		s.Alert(slackConfig.Channel, al.Message, slackConfig.Username, slackConfig.IconEmoji, al.Level)
	}
	return err
}

//NewService function for Slack
func NewService(c Config, d Diagnostic) (*Service, error) {
	tlsConfig, err := Create(c.SSLCA, c.SSLCert, c.SSLKey, c.InsecureSkipVerify)
	if err != nil {
		return nil, err
	}
	if tlsConfig.InsecureSkipVerify {
		d.InsecureSkipVerify()
	}
	s := &Service{
		diag: d,
	}
	s.configValue.Store(c)
	s.clientValue.Store(&http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: tlsConfig,
		},
	})
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
func (s *Service) Alert(channel, message, username, iconEmoji string, level alert.Level) error {
	url, post, err := s.preparePost(channel, message, username, iconEmoji, level)
	if err != nil {
		return err
	}
	client := s.clientValue.Load().(*http.Client)
	resp, err := client.Post(url, "application/json", post)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
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

func (s *Service) preparePost(channel, message, username, iconEmoji string, level alert.Level) (string, io.Reader, error) {
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
	default:
		color = "good"
	}
	a := attachment{
		Fallback: message,
		Text:     message,
		Color:    color,
		MrkdwnIn: []string{"text"},
	}
	postData := make(map[string]interface{})
	postData["as_user"] = false
	postData["channel"] = channel
	postData["text"] = ""
	postData["attachments"] = []attachment{a}

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
		return "", nil, err
	}

	return c.URL, &post, nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

// slack attachment info
type attachment struct {
	Fallback string   `json:"fallback"`
	Color    string   `json:"color"`
	Text     string   `json:"text"`
	MrkdwnIn []string `json:"mrkdwn_in"`
}
