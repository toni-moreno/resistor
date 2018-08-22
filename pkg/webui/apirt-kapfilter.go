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
	"strings"
	"sync/atomic"

	"github.com/Sirupsen/logrus"
	"github.com/go-macaron/binding"
	"github.com/influxdata/kapacitor/alert"
	"github.com/influxdata/kapacitor/keyvalue"
	//kapaPost "github.com/influxdata/kapacitor/services/httppost"
	//kapaSlack "github.com/influxdata/kapacitor/services/slack"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
	//"github.com/toni-moreno/resistor/pkg/data/alertfilter"

	"gopkg.in/macaron.v1"
)

// NewAPIRtKapFilter set the runtime Kapacitor filter  API
func NewAPIRtKapFilter(m *macaron.Macaron) error {

	bind := binding.Bind
	m.Group("/api/rt/kapfilter", func() {
		m.Post("/alert/:outhttp", reqAlertSignedIn, bind(alert.Data{}), RTAlertHandler)
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

	alertevent := makeAlertEvent(al)
	AddAlertEvent(alertevent)

	strouthttp := ctx.Params(":outhttp")
	log.Debugf("outhttp: %s", strouthttp)
	if len(strouthttp) > 0 {
		arouthttp := strings.Split(strouthttp, ",")
		for _, outhttpid := range arouthttp {
			outhttp, err := agent.MainConfig.Database.GetOutHTTPCfgByID(outhttpid)
			if err != nil {
				log.Warningf("Error getting outhttp for id %s. Error: %s.", outhttpid, err)
			} else {
				log.Debugf("Got outhttp: %+v", outhttp)
				err = sendData(al, outhttp)
				if err != nil {
					log.Warningf("Error sending data to endpoint with id %s. Error: %s.", outhttpid, err)
				}
			}
		}
	}

	//alertfilter.ProcessAlert(al)

	ctx.JSON(200, "DONE")
}

func makeAlertEvent(al alert.Data) (dev config.AlertEventCfg) {
	alertevent := config.AlertEventCfg{}
	alertevent.UID = 0
	alertevent.ID = al.ID
	alertevent.Message = al.Message
	alertevent.Details = al.Details
	alertevent.Time = al.Time
	alertevent.Duration = al.Duration
	alertevent.Level = al.Level.String()
	return alertevent
}

// AddAlertEvent Inserts new alert event into the internal DB
func AddAlertEvent(dev config.AlertEventCfg) {
	log.Printf("ADDING alert event %+v", dev)
	affected, err := agent.MainConfig.Database.AddAlertEventCfg(&dev)
	if err != nil {
		log.Warningf("Error on insert for alert event %s , affected : %+v , error: %s", dev.ID, affected, err)
	} else {
		log.Infof("Alert event %s successfully inserted", dev.ID)
	}
}

func sendData(al alert.Data, outhttp config.OutHTTPCfg) error {
	var err error
	strouttype := outhttp.Type
	jsonconfig := outhttp.JSONConfig
	log.Debugf("strouttype: %s", strouttype)
	if strouttype == "logging" {
		err = sendDataToLog(al, jsonconfig)
	} else if strouttype == "httppost" {
		log.Warningf("httppost pending to develop")
	} else if strouttype == "slack" {
		err = sendDataToSlack(al, jsonconfig)
	}
	return err
}

func sendDataToLog(al alert.Data, jsonconfig string) error {
	type logConfig struct {
		File  string `json:"file"`
		Level string `json:"level"`
	}
	logConf := logConfig{}

	err := json.Unmarshal([]byte(jsonconfig), &logConf)
	log.Debugf("logConf: %+v", logConf)
	// New log
	logout := logrus.New()
	//Log format
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logout.Formatter = customFormatter
	customFormatter.FullTimestamp = true
	//Log level
	l, _ := logrus.ParseLevel(logConf.Level)
	logout.Level = l
	//Log file
	if len(logConf.File) > 0 {
		logConfDir, _ := filepath.Split(logConf.File)
		err = os.MkdirAll(logConfDir, 0755)
		if err != nil {
			log.Warningf("sendDataToLog. Error creating logConfDir: %s. Error: %s", logConfDir, err)
		}
		//Log output
		f, err := os.OpenFile(logConf.File, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
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

type Diagnostic interface {
	WithContext(ctx ...keyvalue.T) Diagnostic

	InsecureSkipVerify()

	Error(msg string, err error)
}

type Service struct {
	configValue atomic.Value
	clientValue atomic.Value
	diag        Diagnostic
	client      *http.Client
}

func sendDataToSlack(al alert.Data, jsonconfig string) error {
	/*
		mslackConfig := kapaSlack.Config{}
		log.Debugf("Getting json: %+v", mslackConfig)
		mPostConf := kapaPost.Config{}
		log.Debugf("Getting json: %+v", mPostConf)
	*/

	slackConfig := Config{}

	err := json.Unmarshal([]byte(jsonconfig), &slackConfig)
	log.Debugf("slackConfig: %+v", slackConfig)
	var diag Diagnostic
	s, err := NewService(slackConfig, diag)
	log.Debugf("s: %+v, diag: %+v", s, diag)
	if slackConfig.Enabled {
		s.Alert(slackConfig.Channel, al.Message, slackConfig.Username, slackConfig.IconEmoji, al.Level)
	}
	return err
}

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
