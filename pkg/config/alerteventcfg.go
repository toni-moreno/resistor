package config

import (
	"fmt"
	"strconv"
)

/***************************
	AlertEventCfg Alert events
	-GetAlertEventCfgCfgByID(struct)
	-GetAlertEventCfgMap (map - for interna config use
	-GetAlertEventCfgArray(Array - for web ui use )
***********************************/

/*GetAlertEventCfgByID get device data by id*/
func (dbc *DatabaseCfg) GetAlertEventCfgByID(uid int64) (AlertEventCfg, error) {
	cfgarray, err := dbc.GetAlertEventCfgArray("uid='" + strconv.FormatInt(uid, 10) + "'")
	if err != nil {
		return AlertEventCfg{}, err
	}
	if len(cfgarray) > 1 {
		return AlertEventCfg{}, fmt.Errorf("Error %d results on get AlertEventCfgArray by uid %d", len(cfgarray), uid)
	}
	if len(cfgarray) == 0 {
		return AlertEventCfg{}, fmt.Errorf("Error no values have been returned with this id %d in the influx config table", uid)
	}
	return *cfgarray[0], nil
}

/*GetAlertEventCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetAlertEventCfgMap(filter string) (map[int64]*AlertEventCfg, error) {
	cfgarray, err := dbc.GetAlertEventCfgArray(filter)
	cfgmap := make(map[int64]*AlertEventCfg)
	for _, val := range cfgarray {
		cfgmap[val.UID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetAlertEventCfgArray generate an array of devices with all its information */
func (dbc *DatabaseCfg) GetAlertEventCfgArray(filter string) ([]*AlertEventCfg, error) {
	var err error
	var devices []*AlertEventCfg
	//Get Only data for selected devices
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&devices); err != nil {
			log.Warnf("Fail to get AlertEventCfg  data filteter with %s : %v\n", filter, err)
			return nil, err
		}
	} else {
		if err = dbc.x.Find(&devices); err != nil {
			log.Warnf("Fail to get influxcfg   data: %v\n", err)
			return nil, err
		}
	}
	return devices, nil
}

/*AddAlertEventCfg for adding new devices*/
func (dbc *DatabaseCfg) AddAlertEventCfg(dev *AlertEventCfg) (int64, error) {
	var err error
	var affected int64
	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Insert(dev)
	if err != nil {
		session.Rollback()
		return 0, err
	}
	//no other relation
	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Added new Alert event Successfully with id %s ", dev.ID)
	dbc.addChanges(affected)
	return affected, nil
}
