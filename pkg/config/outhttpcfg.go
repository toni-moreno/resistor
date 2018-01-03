package config

import (
	"fmt"
)

/***************************
	OutHTTPCfg DB backends
	-GetOutHTTPCfgCfgByID(struct)
	-GetOutHTTPCfgMap (map - for interna config use
	-GetOutHTTPCfgArray(Array - for web ui use )
	-AddOutHTTPCfg
	-DelOutHTTPCfg
	-UpdateOutHTTPCfg
  -GetOutHTTPCfgAffectOnDel
***********************************/

/*GetOutHTTPCfgByID get device data by id*/
func (dbc *DatabaseCfg) GetOutHTTPCfgByID(id string) (OutHTTPCfg, error) {
	cfgarray, err := dbc.GetOutHTTPCfgArray("id='" + id + "'")
	if err != nil {
		return OutHTTPCfg{}, err
	}
	if len(cfgarray) > 1 {
		return OutHTTPCfg{}, fmt.Errorf("Error %d results on get OutHTTPCfg by id %s", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return OutHTTPCfg{}, fmt.Errorf("Error no values have been returned with this id %s in the influx config table", id)
	}
	return *cfgarray[0], nil
}

/*GetOutHTTPCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetOutHTTPCfgMap(filter string) (map[string]*OutHTTPCfg, error) {
	cfgarray, err := dbc.GetOutHTTPCfgArray(filter)
	cfgmap := make(map[string]*OutHTTPCfg)
	for _, val := range cfgarray {
		cfgmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetOutHTTPCfgArray generate an array of devices with all its information */
func (dbc *DatabaseCfg) GetOutHTTPCfgArray(filter string) ([]*OutHTTPCfg, error) {
	var err error
	var devices []*OutHTTPCfg
	//Get Only data for selected devices
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&devices); err != nil {
			log.Warnf("Fail to get OutHTTPCfg  data filteter with %s : %v\n", filter, err)
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

/*AddOutHTTPCfg for adding new devices*/
func (dbc *DatabaseCfg) AddOutHTTPCfg(dev *OutHTTPCfg) (int64, error) {
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
	log.Infof("Added new  HTTP Output backend Successfully with id %s ", dev.ID)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelOutHTTPCfg for deleting influx databases from ID*/
func (dbc *DatabaseCfg) DelOutHTTPCfg(id string) (int64, error) {
	var affecteddev, affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affecteddev, err = session.Where("kapacitorid='" + id + "'").Cols("kapacitorid").Update(&AlertIDCfg{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on Delete OutHTTPCfg with id: %s, error: %s", id, err)
	}

	affected, err = session.Where("id='" + id + "'").Delete(&OutHTTPCfg{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Deleted Successfully HTTP Output with ID %s [ %d Devices Affected  ]", id, affecteddev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*UpdateOutHTTPCfg for adding new influxdb*/
func (dbc *DatabaseCfg) UpdateOutHTTPCfg(id string, dev *OutHTTPCfg) (int64, error) {
	var affecteddev, affected int64
	var err error
	session := dbc.x.NewSession()
	defer session.Close()
	if id != dev.ID { //ID has been changed
		affecteddev, err = session.Where("kapacitorid='" + id + "'").Cols("kapacitorid").Update(&AlertIDCfg{KapacitorID: dev.ID})
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

	log.Infof("Updated HTTP OutPut Successfully with id %s and data:%+v, affected", id, dev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*GetOutHTTPCfgAffectOnDel for deleting devices from ID*/
func (dbc *DatabaseCfg) GetOutHTTPCfgAffectOnDel(id string) ([]*DbObjAction, error) {
	var devices []*AlertHTTPOutRel
	var obj []*DbObjAction
	if err := dbc.x.Where("http_out_id='" + id + "'").Find(&devices); err != nil {
		log.Warnf("Error on Get Out HTTP out id %d for devices , error: %s", id, err)
		return nil, err
	}

	for _, val := range devices {
		obj = append(obj, &DbObjAction{
			Type:     "alertidcfg",
			TypeDesc: "alertID's ",
			ObID:     val.AlertID,
			Action:   "Change alert to Other Kapacitor alert",
		})

	}
	return obj, nil
}
