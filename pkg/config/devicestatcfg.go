package config

import (
	"fmt"
	"strconv"
)

/***************************
	DeviceStatCfg DB backends
	-GetDeviceStatCfgByID(struct)
	-GetDeviceStatCfgMap (map - for interna config use
	-GetDeviceStatCfgArray(Array - for web ui use )
	-AddDeviceStatCfg
	-DelDeviceStatCfg
	-UpdateDeviceStatCfg
  -GetDeviceStatCfgAffectOnDel
***********************************/

/*GetDeviceStatCfgByID get device data by id*/
func (dbc *DatabaseCfg) GetDeviceStatCfgByID(id int64) (DeviceStatCfg, error) {
	cfgarray, err := dbc.GetDeviceStatCfgArray("id='" + strconv.FormatInt(id, 10) + "'")
	if err != nil {
		return DeviceStatCfg{}, err
	}
	if len(cfgarray) > 1 {
		return DeviceStatCfg{}, fmt.Errorf("Error %d results on get DeviceStatCfgArray by id %d", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return DeviceStatCfg{}, fmt.Errorf("Error no values have been returned with this id %d in the influx config table", id)
	}
	return *cfgarray[0], nil
}

/*GetDeviceStatCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetDeviceStatCfgMap(filter string) (map[int64]*DeviceStatCfg, error) {
	cfgarray, err := dbc.GetDeviceStatCfgArray(filter)
	cfgmap := make(map[int64]*DeviceStatCfg)
	for _, val := range cfgarray {
		cfgmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetDeviceStatCfgArray generate an array of devices with all its information */
func (dbc *DatabaseCfg) GetDeviceStatCfgArray(filter string) ([]*DeviceStatCfg, error) {
	var err error
	var devices []*DeviceStatCfg
	//Get Only data for selected devices
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&devices); err != nil {
			log.Warnf("Fail to get DeviceStatCfg  data filteter with %s : %v\n", filter, err)
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

/*AddDeviceStatCfg for adding new devices*/
func (dbc *DatabaseCfg) AddDeviceStatCfg(dev DeviceStatCfg) (int64, error) {
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
	log.Infof("Added new Kapacitor backend Successfully with id %s ", dev.ID)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelDeviceStatCfg for deleting influx databases from ID*/
func (dbc *DatabaseCfg) DelDeviceStatCfg(id int64) (int64, error) {
	var affecteddev, affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Where("id='" + strconv.FormatInt(id, 10) + "'").Delete(&DeviceStatCfg{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Deleted Successfully influx db with ID %d [ %d Devices Affected  ]", id, affecteddev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*UpdateDeviceStatCfg for adding new influxdb*/
func (dbc *DatabaseCfg) UpdateDeviceStatCfg(id int64, dev DeviceStatCfg) (int64, error) {
	var affecteddev, affected int64
	var err error
	session := dbc.x.NewSession()
	defer session.Close()
	/*if id != dev.ID { //ID has been changed
		affecteddev, err = session.Where("id='" + id + "'").Cols("id").Update(&DeviceStatCfg{ID: dev.ID})
		if err != nil {
			session.Rollback()
			return 0, fmt.Errorf("Error on Update InfluxConfig on update id(old)  %d with (new): %d, error: %s", id, dev.ID, err)
		}
		log.Infof("Updated Influx Config to %s devices ", affecteddev)
	}*/

	affected, err = session.Where("id='" + strconv.FormatInt(id, 10) + "'").UseBool().AllCols().Update(dev)
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

/*GetDeviceStatCfgAffectOnDel for deleting devices from ID*/
func (dbc *DatabaseCfg) GetDeviceStatCfgAffectOnDel(id string) ([]*DbObjAction, error) {
	//var devices []*DeviceStatCfg
	var obj []*DbObjAction
	/*if err := dbc.x.Where("outdb='" + id + "'").Find(&devices); err != nil {
		log.Warnf("Error on Get Outout db id %d for devices , error: %s", id, err)
		return nil, err
	}

	for _, val := range devices {
		obj = append(obj, &DbObjAction{
			Type:     "snmpdevicecfg",
			TypeDesc: "SNMP Devices",
			ObID:     val.ID,
			Action:   "Reset InfluxDB Server from SNMPDevice to 'default' InfluxDB Server",
		})

	}*/
	return obj, nil
}
