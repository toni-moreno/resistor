package config

import (
	"fmt"
	"strings"
	// _ needed to mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	// _ needed to sqlite3
	_ "github.com/mattn/go-sqlite3"
	"os"
	"sync/atomic"
)

func (dbc *DatabaseCfg) resetChanges() {
	atomic.StoreInt64(&dbc.numChanges, 0)
}

func (dbc *DatabaseCfg) addChanges(n int64) {
	atomic.AddInt64(&dbc.numChanges, n)
}
func (dbc *DatabaseCfg) getChanges() int64 {
	return atomic.LoadInt64(&dbc.numChanges)
}

//DbObjAction measurement groups to assign to devices
type DbObjAction struct {
	Type     string
	TypeDesc string
	ObID     string
	Action   string
}

//InitDB initialize de BD configuration
func (dbc *DatabaseCfg) InitDB() {
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
		//datasource = dbc.User + ":" + dbc.Pass + "@" + dbc.Host + "/" + dbc.Name + "?charset=utf8"
	default:
		log.Errorf("unknown db  type %s", dbc.Type)
		return
	}

	dbc.x, err = xorm.NewEngine(dbtype, datasource)
	if err != nil {
		log.Fatalf("Fail to create engine: %v\n", err)
	}

	if len(dbc.SQLLogFile) != 0 {
		dbc.x.ShowSQL(true)
		f, error := os.Create(logDir + "/" + dbc.SQLLogFile)
		if err != nil {
			log.Errorln("Fail to create log file  ", error)
		}
		dbc.x.SetLogger(xorm.NewSimpleLogger(f))
	}
	if dbc.Debug == "true" {
		dbc.x.Logger().SetLevel(core.LOG_DEBUG)
	}

	/* Sync tables
	if err = dbc.x.Sync(new(InfluxCfg)); err != nil {
		log.Fatalf("Fail to sync database InfluxCfg: %v\n", err)
	}
	if err = dbc.x.Sync(new(SnmpDeviceCfg)); err != nil {
		log.Fatalf("Fail to sync database SnmpDeviceCfg: %v\n", err)
	}
	if err = dbc.x.Sync(new(SnmpMetricCfg)); err != nil {
		log.Fatalf("Fail to sync database SnmpMetricCfg: %v\n", err)
	}
	if err = dbc.x.Sync(new(MeasurementCfg)); err != nil {
		log.Fatalf("Fail to sync database MeasurementCfg: %v\n", err)
	}
	if err = dbc.x.Sync(new(MeasFilterCfg)); err != nil {
		log.Fatalf("Fail to sync database MeasurementFilterCfg : %v\n", err)
	}
	if err = dbc.x.Sync(new(MeasurementFieldCfg)); err != nil {
		log.Fatalf("Fail to sync database MeasurementFieldCfg: %v\n", err)
	}
	if err = dbc.x.Sync(new(MGroupsCfg)); err != nil {
		log.Fatalf("Fail to sync database MGroupCfg: %v\n", err)
	}
	if err = dbc.x.Sync(new(MGroupsMeasurements)); err != nil {
		log.Fatalf("Fail to sync database MGroupsMeasurements: %v\n", err)
	}
	if err = dbc.x.Sync(new(SnmpDevMGroups)); err != nil {
		log.Fatalf("Fail to sync database SnmpDevMGroups: %v\n", err)
	}
	if err = dbc.x.Sync(new(SnmpDevFilters)); err != nil {
		log.Fatalf("Fail to sync database SnmpDevFilters: %v\n", err)
	}
	if err = dbc.x.Sync(new(CustomFilterCfg)); err != nil {
		log.Fatalf("Fail to sync database CustomFilterCfg: %v\n", err)
	}
	if err = dbc.x.Sync(new(CustomFilterItems)); err != nil {
		log.Fatalf("Fail to sync database CustomFilterItems: %v\n", err)
	}
	if err = dbc.x.Sync(new(OidConditionCfg)); err != nil {
		log.Fatalf("Fail to sync database OidConditionCfg: %v\n", err)
	}*/
}

//LoadDbConfig get data from database
func (dbc *DatabaseCfg) LoadDbConfig(cfg *SQLConfig) {
	//var err error
	return

	/*Load Influxdb databases
	cfg.Influxdb, err = dbc.GetInfluxCfgMap("")
	if err != nil {
		log.Warningf("Some errors on get Influx db's :%v", err)
	}

	//Load metrics
	cfg.Metrics, err = dbc.GetSnmpMetricCfgMap("")
	if err != nil {
		log.Warningf("Some errors on get Metrics  :%v", err)
	}

	//Load Measurements
	cfg.Measurements, err = dbc.GetMeasurementCfgMap("")
	if err != nil {
		log.Warningf("Some errors on get Measurements  :%v", err)
	}

	//Load Measurement Filters
	cfg.MFilters, err = dbc.GetMeasFilterCfgMap("")
	if err != nil {
		log.Warningf("Some errors on get Measurement Filters  :%v", err)
	}

	//Load measourement Groups

	cfg.GetGroups, err = dbc.GetMGroupsCfgMap("")
	if err != nil {
		log.Warningf("Some errors on get Measurements Groups  :%v", err)
	}

	//Device

	cfg.SnmpDevice, err = dbc.GetSnmpDeviceCfgMap("")
	if err != nil {
		log.Warningf("Some errors on get SnmpDeviceConf :%v", err)
	}
	dbc.resetChanges()
	*/
}
