package config

import (
	"fmt"
)

/***************************
	IfxMeasurementCfg DB backends
	-GetIfxMeasurementCfgCfgByID(struct)
	-GetIfxMeasurementCfgMap (map - for interna config use
	-GetIfxMeasurementCfgArray(Array - for web ui use )
	-AddIfxMeasurementCfg
	-DelIfxMeasurementCfg
	-UpdateIfxMeasurementCfg
  -GetIfxMeasurementCfgAffectOnDel
***********************************/

/*GetIfxMeasurementCfgByID get device data by id*/
func (dbc *DatabaseCfg) GetIfxMeasurementCfgByID(id string) (IfxMeasurementCfg, error) {
	cfgarray, err := dbc.GetIfxMeasurementCfgArray("id='" + id + "'")
	if err != nil {
		return IfxMeasurementCfg{}, err
	}
	if len(cfgarray) > 1 {
		return IfxMeasurementCfg{}, fmt.Errorf("Error %d results on get IfxMeasurementCfg by id %s", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return IfxMeasurementCfg{}, fmt.Errorf("Error no values have been returned with this id %s in the influx config table", id)
	}
	return *cfgarray[0], nil
}

/*GetIfxMeasurementCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetIfxMeasurementCfgMap(filter string) (map[string]*IfxMeasurementCfg, error) {
	cfgarray, err := dbc.GetIfxMeasurementCfgArray(filter)
	cfgmap := make(map[string]*IfxMeasurementCfg)
	for _, val := range cfgarray {
		cfgmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetIfxMeasurementCfgArray generate an array of devices with all its information */
func (dbc *DatabaseCfg) GetIfxMeasurementCfgArray(filter string) ([]*IfxMeasurementCfg, error) {
	var err error
	var devices []*IfxMeasurementCfg
	//Get Only data for selected devices
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&devices); err != nil {
			log.Warnf("Fail to get IfxMeasurementCfg  data filteter with %s : %v\n", filter, err)
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

/*AddIfxMeasurementCfg for adding new devices*/
func (dbc *DatabaseCfg) AddIfxMeasurementCfg(dev IfxMeasurementCfg) (int64, error) {
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

/*DelIfxMeasurementCfg for deleting influx databases from ID*/
func (dbc *DatabaseCfg) DelIfxMeasurementCfg(id string) (int64, error) {
	var affecteddev, affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()
	/*
		affecteddev, err = session.Where("kapacitorid='" + id + "'").Cols("kapacitorid").Update(&AlertIDCfg{})
		if err != nil {
			session.Rollback()
			return 0, fmt.Errorf("Error on Delete IfxMeasurementCfg with id: %s, error: %s", id, err)
		}
	*/
	affected, err = session.Where("id='" + id + "'").Delete(&IfxMeasurementCfg{})
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

/*UpdateIfxMeasurementCfg for adding new influxdb*/
func (dbc *DatabaseCfg) UpdateIfxMeasurementCfg(id string, dev IfxMeasurementCfg) (int64, error) {
	var affecteddev, affected int64
	var err error
	session := dbc.x.NewSession()
	defer session.Close()
	if id != dev.ID { //ID has been changed
		/*
			affecteddev, err = session.Where("kapacitorid='" + id + "'").Cols("kapacitorid").Update(&AlertIDCfg{KapacitorID: dev.ID})
			if err != nil {
				session.Rollback()
				return 0, fmt.Errorf("Error on Update InfluxConfig on update id(old)  %s with (new): %s, error: %s", id, dev.ID, err)
			}*/
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

	log.Infof("Updated KapacitorID Config Successfully with id %s and data:%+v, affected", id, dev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*GetIfxMeasurementCfgAffectOnDel for deleting devices from ID*/
func (dbc *DatabaseCfg) GetIfxMeasurementCfgAffectOnDel(id string) ([]*DbObjAction, error) {
	var devices []*AlertIDCfg
	var obj []*DbObjAction
	if err := dbc.x.Where("kapacitorid='" + id + "'").Find(&devices); err != nil {
		log.Warnf("Error on Get Outout db id %d for devices , error: %s", id, err)
		return nil, err
	}

	for _, val := range devices {
		obj = append(obj, &DbObjAction{
			Type:     "alertidcfg",
			TypeDesc: "",
			ObID:     val.ID,
			Action:   "Change alert to Other Kapacitor alert",
		})

	}
	return obj, nil
}
