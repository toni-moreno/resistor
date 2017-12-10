package config

import (
	"fmt"
)

/***************************
	ProductGroupCfg DB backends
	-GetProductGroupCfgCfgByID(struct)
	-GetProductGroupCfgMap (map - for interna config use
	-GetProductGroupCfgArray(Array - for web ui use )
	-AddProductGroupCfg
	-DelProductGroupCfg
	-UpdateProductGroupCfg
  -GetProductGroupCfgAffectOnDel
***********************************/

/*GetProductGroupCfgByID get device data by id*/
func (dbc *DatabaseCfg) GetProductGroupCfgByID(id string) (ProductGroupCfg, error) {
	cfgarray, err := dbc.GetProductGroupCfgArray("id='" + id + "'")
	if err != nil {
		return ProductGroupCfg{}, err
	}
	if len(cfgarray) > 1 {
		return ProductGroupCfg{}, fmt.Errorf("Error %d results on get ProductGroupCfg by id %s", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return ProductGroupCfg{}, fmt.Errorf("Error no values have been returned with this id %s in the influx config table", id)
	}
	return *cfgarray[0], nil
}

/*GetProductGroupCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetProductGroupCfgMap(filter string) (map[string]*ProductGroupCfg, error) {
	cfgarray, err := dbc.GetProductGroupCfgArray(filter)
	cfgmap := make(map[string]*ProductGroupCfg)
	for _, val := range cfgarray {
		cfgmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetProductGroupCfgArray generate an array of devices with all its information */
func (dbc *DatabaseCfg) GetProductGroupCfgArray(filter string) ([]*ProductGroupCfg, error) {
	var err error
	var devices []*ProductGroupCfg
	//Get Only data for selected devices
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&devices); err != nil {
			log.Warnf("Fail to get ProductGroupCfg  data filteter with %s : %v\n", filter, err)
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

/*AddProductGroupCfg for adding new devices*/
func (dbc *DatabaseCfg) AddProductGroupCfg(dev ProductGroupCfg) (int64, error) {
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

/*DelProductGroupCfg for deleting influx databases from ID*/
func (dbc *DatabaseCfg) DelProductGroupCfg(id string) (int64, error) {
	var affecteddev, affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affecteddev, err = session.Where("kapacitorid='" + id + "'").Cols("kapacitorid").Update(&AlertIDCfg{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on Delete ProductGroupCfg with id: %s, error: %s", id, err)
	}

	affected, err = session.Where("id='" + id + "'").Delete(&ProductGroupCfg{})
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

/*UpdateProductGroupCfg for adding new influxdb*/
func (dbc *DatabaseCfg) UpdateProductGroupCfg(id string, dev ProductGroupCfg) (int64, error) {
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

	log.Infof("Updated KapacitorID Config Successfully with id %s and data:%+v, affected", id, dev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*GetProductGroupCfgAffectOnDel for deleting devices from ID*/
func (dbc *DatabaseCfg) GetProductGroupCfgAffectOnDel(id string) ([]*DbObjAction, error) {
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
