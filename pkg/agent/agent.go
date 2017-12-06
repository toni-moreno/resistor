package agent

import (
	//	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/toni-moreno/resistor/pkg/agent/output"
	"github.com/toni-moreno/resistor/pkg/agent/selfmon"
	"github.com/toni-moreno/resistor/pkg/config"

	"sync"
	//	"time"
)

var (
	// Version Binari version
	Version string
	// Commit Git short commit id on build time
	Commit string
	// Branch Git branch on build time
	Branch string
	// BuildStamp time stamp on build time
	BuildStamp string
)

// RInfo  Release basic version info for the agent
type RInfo struct {
	InstanceID string
	Version    string
	Commit     string
	Branch     string
	BuildStamp string
}

// GetRInfo get release info
func GetRInfo() *RInfo {
	info := &RInfo{
		InstanceID: MainConfig.General.InstanceID,
		Version:    Version,
		Commit:     Commit,
		Branch:     Branch,
		BuildStamp: BuildStamp,
	}
	return info
}

var (

	// MainConfig has all configuration
	MainConfig config.Config

	// DBConfig db config
	DBConfig config.DBConfig

	log *logrus.Logger
	//mutex for devices map
	mutex sync.RWMutex
	//reload mutex
	reloadMutex   sync.Mutex
	reloadProcess bool

	selfmonProc *selfmon.SelfMon
	// for synchronize  deivce specific goroutines
	gatherWg sync.WaitGroup
	senderWg sync.WaitGroup
)

// SetLogger set log output
func SetLogger(l *logrus.Logger) {
	log = l
}

//Reload Mutex Related Methods.

// CheckAndSetStarted check if this thread is already working and set if not
func CheckReloadProcess() bool {
	reloadMutex.Lock()
	defer reloadMutex.Unlock()
	return reloadProcess
}

// CheckAndSetStarted check if this thread is already working and set if not
func CheckAndSetReloadProcess() bool {
	reloadMutex.Lock()
	defer reloadMutex.Unlock()
	retval := reloadProcess
	reloadProcess = true
	return retval
}

// CheckAndUnSetStarted check if this thread is already working and set if not
func CheckAndUnSetReloadProcess() bool {
	reloadMutex.Lock()
	defer reloadMutex.Unlock()
	retval := reloadProcess
	reloadProcess = false
	return retval
}

func initSelfMonitoring() {

	selfmonProc = selfmon.NewNotInit(&MainConfig.Selfmon)

	val := output.NewNotInitInfluxDB(&MainConfig.Influxdb)

	if MainConfig.Selfmon.Enabled {
		val.Init()
		val.StartSender(&senderWg)

		selfmonProc.Init()
		selfmonProc.SetOutput(val)

		log.Printf("SELFMON enabled %+v", MainConfig.Selfmon)
		//Begin the statistic reporting
		selfmonProc.StartGather(&gatherWg)

	} else {
		log.Printf("SELFMON disabled %+v\n", MainConfig.Selfmon)
	}
}

// LoadConf call to initialize alln configurations
func LoadConf() {
	//Load all database info to Cfg struct
	MainConfig.Database.LoadDbConfig(&DBConfig)
	//Prepare the InfluxDataBases Configuration
	log.Debugf("DB CONFIG LOAD :%#+v", DBConfig)

	// beginning self monitoring process if needed.( before each other gorotines could begin)

	initSelfMonitoring()

	//Initialize Device Metrics CFG

	config.Init(&DBConfig)

	//beginning  the gather process
}
