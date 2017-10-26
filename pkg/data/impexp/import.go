package impexp

import (
	"encoding/json"
	"errors"
	"fmt"
	//	"github.com/go-macaron/binding"
	//	"github.com/toni-moreno/resistor/pkg/config"
	"strconv"
	"time"
)

func (e *ExportData) ImportCheck() (*ExportData, error) {

	var duplicated []*ExportObject

	for i := 0; i < len(e.Objects); i++ {
		o := e.Objects[i]
		log.Debugf("Checking object %+v", o)
		if o.ObjectCfg == nil {
			o.Error = fmt.Sprintf("Error inconsistent data not ObjectCfg found on Imported data for id: %s", o.ObjectID)
			return nil, errors.New(o.Error)
		}
		raw, err := json.Marshal(o.ObjectCfg)
		if err != nil {
			o.Error = fmt.Sprintf("error on reformating object %s: error: %s ", o.ObjectID, err)
			return nil, errors.New(o.Error)
		}
		log.Debugf("RAW: %s", raw)
		switch o.ObjectTypeID {
		case "xxxxxx":
		/*	data := config.SnmpDeviceCfg{}
			json.Unmarshal(raw, &data)
			ers := binding.RawValidate(data)
			if ers.Len() > 0 {
				e, _ := json.Marshal(ers)
				o.Error = string(e)
				duplicated = append(duplicated, o)
				break
			}
			_, err := dbc.GetSnmpDeviceCfgByID(o.ObjectID)
			if err == nil {
				o.Error = fmt.Sprintf("Duplicated object %s in the database", o.ObjectID)
				duplicated = append(duplicated, o)
			}*/
		case "yyyyyyyy":
			/*data := config.InfluxCfg{}
			json.Unmarshal(raw, &data)
			ers := binding.RawValidate(data)
			if ers.Len() > 0 {
				e, _ := json.Marshal(ers)
				o.Error = string(e)
				duplicated = append(duplicated, o)
				break
			}
			_, err := dbc.GetInfluxCfgByID(o.ObjectID)
			if err == nil {
				o.Error = fmt.Sprintf("Duplicated object %s in the database", o.ObjectID)
				duplicated = append(duplicated, o)
			}*/
		case "zzzzzzzzz":
			//
		case "customfiltercfg":
			//
		case "oidconditioncfg":
			//
		case "measurementcfg":
			//
		case "snmpmetriccfg":
			//
		case "measgroupcfg":
			//
		default:
			return &ExportData{Info: e.Info, Objects: duplicated}, fmt.Errorf("Unknown type object type %s ", o.ObjectTypeID)
		}
	}

	if len(duplicated) > 0 {
		return &ExportData{Info: e.Info, Objects: duplicated}, fmt.Errorf("There is %d objects with errors in the imported file", len(duplicated))
	}

	return &ExportData{Info: e.Info, Objects: duplicated}, nil
}

func (e *ExportData) Import(overwrite bool, autorename bool) error {

	var suffix string
	if autorename == true {
		timestamp := time.Now().Unix()
		suffix = "_" + strconv.FormatInt(timestamp, 10)
	}
	log.Debugf("suffix: %s", suffix)
	for i := 0; i < len(e.Objects); i++ {
		o := e.Objects[i]
		o.Error = "" //reset error if exist becaouse we
		log.Debugf("Importing object %+v", o)
		if o.ObjectCfg == nil {
			o.Error = fmt.Sprintf("Error inconsistent data not ObjectCfg found on Imported data for id: %s", o.ObjectID)
			return errors.New(o.Error)
		}
		raw, err := json.Marshal(o.ObjectCfg)
		if err != nil {
			o.Error = fmt.Sprintf("error on reformating object %s: error: %s ", o.ObjectID, err)
			return errors.New(o.Error)
		}
		log.Debugf("RAW: %s", raw)
		switch o.ObjectTypeID {
		case "xxxxxx":
			/*
				log.Debugf("Importing snmpdevicecfg : %+v", o.ObjectCfg)
				data := config.SnmpDeviceCfg{}
				json.Unmarshal(raw, &data)
				var err error
				_, err = dbc.GetSnmpDeviceCfgByID(o.ObjectID)
				if err == nil { //value exist already in the database
					if overwrite == true {
						_, err2 := dbc.UpdateSnmpDeviceCfg(o.ObjectID, data)
						if err2 != nil {
							return fmt.Errorf("Error on overwrite object [%s] %s : %s", o.ObjectTypeID, o.ObjectID, err2)
						}
						break
					}
				}
				if autorename == true {
					data.ID = data.ID + suffix
				}
				_, err = dbc.AddSnmpDeviceCfg(data)
				if err != nil {
					return err
				}*/

		case "yyyyyyy":
			/*
				log.Debugf("Importing influxcfg : %+v", o.ObjectCfg)
				data := config.InfluxCfg{}
				json.Unmarshal(raw, &data)
				var err error
				_, err = dbc.GetInfluxCfgByID(o.ObjectID)
				if err == nil { //value exist already in the database
					if overwrite == true {
						_, err2 := dbc.UpdateInfluxCfg(o.ObjectID, data)
						if err2 != nil {
							return fmt.Errorf("Error on overwrite object [%s] %s : %s", o.ObjectTypeID, o.ObjectID, err2)
						}
						break
					}
				}
				if autorename == true {
					data.ID = data.ID + suffix
				}
				_, err = dbc.AddInfluxCfg(data)
				if err != nil {
					return err
				}
			*/
		case "zzzzzzzzzz":
			//
		default:
			return fmt.Errorf("Unknown type object type %s ", o.ObjectTypeID)
		}
	}
	return nil
}
