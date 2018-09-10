package config

import (
	"fmt"
)

/***************************
	EndpointCfg DB backends
	-GetEndpointCfgCfgByID(struct)
	-GetEndpointCfgMap (map - for interna config use
	-GetEndpointCfgArray(Array - for web ui use )
	-AddEndpointCfg
	-DelEndpointCfg
	-UpdateEndpointCfg
  -GetEndpointCfgAffectOnDel
***********************************/

/*GetEndpointCfgByID get device data by id*/
func (dbc *DatabaseCfg) GetEndpointCfgByID(id string) (EndpointCfg, error) {
	cfgarray, err := dbc.GetEndpointCfgArray("id='" + id + "'")
	if err != nil {
		return EndpointCfg{}, err
	}
	if len(cfgarray) > 1 {
		return EndpointCfg{}, fmt.Errorf("Error %d results on get EndpointCfg by id %s", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return EndpointCfg{}, fmt.Errorf("Error no values have been returned with this id %s in the influx config table", id)
	}
	return *cfgarray[0], nil
}

/*GetEndpointCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetEndpointCfgMap(filter string) (map[string]*EndpointCfg, error) {
	cfgarray, err := dbc.GetEndpointCfgArray(filter)
	cfgmap := make(map[string]*EndpointCfg)
	for _, val := range cfgarray {
		cfgmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetEndpointCfgArray generate an array of devices with all its information */
func (dbc *DatabaseCfg) GetEndpointCfgArray(filter string) ([]*EndpointCfg, error) {
	var err error
	var devices []*EndpointCfg
	//Get Only data for selected devices
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&devices); err != nil {
			log.Warnf("Fail to get EndpointCfg  data filteter with %s : %v\n", filter, err)
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

/*AddEndpointCfg for adding new devices*/
func (dbc *DatabaseCfg) AddEndpointCfg(dev *EndpointCfg) (int64, error) {
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
	log.Infof("Added new  Endpoint backend Successfully with id %s ", dev.ID)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelEndpointCfg for deleting influx databases from ID*/
func (dbc *DatabaseCfg) DelEndpointCfg(id string) (int64, error) {
	var affecteddev, affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affecteddev, err = session.Where("kapacitorid='" + id + "'").Cols("kapacitorid").Update(&AlertIDCfg{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on Delete EndpointCfg with id: %s, error: %s", id, err)
	}

	affected, err = session.Where("id='" + id + "'").Delete(&EndpointCfg{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Deleted Successfully Endpoint with ID %s [ %d Devices Affected  ]", id, affecteddev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*UpdateEndpointCfg for adding new influxdb*/
func (dbc *DatabaseCfg) UpdateEndpointCfg(id string, dev *EndpointCfg) (int64, error) {
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

	log.Infof("Updated Endpoint Successfully with id %s and data:%+v, affected", id, dev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*GetEndpointCfgAffectOnDel for deleting devices from ID*/
func (dbc *DatabaseCfg) GetEndpointCfgAffectOnDel(id string) ([]*DbObjAction, error) {
	var devices []*AlertEndpointRel
	var obj []*DbObjAction
	if err := dbc.x.Where("endpoint_id='" + id + "'").Find(&devices); err != nil {
		log.Warnf("Error on Get Out Endpoint id %d for devices , error: %s", id, err)
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
