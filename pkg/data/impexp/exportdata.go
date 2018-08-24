package impexp

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/toni-moreno/resistor/pkg/agent"
	"github.com/toni-moreno/resistor/pkg/config"
)

var (
	log     *logrus.Logger
	confDir string              //Needed to get File Filters data
	dbc     *config.DatabaseCfg //Needed to get Custom Filter  data
)

// SetConfDir  enable load File Filters from anywhere in the our FS.
func SetConfDir(dir string) {
	confDir = dir
}

// SetDB load database config to load data if needed (used in filters)
func SetDB(db *config.DatabaseCfg) {
	dbc = db
}

// SetLogger set output log
func SetLogger(l *logrus.Logger) {
	log = l
}

type ExportInfo struct {
	FileName      string
	Description   string
	Author        string
	Tags          string
	AgentVersion  string
	ExportVersion string
	CreationDate  time.Time
}

type EIOptions struct {
	Recursive   bool   //Export Option
	AutoRename  bool   //Import Option
	AlternateID string //Import Option
}

type ExportObject struct {
	ObjectTypeID string
	ObjectID     string
	Options      *EIOptions
	ObjectCfg    interface{}
	Error        string
}

// ExportData the runtime measurement config
type ExportData struct {
	Info       *ExportInfo
	Objects    []*ExportObject
	tmpObjects []*ExportObject //only for temporal use
}

func NewExport(info *ExportInfo) *ExportData {
	if len(agent.Version) > 0 {
		info.AgentVersion = agent.Version
	} else {
		info.AgentVersion = "debug"
	}

	info.ExportVersion = "1.0"
	info.CreationDate = time.Now()
	return &ExportData{
		Info: info,
	}
}

func checkIfExistOnArray(list []*ExportObject, ObjType string, id string) bool {
	for _, v := range list {
		if v.ObjectTypeID == ObjType && v.ObjectID == id {
			return true
		}
	}
	return false
}

func (e *ExportData) PrependObject(obj *ExportObject) {
	if checkIfExistOnArray(e.Objects, obj.ObjectTypeID, obj.ObjectID) == true {
		return
	}
	e.tmpObjects = append([]*ExportObject{obj}, e.tmpObjects...)
}

func (e *ExportData) UpdateTmpObject() {
	//we need remove duplicated objects on the auxiliar array
	objectList := []*ExportObject{}
	for i := 0; i < len(e.tmpObjects); i++ {
		v := e.tmpObjects[i]
		if checkIfExistOnArray(objectList, v.ObjectTypeID, v.ObjectID) == false {
			objectList = append(objectList, v)
		}
	}
	e.Objects = append(e.Objects, objectList...)
	e.tmpObjects = nil
}

// Export  exports data
func (e *ExportData) Export(ObjType string, id string, recursive bool, level int) error {

	log.Debugf("Entering Export with ObjType: %s, id: %s, recursive: %t, level: %d", ObjType, id, recursive, level)
	switch ObjType {
	case "rangetimecfg":
		//contains sensible data
		v, err := dbc.GetRangeTimeCfgByID(id)
		if err != nil {
			return err
		}
		e.PrependObject(&ExportObject{ObjectTypeID: "rangetimecfg", ObjectID: id, ObjectCfg: v})
	case "ifxservercfg":
		//contains sensible data
		v, err := dbc.GetIfxServerCfgByID(id)
		if err != nil {
			return err
		}
		e.PrependObject(&ExportObject{ObjectTypeID: "ifxservercfg", ObjectID: id, ObjectCfg: v})
	case "kapacitorcfg":
		v, err := dbc.GetKapacitorCfgByID(id)
		if err != nil {
			return err
		}
		e.PrependObject(&ExportObject{ObjectTypeID: "kapacitorcfg", ObjectID: id, ObjectCfg: v})
	case "productcfg":
		v, err := dbc.GetProductCfgByID(id)
		if err != nil {
			return err
		}
		e.PrependObject(&ExportObject{ObjectTypeID: "productcfg", ObjectID: id, ObjectCfg: v})
	case "productgroupcfg":
		//contains sensible data
		v, err := dbc.GetProductGroupCfgByID(id)
		if err != nil {
			return err
		}
		e.PrependObject(&ExportObject{ObjectTypeID: "productgroupcfg", ObjectID: id, ObjectCfg: v})
		if !recursive {
			break
		}
		for _, val := range v.Products {
			e.Export("productcfg", val, recursive, level+1)
		}
	case "outhttpcfg":
		v, err := dbc.GetOutHTTPCfgByID(id)
		if err != nil {
			return err
		}
		e.PrependObject(&ExportObject{ObjectTypeID: "outhttpcfg", ObjectID: id, ObjectCfg: v})
	case "alertcfg":
		//contains sensible data
		v, err := dbc.GetAlertIDCfgByID(id)
		if err != nil {
			return err
		}
		e.PrependObject(&ExportObject{ObjectTypeID: "alertcfg", ObjectID: id, ObjectCfg: v})
		if !recursive {
			break
		}
		for _, val := range v.OutHTTP {
			e.Export("outhttpcfg", val, recursive, level+1)
		}
		e.Export("kapacitorcfg", v.KapacitorID, recursive, level+1)
		e.Export("productcfg", v.ProductID, recursive, level+1)

		if v.TriggerType != "DEADMAN" {
			e.Export("rangetimecfg", v.ThCritRangeTimeID, recursive, level+1)
			e.Export("rangetimecfg", v.ThWarnRangeTimeID, recursive, level+1)
			e.Export("rangetimecfg", v.ThInfoRangeTimeID, recursive, level+1)
		}

	case "templatecfg":
		v, err := dbc.GetTemplateCfgByID(id)
		if err != nil {
			return err
		}
		e.PrependObject(&ExportObject{ObjectTypeID: "templatecfg", ObjectID: id, ObjectCfg: v})
	case "devicestatcfg":
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}
		v, err := dbc.GetDeviceStatCfgByID(idInt64)
		if err != nil {
			return err
		}
		e.PrependObject(&ExportObject{ObjectTypeID: "devicestatcfg", ObjectID: id, ObjectCfg: v})
	default:
		return fmt.Errorf("Unknown type object type %s ", ObjType)
	}
	if level == 0 {
		e.UpdateTmpObject()
	}
	return nil
}
