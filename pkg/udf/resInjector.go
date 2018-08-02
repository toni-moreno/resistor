package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/influxdata/kapacitor/udf/agent"
	"github.com/spf13/viper"
	"github.com/toni-moreno/resistor/pkg/config"
)

/*
*
*
*
 */

//GeneralConfig --pending comment--
type GeneralConfig struct {
	InstanceID string `toml:"instanceID"`
	LogDir     string `toml:"logdir"`
	HomeDir    string `toml:"homedir"`
	DataDir    string `toml:"datadir"`
	LogLevel   string `toml:"loglevel"`
}

//DatabaseCfg --pending comment--
type DatabaseCfg struct {
	Type       string `toml:"type"`
	Host       string `toml:"host"`
	Name       string `toml:"name"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	SQLLogFile string `toml:"sqllogfile"`
	Debug      string `toml:"debug"`
	x          *xorm.Engine
	Path       string        `toml:"path"`
	Period     time.Duration `toml:"period"`
}

//Config --pending comment--
type Config struct {
	General  GeneralConfig
	Database DatabaseCfg
}

var (
	cfg        Config
	log        = logrus.New()
	appdir     = os.Getenv("PWD")
	logDir     = filepath.Join(appdir, "log")
	confDir    = filepath.Join(appdir, "conf")
	dataDir    = confDir
	configFile = filepath.Join(confDir, "resinjector.toml")
	socketFile = "/tmp/resInjector.sock"
	// now load up config settings
	homeDir string
	pidFile string
	//DevDB Map with device stats
	DevDB      map[string][]config.DeviceStatCfg
	mutex      sync.RWMutex
	getversion bool
	// Version Binary version
	Version string
	// Commit Git short commit id on build time
	Commit string
	// Branch Git branch on build time
	Branch string
	// BuildStamp time stamp on build time
	BuildStamp string
)

func writePIDFile() {
	if pidFile == "" {
		return
	}

	// Ensure the required directory structure exists.
	err := os.MkdirAll(filepath.Dir(pidFile), 0700)
	if err != nil {
		log.Fatal(3, "Failed to verify pid directory", err)
	}

	// Retrieve the PID and write it.
	pid := strconv.Itoa(os.Getpid())
	if err := ioutil.WriteFile(pidFile, []byte(pid), 0644); err != nil {
		log.Fatal(3, "Failed to write pidfile", err)
	}
}

func flags() *flag.FlagSet {
	var f flag.FlagSet
	f.BoolVar(&getversion, "version", getversion, "display the version")
	f.StringVar(&configFile, "config", configFile, "config file")
	f.StringVar(&logDir, "logs", logDir, "log directory")
	f.StringVar(&dataDir, "data", dataDir, "Data directory")
	f.StringVar(&pidFile, "pidfile", pidFile, "path to pid file")
	f.StringVar(&socketFile, "socket", socketFile, "path to create the unix socket")
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		f.VisitAll(func(flag *flag.Flag) {
			format := "%10s: %s\n"
			fmt.Fprintf(os.Stderr, format, "-"+flag.Name, flag.Usage)
		})
		fmt.Fprintf(os.Stderr, "\nAll settings can be set in config file: %s\n", configFile)
		os.Exit(1)

	}
	return &f
}

//init Reads configuration file
func init() {
	//Log format
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.Formatter = customFormatter
	customFormatter.FullTimestamp = true
	//fmt.Printf("Init Default directories : \n   - Exec: %s\n   - Config: %s\n   -Logs: %s\n -Data: %s\n", appdir, confDir, logDir, dataDir)

	// parse first time to see if config file is being specified
	f := flags()
	f.Parse(os.Args[1:])
	//fmt.Printf("After flags Default directories : \n   - Exec: %s\n   - Config: %s\n   -Logs: %s\n -Data: %s\n", appdir, confDir, logDir, dataDir)

	if getversion {
		t, _ := strconv.ParseInt(BuildStamp, 10, 64)
		fmt.Printf("resinjector v%s (git: %s ) built at [%s]\n", Version, Commit, time.Unix(t, 0).Format("2006-01-02 15:04:05"))
		os.Exit(0)
	}

	// now load up config settings
	if _, err := os.Stat(configFile); err == nil {
		viper.SetConfigFile(configFile)
		confDir = filepath.Dir(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/resistor/")
		viper.AddConfigPath("/opt/resistor/conf/")
		viper.AddConfigPath("./conf/")
		viper.AddConfigPath(".")
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Warnf("Fatal error config file: %s \n", err)
		os.Exit(1)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Warnf("Fatal error config file: %s \n", err)
		os.Exit(1)
	}
	//cfg := &MainConfig
	//fmt.Printf("After reading config Default directories : \n   %+v\n", cfg.General)

	if len(cfg.General.LogDir) > 0 {
		logDir = cfg.General.LogDir
	}
	if len(cfg.General.LogLevel) > 0 {
		l, _ := logrus.ParseLevel(cfg.General.LogLevel)
		log.Level = l
	}
	if len(cfg.General.DataDir) > 0 {
		dataDir = cfg.General.DataDir
	}
	if len(cfg.General.HomeDir) > 0 {
		homeDir = cfg.General.HomeDir
	}

	// parse again to overwrite values received as parameters
	f = flags()
	f.Parse(os.Args[1:])
	//fmt.Printf("Exiting init with directories : \n   - Exec: %s\n   - Config: %s\n   -Logs: %s\n -Data: %s\n", appdir, confDir, logDir, dataDir)

	if len(logDir) > 0 {
		os.Mkdir(logDir, 0755)
		//Log output
		f, _ := os.OpenFile(logDir+"/resinjector.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		log.Out = f
	}
	if len(dataDir) > 0 {
		os.Mkdir(dataDir, 0755)
	}
	//Initialize the DB configuration
	err = cfg.Database.InitDB()
	if err != nil {
		log.Warnf("Fatal error on InitDB: %s \n", err)
		os.Exit(1)
	}

}

//InitDB initialize the DB configuration
func (dbc *DatabaseCfg) InitDB() error {
	// Create ORM engine and database
	var err error
	var dbtype string
	var datasource string

	log.Debugf("Database config: %+v", dbc)

	switch dbc.Type {
	case "sqlite3":
		dbtype = "sqlite3"
		datasource = dataDir + "/" + dbc.Name + ".db"
	case "mysql":
		dbtype = "mysql"
		protocol := "tcp"
		if strings.HasPrefix(dbc.Host, "/") {
			protocol = "unix"
		}
		datasource = fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8", dbc.User, dbc.Password, protocol, dbc.Host, dbc.Name)

	default:
		log.Warnf("unknown db type %s", dbc.Type)
		return fmt.Errorf("unknown db type %s", dbc.Type)
	}

	dbc.x, err = xorm.NewEngine(dbtype, datasource)
	if err != nil {
		log.Warnf("Fail to create engine: %v\n", err)
	}

	if len(dbc.SQLLogFile) != 0 {
		dbc.x.ShowSQL(true)
		f, err := os.Create(logDir + "/" + dbc.SQLLogFile)
		if err != nil {
			log.Warnf("Fail to create log file: %s.", err)
		}
		dbc.x.SetLogger(xorm.NewSimpleLogger(f))
	}
	if dbc.Debug == "true" {
		dbc.x.Logger().SetLevel(core.LOG_DEBUG)
	}
	return err
}

// resInjector
//     @resInjector()
//        .alertId(ID)
//        .productID(ID_PRODUCT)
//        .searchByTag(DEVICEID_TAG)
//        .setLine(Line)
//        .timeCrit(weekdays,hourmin,hourmax) => check_crit (true/false)
//        .timeWarn(weekdays,hourmin,hourmax) => check_warn (true/false)
//        .timeInfo(weekdays,hourmin,hourmax) => check_info (true/false)
// 		  .injectAsTag()
// 	Generates as boolean/integer Fields or tags with the following data
// 			mon_exc = integer
// 			check_crit = true/false
//      	check_warn = true/false
// 			check_info = true/false

type resInjectorHandler struct {
	agent       *agent.Agent
	alertId     string
	productID   string
	searchByTag string
	injectAsTag bool
	critHmin    int
	critHmax    int
	warnHmin    int
	warnHmax    int
	infoHmin    int
	infoHmax    int
	critWeekDay string
	warnWeekDay string
	infoWeekDay string
	line        string
}

func newresInjectorHandler(agent *agent.Agent) *resInjectorHandler {
	return &resInjectorHandler{agent: agent}
}

// Return the InfoResponse. Describing the properties of this UDF agent.
func (*resInjectorHandler) Info() (*agent.InfoResponse, error) {
	info := &agent.InfoResponse{
		Wants:    agent.EdgeType_STREAM,
		Provides: agent.EdgeType_STREAM,
		Options: map[string]*agent.OptionInfo{
			"alertId":     {ValueTypes: []agent.ValueType{agent.ValueType_STRING}},
			"productID":   {ValueTypes: []agent.ValueType{agent.ValueType_STRING}},
			"searchByTag": {ValueTypes: []agent.ValueType{agent.ValueType_STRING}},
			"setLine":     {ValueTypes: []agent.ValueType{agent.ValueType_STRING}},
			"timeCrit":    {ValueTypes: []agent.ValueType{agent.ValueType_STRING, agent.ValueType_INT, agent.ValueType_INT}},
			"timeWarn":    {ValueTypes: []agent.ValueType{agent.ValueType_STRING, agent.ValueType_INT, agent.ValueType_INT}},
			"timeInfo":    {ValueTypes: []agent.ValueType{agent.ValueType_STRING, agent.ValueType_INT, agent.ValueType_INT}},
			"injectAsTag": {},
		},
	}
	return info, nil
}

// Init Initialize the handler based of the provided options.
// required options
// alertId string
// searchByTag string
// default values if options are not provided
// m.productID = ""
// m.critHmax = 23
// m.warnHmax = 23
// m.infoHmax = 23
// m.critHmin = 0
// m.warnHmin = 0
// m.infoHmin = 0
// m.critWeekDay = "0123456" //all days
// m.warnWeekDay = "0123456" //all days
// m.infoWeekDay = "0123456" //all days
// m.line = "LB"
// m.injectAsTag = false
func (m *resInjectorHandler) Init(r *agent.InitRequest) (*agent.InitResponse, error) {
	init := &agent.InitResponse{
		Success: true,
		Error:   "",
	}
	//default time values
	m.critHmax = 23
	m.warnHmax = 23
	m.infoHmax = 23
	m.critHmin = 0
	m.warnHmin = 0
	m.infoHmin = 0
	m.critWeekDay = "0123456" //all days
	m.warnWeekDay = "0123456" //all days
	m.infoWeekDay = "0123456" //all days
	m.line = "LB"
	m.productID = ""

	for _, opt := range r.Options {
		log.Infof("Init options: %+v", opt)
		switch opt.Name {
		case "alertId":
			m.alertId = opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue
		case "productID":
			m.productID = opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue
		case "searchByTag":
			m.searchByTag = opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue
		case "injectAsTag":
			//m.injectAsTag = opt.Values[0].Value.(*agent.OptionValue_BoolValue).BoolValue
			m.injectAsTag = true
		case "setLine":
			m.line = opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue
		case "timeCrit":
			m.critWeekDay = opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue
			m.critHmax = int(opt.Values[1].Value.(*agent.OptionValue_IntValue).IntValue)
			m.critHmin = int(opt.Values[2].Value.(*agent.OptionValue_IntValue).IntValue)
		case "timeWarn":
			m.warnWeekDay = opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue
			m.warnHmax = int(opt.Values[1].Value.(*agent.OptionValue_IntValue).IntValue)
			m.warnHmin = int(opt.Values[2].Value.(*agent.OptionValue_IntValue).IntValue)
		case "timeInfo":
			m.infoWeekDay = opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue
			m.infoHmax = int(opt.Values[1].Value.(*agent.OptionValue_IntValue).IntValue)
			m.infoHmin = int(opt.Values[2].Value.(*agent.OptionValue_IntValue).IntValue)
		}
	}

	if m.alertId == "" {
		init.Success = false
		init.Error += " must supply an AlertId"
		log.Warnf("Error on init: %s", init.Error)
	}
	if m.searchByTag == "" {
		init.Success = false
		init.Error += " must supply SearchByTag"
		log.Warnf("Error on init: %s", init.Error)
	}
	//if not set injectAsTag will be false by default

	return init, nil
}

// Create a snapshot of the running state of the process.
func (*resInjectorHandler) Snapshot() (*agent.SnapshotResponse, error) {
	return &agent.SnapshotResponse{}, nil
}

// Restore a previous snapshot.
func (*resInjectorHandler) Restore(req *agent.RestoreRequest) (*agent.RestoreResponse, error) {
	return &agent.RestoreResponse{
		Success: true,
	}, nil
}

// Start working with the next batch
func (*resInjectorHandler) BeginBatch(begin *agent.BeginBatch) error {
	return errors.New("batching not supported")
}

// ApplyRules Apply Rules from DB to received point
// adding fields or tags to the point with the following data
// 			mon_exc = integer
// 			check_crit = true/false
//      	check_warn = true/false
// 			check_info = true/false
func (m *resInjectorHandler) ApplyRules(deviceid string, rules []config.DeviceStatCfg, p *agent.Point, critok bool, warnok bool, infook bool) {
	defer timeTrack(time.Now(), "ApplyRules")
	log.Debugf("Entering ApplyRules with: deviceid: %s. AlertId: %s. ProductID: %s. rules: %+v. point: %+v.", deviceid, m.alertId, m.productID, rules, p)
	for i, r := range rules {
		log.Debugf("Rule number %d: %+v.", i, r)
		ruleAlertID := regexp.MustCompile(r.AlertID)
		if ruleAlertID.MatchString(m.alertId) {
			log.Debugf("AlertId %s received from kapacitor matches AlertId %s from rules.", m.alertId, r.AlertID)

			//ProductID
			if len(m.productID) > 0 {
				ruleProductID := regexp.MustCompile(r.ProductID)
				if ruleProductID.MatchString(m.productID) {
					log.Debugf("ProductID %s received from kapacitor matches ProductID %s from rules.", m.productID, r.ProductID)
				} else {
					//next iteration => this point should not be applied
					log.Debugf("ProductID %s received from kapacitor does not match ProductID %s from rules.", m.productID, r.ProductID)
					continue
				}
			}

			//check if any other tag to apply a filter match
			if len(r.FilterTagKey) > 0 && len(r.FilterTagValue) > 0 {
				log.Debugf("a new filter (%s=%s) exists. Checking...", r.FilterTagKey, r.FilterTagValue)
				// a new filter exist and we should check if this point has
				if tagval, oktag := p.Tags[r.FilterTagKey]; oktag {
					//tag exists in this serie
					ruleTagValue := regexp.MustCompile(r.FilterTagValue)
					if !ruleTagValue.MatchString(tagval) {
						//next iteration => this tag should not be applied
						log.Debugf("Expecting tag value %s and this point has %s", r.FilterTagValue, tagval)
						continue
					} else {
						log.Debugf("Tag (%s=%s) Matching OK with rule (%s=%s) !!!", r.FilterTagKey, tagval, r.FilterTagKey, r.FilterTagValue)
					}
				} else {
					//tag doesn't exist
					log.Warnf("There is not expected tag %s in this AlertID %s in this point.", r.FilterTagKey, m.alertId)
					continue
				}
				//in this point filter tag is matching so we will continue with tag injection
			}

			//if r.Exception=-1, then tmpbool=false, then kapacitor will not send alert
			tmpbool := r.Active && (r.ExceptionID >= 0) && strings.Contains(r.BaseLine, m.line)

			log.Debugf("Active: %t | ExceptionID : %d | LINE: %s (%s) | RESULT :%t", r.Active, r.ExceptionID, r.BaseLine, m.line, tmpbool)
			log.Debugf("CRIT : %t", critok)
			log.Debugf("WARN : %t", warnok)
			log.Debugf("INFO : %t", infook)

			if m.injectAsTag {
				// inject data as tags
				if tmpbool && critok {
					p.Tags["check_crit"] = "1"
				} else {
					p.Tags["check_crit"] = "0"
				}
				if tmpbool && warnok {
					p.Tags["check_warn"] = "1"
				} else {
					p.Tags["check_warn"] = "0"
				}
				if tmpbool && infook {
					p.Tags["check_info"] = "1"
				} else {
					p.Tags["check_info"] = "0"
				}
				p.Tags["mon_exc"] = strconv.FormatInt(r.ExceptionID, 10)
				log.Infof("Point received with data: %+v.\nCalling parameters: %+v.\nRule applied: %+v.\nInjected values: [check_crit: %s, check_warn: %s, check_info : %s, mon_exc: %s].",
					p, m, r,
					p.Tags["check_crit"],
					p.Tags["check_warn"],
					p.Tags["check_info"],
					p.Tags["mon_exc"])
			} else {
				//inject data as Fields

				if p.FieldsBool == nil {
					p.FieldsBool = make(map[string]bool)
				}
				p.FieldsBool["check_crit"] = tmpbool && critok
				p.FieldsBool["check_warn"] = tmpbool && warnok
				p.FieldsBool["check_info"] = tmpbool && infook

				if p.FieldsInt == nil {
					p.FieldsInt = make(map[string]int64)
				}
				p.FieldsInt["mon_exc"] = r.ExceptionID
				log.Infof("Point received with data: %+v.\nCalling parameters: %+v.\nRule applied: %+v.\nInjected values: [check_crit: %t, check_warn: %t, check_info : %t, mon_exc: %d].",
					p, m, r,
					p.FieldsBool["check_crit"],
					p.FieldsBool["check_warn"],
					p.FieldsBool["check_info"],
					p.FieldsInt["mon_exc"])
			}

			log.Debugf("Applying DATA for device: %s,  AlertID: %s, ProductID: %s, (Filter: %s/%s ):  [crit: %t| warn: %t| info : %t | exc: %d]",
				deviceid,
				r.AlertID,
				r.ProductID,
				r.FilterTagKey,
				r.FilterTagValue,
				tmpbool && critok,
				tmpbool && warnok,
				tmpbool && infook,
				r.ExceptionID)

		} else {
			log.Debugf("AlertId %s received from kapacitor does not match AlertId %s from rules.", m.alertId, r.AlertID)
		}
	}

}

//CheckTime Checks if the point has been received in a day and hour that allows sending the alert
func (m *resInjectorHandler) CheckTime(p *agent.Point) (bool, bool, bool, error) {
	tm := p.Time

	t := time.Unix(0, tm)

	log.Debugf("Time %s", t)

	h, _, _ := t.Clock()
	wd := strconv.Itoa(int(t.Weekday()))

	log.Debugf("Point WeekDay: %s Hour: %d", wd, h)

	critok := (h >= m.critHmin) && (h <= m.critHmax) && strings.Contains(m.critWeekDay, wd)
	warnok := (h >= m.warnHmin) && (h <= m.warnHmax) && strings.Contains(m.warnWeekDay, wd)
	infook := (h >= m.infoHmin) && (h <= m.infoHmax) && strings.Contains(m.infoWeekDay, wd)

	log.Debugf("Point TimeCheck CRIT: %t  WARN: %t INFO: %t", critok, warnok, infook)

	return critok, warnok, infook, nil
}

//SetDefault Sets default values on mon_exc, check_crit, check_warn and check_info
// Why not mon_exc = 0 ???
func (m *resInjectorHandler) SetDefault(p *agent.Point) {
	// 		mon exc = integer
	// 		check_crit = true/false
	//      check_warn = true/false
	// 		check_info = true/false
	defer timeTrack(time.Now(), "SetDefault")
	log.Debugf("Entering SetDefault with injectAsTag=%v and Point: %+v", m.injectAsTag, p)
	if m.injectAsTag {
		if _, ok := p.Tags["mon_exc"]; !ok {
			p.Tags["mon_exc"] = "0"
		}
		if _, ok := p.Tags["check_crit"]; !ok {
			p.Tags["check_crit"] = "1"
		}
		if _, ok := p.Tags["check_warn"]; !ok {
			p.Tags["check_warn"] = "1"
		}
		if _, ok := p.Tags["check_info"]; !ok {
			p.Tags["check_info"] = "1"
		}
	} else {
		//inject data as Fields
		if p.FieldsBool == nil {
			p.FieldsBool = make(map[string]bool)
		}

		if _, ok := p.FieldsBool["check_crit"]; !ok {
			p.FieldsBool["check_crit"] = true
		}
		if _, ok := p.FieldsBool["check_warn"]; !ok {
			p.FieldsBool["check_warn"] = true
		}
		if _, ok := p.FieldsBool["check_info"]; !ok {
			p.FieldsBool["check_info"] = true
		}

		if p.FieldsInt == nil {
			p.FieldsInt = make(map[string]int64)
		}
		if _, ok := p.FieldsInt["mon_exc"]; !ok {
			p.FieldsInt["mon_exc"] = 0
		}
	}
	log.Debugf("Exiting SetDefault with Point: %+v", p)
}

//Point Sends back the received point if it pass filter rules specified on injectdb
func (m *resInjectorHandler) Point(p *agent.Point) error {
	// Send back the point we just received
	defer timeTrack(time.Now(), "Point")
	log.Debugf("Receiving POINT with data: %+v", p)
	critok, warnok, infook, err := m.CheckTime(p)
	if err != nil {
		return err
	}

	var deviceid string
	var ok bool

	if deviceid, ok = p.Tags[m.searchByTag]; !ok {
		return fmt.Errorf("Tag %s doesn't exist in point", m.searchByTag)
	}
	//check if exist on the db
	mutex.RLock()
	//Generic rules first
	if rules, ok := DevDB["*"]; ok {
		m.ApplyRules(deviceid, rules, p, critok, warnok, infook)
		//do something here
	} else {
		log.Infof("there are no generic rules for device %s.", deviceid)
	}
	//specific rules after
	if rules, ok := DevDB[deviceid]; ok {
		m.ApplyRules(deviceid, rules, p, critok, warnok, infook)
		//do something here
	} else {
		log.Infof("there are no specific rules for device %s.", deviceid)
	}
	mutex.RUnlock()

	m.SetDefault(p)

	m.agent.Responses <- &agent.Response{
		Message: &agent.Response_Point{
			Point: p,
		},
	}
	log.Debugf("Returning POINT with data: %+v", p)

	return nil
}

func (*resInjectorHandler) EndBatch(end *agent.EndBatch) error {
	return nil
}

// Stop the handler gracefully.
func (m *resInjectorHandler) Stop() {
	close(m.agent.Responses)
}

type accepter struct {
	count int64
}

// Create a new agent/handler for each new connection.
// Count and log each new connection and termination.
func (acc *accepter) Accept(conn net.Conn) {
	count := acc.count
	acc.count++
	a := agent.New(conn, conn)
	h := newresInjectorHandler(a)
	a.Handler = h

	log.Debugf("Starting agent %d for connection", count)
	a.Start()
	go func() {
		err := a.Wait()
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("Agent for connection %d finished", count)
	}()
}

// ---------- Reload data from DB table ---
func reloadDbData() error {
	defer timeTrack(time.Now(), "reloadDbData")
	defer mutex.Unlock()
	mutex.Lock()

	var err error
	var devices []*config.DeviceStatCfg
	if err = cfg.Database.x.Where("`active` = 1").OrderBy("`deviceid`, `alertid`, `orderid`").Find(&devices); err != nil {
		log.Warnf("Getting devices from DB failed with error: %+v.", err)
	}
	DevDB = make(map[string][]config.DeviceStatCfg)
	for _, dev := range devices {
		DevDB[dev.DeviceID] = append(DevDB[dev.DeviceID], *dev)
	}
	log.Debugf("Getting devices from DB returns: %+v", DevDB)
	return err
}

func startRefreshProc() {
	log.Debugf("Init refresh Proc ...")
	t := time.NewTicker(cfg.Database.Period)
	for {
		log.Debugf("Beginning refresh proc again...")
		err := reloadDbData()
		if err != nil {
			log.Warnf("Error on reload DB data : %s", err)
		}

	LOOP:
		for {
			select {
			case <-t.C:
				log.Debugf("tick received...")
				break LOOP
			}
		}
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Infof("TIMELOG: %s took %s", name, elapsed)
}

func main() {
	writePIDFile()
	syscall.Umask(0000)
	// Create unix socket
	addr, err := net.ResolveUnixAddr("unix", socketFile)
	if err != nil {
		log.Warnf("Error on ResolveUnixAddr: %s", err)
	}
	// Remove the old socket before creating a new one
	os.Remove(socketFile)

	l, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Warnf("Error on ListenUnix: %s", err)
	}

	//begin data reload.
	go startRefreshProc()

	// Create server that listens on the socket
	s := agent.NewServer(l, &accepter{})

	// Setup signal handler to stop Server on various signals
	s.StopOnSignals(os.Interrupt, syscall.SIGTERM)

	log.Infof("Server listening on %s", addr.String())
	err = s.Serve()
	if err != nil {
		log.Warnf("Error on s.Serve(): %s", err)
	}
	log.Info("Server stopped")
}
