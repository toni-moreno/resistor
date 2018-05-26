package impexp

import (
	"encoding/json"
	"errors"
	"fmt"
	//	"github.com/go-macaron/binding"
	//	"github.com/toni-moreno/resistor/pkg/config"
	"strconv"
	"time"

	"github.com/go-macaron/binding"
	"github.com/toni-moreno/resistor/pkg/config"
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
		case "rangetimecfg":
			data := config.RangeTimeCfg{}
			json.Unmarshal(raw, &data)
			ers := binding.RawValidate(data)
			if ers.Len() > 0 {
				e, _ := json.Marshal(ers)
				o.Error = string(e)
				duplicated = append(duplicated, o)
				break
			}
			_, err := dbc.GetRangeTimeCfgByID(o.ObjectID)
			if err == nil {
				o.Error = fmt.Sprintf("Duplicated object %s in the database", o.ObjectID)
				duplicated = append(duplicated, o)
			}
		case "ifxservercfg":
			data := config.IfxServerCfg{}
			json.Unmarshal(raw, &data)
			ers := binding.RawValidate(data)
			if ers.Len() > 0 {
				e, _ := json.Marshal(ers)
				o.Error = string(e)
				duplicated = append(duplicated, o)
				break
			}
			_, err := dbc.GetIfxServerCfgByID(o.ObjectID)
			if err == nil {
				o.Error = fmt.Sprintf("Duplicated object %s in the database", o.ObjectID)
				duplicated = append(duplicated, o)
			}
		case "kapacitorcfg":
			data := config.KapacitorCfg{}
			json.Unmarshal(raw, &data)
			ers := binding.RawValidate(data)
			if ers.Len() > 0 {
				e, _ := json.Marshal(ers)
				o.Error = string(e)
				duplicated = append(duplicated, o)
				break
			}
			_, err := dbc.GetKapacitorCfgByID(o.ObjectID)
			if err == nil {
				o.Error = fmt.Sprintf("Duplicated object %s in the database", o.ObjectID)
				duplicated = append(duplicated, o)
			}
		case "productcfg":
			data := config.ProductCfg{}
			json.Unmarshal(raw, &data)
			ers := binding.RawValidate(data)
			if ers.Len() > 0 {
				e, _ := json.Marshal(ers)
				o.Error = string(e)
				duplicated = append(duplicated, o)
				break
			}
			_, err := dbc.GetProductCfgByID(o.ObjectID)
			if err == nil {
				o.Error = fmt.Sprintf("Duplicated object %s in the database", o.ObjectID)
				duplicated = append(duplicated, o)
			}
		case "productgroupcfg":
			data := config.ProductGroupCfg{}
			json.Unmarshal(raw, &data)
			ers := binding.RawValidate(data)
			if ers.Len() > 0 {
				e, _ := json.Marshal(ers)
				o.Error = string(e)
				duplicated = append(duplicated, o)
				break
			}
			_, err := dbc.GetProductGroupCfgByID(o.ObjectID)
			if err == nil {
				o.Error = fmt.Sprintf("Duplicated object %s in the database", o.ObjectID)
				duplicated = append(duplicated, o)
			}
		case "outhttpcfg":
			data := config.OutHTTPCfg{}
			json.Unmarshal(raw, &data)
			ers := binding.RawValidate(data)
			if ers.Len() > 0 {
				e, _ := json.Marshal(ers)
				o.Error = string(e)
				duplicated = append(duplicated, o)
				break
			}
			_, err := dbc.GetOutHTTPCfgByID(o.ObjectID)
			if err == nil {
				o.Error = fmt.Sprintf("Duplicated object %s in the database", o.ObjectID)
				duplicated = append(duplicated, o)
			}
		case "alertcfg":
			data := config.AlertIDCfg{}
			json.Unmarshal(raw, &data)
			ers := binding.RawValidate(data)
			if ers.Len() > 0 {
				e, _ := json.Marshal(ers)
				o.Error = string(e)
				duplicated = append(duplicated, o)
				break
			}
			_, err := dbc.GetAlertIDCfgByID(o.ObjectID)
			if err == nil {
				o.Error = fmt.Sprintf("Duplicated object %s in the database", o.ObjectID)
				duplicated = append(duplicated, o)
			}
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
		case "rangetimecfg":
			log.Debugf("Importing rangetimecfg : %+v", o.ObjectCfg)
			data := config.RangeTimeCfg{}
			json.Unmarshal(raw, &data)
			var err error
			_, err = dbc.GetRangeTimeCfgByID(o.ObjectID)
			if err == nil { //value exist already in the database
				if overwrite == true {
					_, err2 := dbc.UpdateRangeTimeCfg(o.ObjectID, &data)
					if err2 != nil {
						return fmt.Errorf("Error on overwrite object [%s] %s : %s", o.ObjectTypeID, o.ObjectID, err2)
					}
					break
				}
			}
			if autorename == true {
				data.ID = data.ID + suffix
			}
			_, err = dbc.AddRangeTimeCfg(&data)
			if err != nil {
				return err
			}
		case "ifxservercfg":
			log.Debugf("Importing ifxservercfg : %+v", o.ObjectCfg)
			data := config.IfxServerCfg{}
			json.Unmarshal(raw, &data)
			var err error
			_, err = dbc.GetIfxServerCfgByID(o.ObjectID)
			if err == nil { //value exist already in the database
				if overwrite == true {
					_, err2 := dbc.UpdateIfxServerCfg(o.ObjectID, &data)
					if err2 != nil {
						return fmt.Errorf("Error on overwrite object [%s] %s : %s", o.ObjectTypeID, o.ObjectID, err2)
					}
					break
				}
			}
			if autorename == true {
				data.ID = data.ID + suffix
			}
			_, err = dbc.AddIfxServerCfg(&data)
			if err != nil {
				return err
			}
		case "kapacitorcfg":
			log.Debugf("Importing kapacitorcfg : %+v", o.ObjectCfg)
			data := config.KapacitorCfg{}
			json.Unmarshal(raw, &data)
			var err error
			_, err = dbc.GetKapacitorCfgByID(o.ObjectID)
			if err == nil { //value exist already in the database
				if overwrite == true {
					_, err2 := dbc.UpdateKapacitorCfg(o.ObjectID, &data)
					if err2 != nil {
						return fmt.Errorf("Error on overwrite object [%s] %s : %s", o.ObjectTypeID, o.ObjectID, err2)
					}
					break
				}
			}
			if autorename == true {
				data.ID = data.ID + suffix
			}
			_, err = dbc.AddKapacitorCfg(&data)
			if err != nil {
				return err
			}
		case "productcfg":
			log.Debugf("Importing productcfg : %+v", o.ObjectCfg)
			data := config.ProductCfg{}
			json.Unmarshal(raw, &data)
			var err error
			_, err = dbc.GetProductCfgByID(o.ObjectID)
			if err == nil { //value exist already in the database
				if overwrite == true {
					_, err2 := dbc.UpdateProductCfg(o.ObjectID, &data)
					if err2 != nil {
						return fmt.Errorf("Error on overwrite object [%s] %s : %s", o.ObjectTypeID, o.ObjectID, err2)
					}
					break
				}
			}
			if autorename == true {
				data.ID = data.ID + suffix
			}
			_, err = dbc.AddProductCfg(&data)
			if err != nil {
				return err
			}
		case "productgroupcfg":
			log.Debugf("Importing productgroupcfg : %+v", o.ObjectCfg)
			data := config.ProductGroupCfg{}
			json.Unmarshal(raw, &data)
			var err error
			_, err = dbc.GetProductGroupCfgByID(o.ObjectID)
			if err == nil { //value exist already in the database
				if overwrite == true {
					_, err2 := dbc.UpdateProductGroupCfg(o.ObjectID, &data)
					if err2 != nil {
						return fmt.Errorf("Error on overwrite object [%s] %s : %s", o.ObjectTypeID, o.ObjectID, err2)
					}
					break
				}
			}
			if autorename == true {
				data.ID = data.ID + suffix
			}
			_, err = dbc.AddProductGroupCfg(&data)
			if err != nil {
				return err
			}
		case "outhttpcfg":
			log.Debugf("Importing outhttpcfg : %+v", o.ObjectCfg)
			data := config.OutHTTPCfg{}
			json.Unmarshal(raw, &data)
			var err error
			_, err = dbc.GetOutHTTPCfgByID(o.ObjectID)
			if err == nil { //value exist already in the database
				if overwrite == true {
					_, err2 := dbc.UpdateOutHTTPCfg(o.ObjectID, &data)
					if err2 != nil {
						return fmt.Errorf("Error on overwrite object [%s] %s : %s", o.ObjectTypeID, o.ObjectID, err2)
					}
					break
				}
			}
			if autorename == true {
				data.ID = data.ID + suffix
			}
			_, err = dbc.AddOutHTTPCfg(&data)
			if err != nil {
				return err
			}
		case "alertcfg":
			log.Debugf("Importing alertcfg : %+v", o.ObjectCfg)
			data := config.AlertIDCfg{}
			json.Unmarshal(raw, &data)
			var err error
			_, err = dbc.GetAlertIDCfgByID(o.ObjectID)
			if err == nil { //value exist already in the database
				if overwrite == true {
					_, err2 := dbc.UpdateAlertIDCfg(o.ObjectID, &data)
					if err2 != nil {
						return fmt.Errorf("Error on overwrite object [%s] %s : %s", o.ObjectTypeID, o.ObjectID, err2)
					}
					break
				}
			}
			if autorename == true {
				data.ID = data.ID + suffix
			}
			_, err = dbc.AddAlertIDCfg(&data)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unknown type object type %s ", o.ObjectTypeID)
		}
	}
	return nil
}
