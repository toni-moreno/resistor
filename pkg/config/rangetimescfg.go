package config

import (
	"fmt"
)

/***************************
	RangeTimeCfg DB backends
	-GetRangeTimeCfgCfgByID(struct)
	-GetRangeTimeCfgMap (map - for interna config use
	-GetRangeTimeCfgArray(Array - for web ui use )
	-AddRangeTimeCfg
	-DelRangeTimeCfg
	-UpdateRangeTimeCfg
  -GetRangeTimeCfgAffectOnDel
***********************************/

/*GetRangeTimeCfgByID get device data by id*/
func (dbc *DatabaseCfg) GetRangeTimeCfgByID(id string) (RangeTimeCfg, error) {
	cfgarray, err := dbc.GetRangeTimeCfgArray("id='" + id + "'")
	if err != nil {
		return RangeTimeCfg{}, err
	}
	if len(cfgarray) > 1 {
		return RangeTimeCfg{}, fmt.Errorf("Error %d results on get RangeTimeCfgArray by id %s", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return RangeTimeCfg{}, fmt.Errorf("Error no values have been returned with this id %s in the influx config table", id)
	}
	return *cfgarray[0], nil
}

/*GetRangeTimeCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetRangeTimeCfgMap(filter string) (map[string]*RangeTimeCfg, error) {
	cfgarray, err := dbc.GetRangeTimeCfgArray(filter)
	cfgmap := make(map[string]*RangeTimeCfg)
	for _, val := range cfgarray {
		cfgmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetRangeTimeCfgArray generate an array of devices with all its information */
func (dbc *DatabaseCfg) GetRangeTimeCfgArray(filter string) ([]*RangeTimeCfg, error) {
	var err error
	var devices []*RangeTimeCfg
	//Get Only data for selected devices
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&devices); err != nil {
			log.Warnf("Fail to get RangeTimeCfg  data filteter with %s : %v\n", filter, err)
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

/*AddRangeTimeCfg for adding new devices*/
func (dbc *DatabaseCfg) AddRangeTimeCfg(dev RangeTimeCfg) (int64, error) {
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

/*DelRangeTimeCfg for deleting influx databases from ID*/
func (dbc *DatabaseCfg) DelRangeTimeCfg(id string) (int64, error) {
	var affecteddev, affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affecteddev, err = session.Where("th_crit_rangetime_id='" + id + "'").Cols("th_crit_rangetime_id").Update(&AlertIdCfg{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on Delete Alert with id on delete RangeTimeCfg with id: %s, error: %s", id, err)
	}
	affecteddev, err = session.Where("th_warn_rangetime_id='" + id + "'").Cols("th_warn_rangetime_id").Update(&AlertIdCfg{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on Delete Alert with id on delete RangeTimeCfg with id: %s, error: %s", id, err)
	}
	affecteddev, err = session.Where("th_info_rangetime_id='" + id + "'").Cols("th_info_rangetime_id").Update(&AlertIdCfg{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on Delete Alert with id on delete RangeTimeCfg with id: %s, error: %s", id, err)
	}

	affected, err = session.Where("id='" + id + "'").Delete(&RangeTimeCfg{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Deleted Successfully influx db with ID %s [ %d Devices Affected  ]", id, affecteddev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*UpdateRangeTimeCfg for adding new influxdb*/
func (dbc *DatabaseCfg) UpdateRangeTimeCfg(id string, dev RangeTimeCfg) (int64, error) {
	var affecteddev, affected int64
	var err error
	session := dbc.x.NewSession()
	defer session.Close()
	if id != dev.ID { //ID has been changed
		affecteddev, err = session.Where("th_crit_rangetime_id='" + id + "'").Cols("th_crit_rangetime_id").Update(&AlertIdCfg{ThCritRangeTimeID: dev.ID})
		if err != nil {
			session.Rollback()
			return 0, fmt.Errorf("Error on Update InfluxConfig on update id(old)  %s with (new): %s, error: %s", id, dev.ID, err)
		}
		affecteddev, err = session.Where("th_warn_rangetime_id='" + id + "'").Cols("th_warn_rangetime_id").Update(&AlertIdCfg{ThWarnRangeTimeID: dev.ID})
		if err != nil {
			session.Rollback()
			return 0, fmt.Errorf("Error on Update InfluxConfig on update id(old)  %s with (new): %s, error: %s", id, dev.ID, err)
		}
		affecteddev, err = session.Where("th_info_rangetime_id='" + id + "'").Cols("th_info_rangetime_id").Update(&AlertIdCfg{ThInfoRangeTimeID: dev.ID})
		if err != nil {
			session.Rollback()
			return 0, fmt.Errorf("Error on Update InfluxConfig on update id(old)  %s with (new): %s, error: %s", id, dev.ID, err)
		}
		log.Infof("Updated Influx Config to %s devices ", affecteddev)
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

	log.Infof("Updated RangeTime Config Successfully with id %s and data:%+v, affected", id, dev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*GetRangeTimeCfgAffectOnDel for deleting devices from ID*/
func (dbc *DatabaseCfg) GetRangeTimeCfgAffectOnDel(id string) ([]*DbObjAction, error) {
	var devices []*RangeTimeCfg
	var obj []*DbObjAction
	if err := dbc.x.Where("th_crit_rangetime_id='" + id + "'").Find(&devices); err != nil {
		log.Warnf("Error on Get Outout db id %d for devices , error: %s", id, err)
		return nil, err
	}

	for _, val := range devices {
		obj = append(obj, &DbObjAction{
			Type:     "rangetime",
			TypeDesc: "Crit Range Time",
			ObID:     val.ID,
			Action:   "Reset InfluxDB Server from SNMPDevice to 'default' InfluxDB Server",
		})

	}
	return obj, nil
}
