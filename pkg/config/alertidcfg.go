package config

import (
	"fmt"
)

/***************************
	AlertIDCfg DB backends
	-GetAlertIDCfgCfgByID(struct)
	-GetAlertIDCfgMap (map - for interna config use
	-GetAlertIDCfgArray(Array - for web ui use )
	-AddAlertIDCfg
	-DelAlertIDCfg
	-UpdateAlertIDCfg
  -GetAlertIDCfgAffectOnDel
***********************************/

/*GetAlertIDCfgByID get device data by id*/
func (dbc *DatabaseCfg) GetAlertIDCfgByID(id string) (AlertIDCfg, error) {
	cfgarray, err := dbc.GetAlertIDCfgArray("id='" + id + "'")
	if err != nil {
		return AlertIDCfg{}, err
	}
	if len(cfgarray) > 1 {
		return AlertIDCfg{}, fmt.Errorf("Error %d results on get AlertIDCfgArray by id %s", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return AlertIDCfg{}, fmt.Errorf("Error no values have been returned with this id %s in the influx config table", id)
	}
	return *cfgarray[0], nil
}

/*GetAlertIDCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetAlertIDCfgMap(filter string) (map[string]*AlertIDCfg, error) {
	cfgarray, err := dbc.GetAlertIDCfgArray(filter)
	cfgmap := make(map[string]*AlertIDCfg)
	for _, val := range cfgarray {
		cfgmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetAlertIDCfgArray generate an array of devices with all its information */
func (dbc *DatabaseCfg) GetAlertIDCfgArray(filter string) ([]*AlertIDCfg, error) {
	var err error
	var devices []*AlertIDCfg
	//Get Only data for selected devices
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&devices); err != nil {
			log.Warnf("Fail to get AlertIDCfg  data filteter with %s : %v\n", filter, err)
			return nil, err
		}
	} else {
		if err = dbc.x.Find(&devices); err != nil {
			log.Warnf("Fail to get influxcfg   data: %v\n", err)
			return nil, err
		}
	}

	//Asign HTTP Outs to AlertID's
	var alerthttp []*AlertHTTPOutRel
	if err = dbc.x.Find(&alerthttp); err != nil {
		log.Warnf("Fail to get Output Alert and its HTTP output relationship data: %v\n", err)
		return devices, err
	}

	//Load Measurements and metrics relationship
	//We assign field metric ID to each measurement
	for _, v := range devices {
		for _, r := range alerthttp {
			if r.AlertID == v.ID {
				v.OutHTTP = append(v.OutHTTP, r.HTTPOutID)
			}
		}
	}

	return devices, nil
}

/*AddAlertIDCfg for adding new devices*/
func (dbc *DatabaseCfg) AddAlertIDCfg(dev AlertIDCfg) (int64, error) {
	var err error
	var affected, newo int64
	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Insert(dev)
	if err != nil {
		session.Rollback()
		return 0, err
	}
	//save new alert ot outhttp relations
	for _, o := range dev.OutHTTP {

		ostruct := AlertHTTPOutRel{
			AlertID:   dev.ID,
			HTTPOutID: o,
		}
		newo, err = session.Insert(&ostruct)
		if err != nil {
			session.Rollback()
			return 0, err
		}
	}

	//no other relation
	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Added new Kapacitor backend Successfully with id %s [ %d HTTP Output] ", dev.ID, newo)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelAlertIDCfg for deleting influx databases from ID*/
func (dbc *DatabaseCfg) DelAlertIDCfg(id string) (int64, error) {
	var affecteddev, affectedouts, affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affecteddev, err = session.Where("kapacitorid='" + id + "'").Cols("kapacitorid").Update(&AlertIDCfg{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on Delete Alert with id on delete AlertIDCfg with id: %s, error: %s", id, err)
	}

	//first deleting references in AlertHTTPOutRel
	affectedouts, err = session.Where("alert_id='" + id + "'").Delete(&AlertHTTPOutRel{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on Delete Device with id on delete AlertHTTPOutRel with id: %s, error: %s", id, err)
	}

	affected, err = session.Where("id='" + id + "'").Delete(&AlertIDCfg{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Deleted Successfully Alert with ID %s [ %d Affected alerts   ] [%d affected Outs]", id, affecteddev, affectedouts)
	dbc.addChanges(affected + affecteddev + affectedouts)
	return affected, nil
}

/*UpdateAlertIDCfg for adding new influxdb*/
func (dbc *DatabaseCfg) UpdateAlertIDCfg(id string, dev AlertIDCfg) (int64, error) {
	var affecteddev, affected int64
	var err error
	session := dbc.x.NewSession()
	defer session.Close()
	//first deleting references in AlertHTTPOutRel
	_, err = session.Where("alert_id='" + id + "'").Delete(&AlertHTTPOutRel{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on Delete Device with id on delete AlertHTTPOutRel with id: %s, error: %s", id, err)
	}
	//save again aññ alert ot outhttp relations
	for _, o := range dev.OutHTTP {

		ostruct := AlertHTTPOutRel{
			AlertID:   dev.ID,
			HTTPOutID: o,
		}
		_, err = session.Insert(&ostruct)
		if err != nil {
			session.Rollback()
			return 0, err
		}
	}

	affected, err = session.Where("id='" + id + "'").UseBool().AllCols().Update(dev)
	if err != nil {
		session.Rollback()
		return 0, err
	}
	err = session.Commit()
	if err != nil {
		return 0, err
	}

	log.Infof("Updated KapacitorID Config Successfully with id %s and data:%+v, affected", id, dev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*GetAlertIDCfgAffectOnDel for deleting devices from ID*/
func (dbc *DatabaseCfg) GetAlertIDCfgAffectOnDel(id string) ([]*DbObjAction, error) {
	//var devices []*AlertIDCfg
	var obj []*DbObjAction

	return obj, nil
}
