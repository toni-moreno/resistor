package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/influxdata/kapacitor/udf/agent"
	"github.com/spf13/viper"
)

type GeneralConfig struct {
	InstanceID string `toml:"instanceID"`
	LogDir     string `toml:"logdir"`
	HomeDir    string `toml:"homedir"`
	DataDir    string `toml:"datadir"`
	LogLevel   string `toml:"loglevel"`
}

type DatabaseCfg struct {
	/*Type       string        `toml:"type"`
	Host       string        `toml:"host"`
	Name       string        `toml:"name"`
	User       string        `toml:"user"`
	Password   string        `toml:"password"`
	SQLLogFile string        `toml:"sqllogfile"`
	Debug      string        `toml:"debug"`*/
	Path   string        `toml:"path"`
	Period time.Duration `toml:"period"`
}

type DeviceMonStat struct {
	alertId           string
	monExc            int64
	monActive         bool
	monLinia          string
	monFilterTagKey   string
	monFilterTagValue string
}

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
	// now load up config settings
	homeDir string
	pidFile string
	DevDB   map[string][]DeviceMonStat
	mutex   sync.RWMutex
)

func init() {
	//Log format
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.Formatter = customFormatter
	customFormatter.FullTimestamp = true

	if _, err := os.Stat(configFile); err == nil {
		viper.SetConfigFile(configFile)
		confDir = filepath.Dir(configFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("./conf/")
		viper.AddConfigPath(".")
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("Fatal error config file: %s \n", err)
		os.Exit(1)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Errorf("Fatal error config file: %s \n", err)
		os.Exit(1)
	}
	//cfg := &MainConfig

	if len(cfg.General.LogDir) > 0 {
		logDir = cfg.General.LogDir
		os.Mkdir(logDir, 0755)
		//Log output
		f, _ := os.OpenFile(logDir+"/resinjector.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		log.Out = f
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

}

// resInjector
//     @resInjector()
//        .alertId(ID)
//        .searchByTag(DEVICEID_TAG)
//        .setLine(Line)
//        .timeCrit(weekdays,hourmin,hourmax) => check_crit (true/false)
//        .timeWarn(weekdays,hourmin,hourmax) => check_warn (true/false)
//        .timeInfo(wekdays,hourmin,hourmax)  => check_info (true/false)
// 				.injectAsTag()
// 	Generates as a booleager/intee Fields or tags the folowinga data
// 			mon exc = integer
// 			check_crit = true/false
//      check_warn = true/false
// 		, check_info = true/false
//

type resInjectorHandler struct {
	agent       *agent.Agent
	alertId     string
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

// Initialze the handler based of the provided options.
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

	for _, opt := range r.Options {
		log.Infof("Iniciando opciones: %+v", opt)
		switch opt.Name {
		case "alertId":
			m.alertId = opt.Values[0].Value.(*agent.OptionValue_StringValue).StringValue
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
		log.Errorf("Error on init: %s", init.Error)
	}
	if m.searchByTag == "" {
		init.Success = false
		init.Error += " must supply an SearchByTag"
		log.Errorf("Error on init: %s", init.Error)
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

func (m *resInjectorHandler) ApplyRules(deviceid string, rules []DeviceMonStat, p *agent.Point, critok bool, warnok bool, infook bool) {
	for _, r := range rules {
		if r.alertId == m.alertId {
			//check if any other tag to apply a filter match
			if len(r.monFilterTagKey) > 0 && len(r.monFilterTagValue) > 0 {
				log.Debugf("a new filter (%s/%s)exist checking...", r.monFilterTagKey, r.monFilterTagValue)
				// a new filter exist and we should check if this point has
				if tagval, oktag := p.Tags[r.monFilterTagKey]; oktag {
					//tag exist in this serie
					if tagval != r.monFilterTagValue {
						//next iteration => this tags should not apply
						log.Debugf("Expecting tag value %s and this point has %s", r.monFilterTagValue, tagval)
						continue
					} else {
						log.Debug("Tag Matching OK!!!")
					}
				} else {
					//tag doesn't exist
					log.Warnf("There is not expected tag %s in this AlertiD %s  in this point ", r.monFilterTagKey)
					continue
				}
				//in this point filter tag is matching so we will continue with tag injection
			}

			//"mon_activo" == TRUE AND "mon_exc" >= 0 AND strContains("mon_linea",ID_LINIA), 1, 0)
			tmpbool := r.monActive && (r.monExc > 0) && strings.Contains(r.monLinia, m.line)

			log.Debugf("MONACTIVE: %t | MONEXC : %d | LINIA: %s (%s) | RESULT :%t", r.monActive, r.monExc, r.monLinia, m.line, tmpbool)
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
				p.Tags["mon_exc"] = strconv.FormatInt(r.monExc, 10)

			} else {
				//inject data as Fields

				if p.FieldsBool == nil {
					p.FieldsBool = make(map[string]bool)
				}
				p.FieldsBool["check_crit"] = tmpbool && critok
				p.FieldsBool["check_warn"] = tmpbool && warnok
				p.FieldsBool["check_info"] = tmpbool && warnok

				if p.FieldsInt == nil {
					p.FieldsInt = make(map[string]int64)
				}
				p.FieldsInt["mon_exc"] = r.monExc
			}

			log.Debugf("Applying DATA for device %s | %s (Filter: %s/%s ):  [crit: %t| warn: %t| info : %t | exc: %d]",
				deviceid,
				r.alertId,
				r.monFilterTagKey,
				r.monFilterTagValue,
				tmpbool && critok,
				tmpbool && warnok,
				tmpbool && warnok,
				r.monExc)
		}
	}

}

func (m *resInjectorHandler) CheckTime(p *agent.Point) (bool, bool, bool, error) {
	tm := p.Time

	t := time.Unix(0, tm)

	log.Debugf("Time %s", t)

	h, _, _ := t.Clock()
	wd := strconv.Itoa(int(t.Weekday()))

	log.Debugf("Point  Day: %s Hour: %d", wd, h)

	critok := (h >= m.critHmin) && (h <= m.critHmax) && strings.Contains(m.critWeekDay, wd)
	warnok := (h >= m.warnHmin) && (h <= m.warnHmax) && strings.Contains(m.warnWeekDay, wd)
	infook := (h >= m.infoHmin) && (h <= m.infoHmax) && strings.Contains(m.infoWeekDay, wd)

	log.Debugf("Point TimeCheck CRIT: %t  WARN: %t INFO: %t", critok, warnok, infook)

	return critok, warnok, infook, nil
}

func (m *resInjectorHandler) SetDefault(p *agent.Point) {
	// 			mon exc = integer
	// 			check_crit = true/false
	//      check_warn = true/false
	// 		, check_info = true/false

	if m.injectAsTag {
		if _, ok := p.Tags["mon_ext"]; !ok {
			p.Tags["mon_ext"] = "1"
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
			p.FieldsInt["mon_exc"] = 1
		}
	}
}

func (m *resInjectorHandler) Point(p *agent.Point) error {
	// Send back the point we just received
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
		log.Infof("there is not generic rules for device %s  ", deviceid)
	}
	//specific rules after
	if rules, ok := DevDB[deviceid]; ok {
		m.ApplyRules(deviceid, rules, p, critok, warnok, infook)
		//do something here
	} else {
		log.Infof("there is not info related to the %s device", deviceid)
	}
	mutex.RUnlock()

	m.SetDefault(p)

	m.agent.Responses <- &agent.Response{
		Message: &agent.Response_Point{
			Point: p,
		},
	}
	log.Debugf("POINT: %+v", p)
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

	log.Println("Starting agent for connection", count)
	a.Start()
	go func() {
		err := a.Wait()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Agent for connection %d finished", count)
	}()
}

// ---------- Reload data file ---

func reloadData(fn string) error {
	defer mutex.Unlock()
	mutex.Lock()

	file, err := os.Open(fn)
	defer file.Close()

	if err != nil {
		return err
	}

	DevDB = make(map[string][]DeviceMonStat)

	// Start reading from the file with a reader.

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		log.Debugf(" > %s", line)
		if strings.HasPrefix(line, "#") {
			//is a comment
			continue
		}
		dsArray := strings.Split(line, ",")
		if len(dsArray) < 7 {
			log.Warnf("Device data Line not correct: %s", line)
			continue
		}
		var Active bool
		var exc int64
		var err error
		//Active
		if len(dsArray[2]) > 0 {
			Active, err = strconv.ParseBool(dsArray[2])
			if err != nil {
				log.Warnf("Device data Line not correct Active value is not a correct boolean: %s", err)
				continue
			}
		}
		//Exc
		if len(dsArray[4]) > 0 {
			exc, err = strconv.ParseInt(dsArray[4], 10, 64)
			if err != nil {
				log.Warnf("Device data Line not correct Active value is not a correct boolean: %s", err)
				continue
			}
		}

		DevDB[dsArray[0]] = append(DevDB[dsArray[0]], DeviceMonStat{alertId: dsArray[1],
			monActive:         Active,
			monLinia:          dsArray[3],
			monExc:            exc,
			monFilterTagKey:   dsArray[5],
			monFilterTagValue: dsArray[6]})

	}
	if err := scanner.Err(); err != nil {
		log.Warnf("reading standard input: %s", err)
	}

	log.Debugf("DATADB: %+v", DevDB)

	return nil
}

func startRefreshProc() {
	log.Printf("Init refresh Proc ...")
	t := time.NewTicker(cfg.Database.Period)
	for {
		log.Info("Beginning refresh proc again...")
		err := reloadData(cfg.Database.Path)
		if err != nil {
			log.Errorf("Error on reload data : %s", err)
		}

	LOOP:
		for {
			select {
			case <-t.C:
				log.Printf("tick received...")
				break LOOP
			}
		}
	}
}

var socketPath = flag.String("socket", "/tmp/resInjector.sock", "Where to create the unix socket")

func main() {
	flag.Parse()
	syscall.Umask(0000)
	// Create unix socket
	addr, err := net.ResolveUnixAddr("unix", *socketPath)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	log.Info("Server stopped")
}
