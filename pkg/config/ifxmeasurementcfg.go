package config

import (
	"fmt"
	"strconv"
	"strings"
)

/***************************
	IfxMeasurementCfg DB backends
	-GetIfxMeasurementCfgByID(struct)
	-GetIfxMeasurementCfgMap (map - for interna config use
	-GetIfxMeasurementCfgArray(Array - for web ui use )
	-GetIfxMeasurementCfgBySQLQuery(struct - for web ui use )
	-GetIfxMeasurementCfgDistinctNamesArray
	-GetIfxMeasurementTagsArray
	-AddIfxMeasurementCfg
	-DelIfxMeasurementCfg
	-UpdateIfxMeasurementCfg
  -GetIfxMeasurementCfgAffectOnDel
***********************************/

/*GetIfxMeasurementCfgByID get device data by id*/
func (dbc *DatabaseCfg) GetIfxMeasurementCfgByID(id int64) (IfxMeasurementCfg, error) {
	cfgarray, err := dbc.GetIfxMeasurementCfgArray("id='" + strconv.FormatInt(id, 10) + "'")
	if err != nil {
		return IfxMeasurementCfg{}, err
	}
	if len(cfgarray) > 1 {
		return IfxMeasurementCfg{}, fmt.Errorf("Error %d results on get IfxMeasurementCfg by id %d", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return IfxMeasurementCfg{}, fmt.Errorf("Error no values have been returned with this id %d in the influx config table", id)
	}
	return *cfgarray[0], nil
}

/*GetIfxMeasurementCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetIfxMeasurementCfgMap(filter string) (map[int64]*IfxMeasurementCfg, error) {
	cfgarray, err := dbc.GetIfxMeasurementCfgArray(filter)
	cfgmap := make(map[int64]*IfxMeasurementCfg)
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
			log.Warnf("Fail to get IfxMeasurementCfg data filtered with %s : %v\n", filter, err)
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

/*GetIfxMeasurementCfgBySQLQuery Gets an IfxMeasurementCfg with all its information */
func (dbc *DatabaseCfg) GetIfxMeasurementCfgBySQLQuery(sqlquery string) (IfxMeasurementCfg, error) {
	var err error
	var cfgarray []*IfxMeasurementCfg
	//Get Only data for selected devices
	if err = dbc.x.SQL(sqlquery).Find(&cfgarray); err != nil {
		log.Warnf("Failed getting IfxMeasurementCfg data with sql %s. Error: %s\n", sqlquery, err)
		return IfxMeasurementCfg{}, err
	}
	if len(cfgarray) > 1 {
		return IfxMeasurementCfg{}, fmt.Errorf("Error %d results getting IfxMeasurementCfg by sqlquery %s", len(cfgarray), sqlquery)
	}
	if len(cfgarray) == 0 {
		return IfxMeasurementCfg{}, fmt.Errorf("Error no values have been returned with this sqlquery %s in the config table", sqlquery)
	}
	return *cfgarray[0], nil
}

/*GetIfxMeasurementCfgDistinctNamesArray generate an array of IfxMeasurementCfg with distinct names */
func (dbc *DatabaseCfg) GetIfxMeasurementCfgDistinctNamesArray(filter string) ([]*IfxMeasurementCfg, error) {
	var err error
	var msmts []*IfxMeasurementCfg
	if err = dbc.x.Distinct("name").Where(filter).Find(&msmts); err != nil {
		log.Warnf("Failed to get IfxMeasurementCfg data filtered with %s : %v\n", filter, err)
		return nil, err
	}
	return msmts, nil
}

/*GetIfxMeasurementTagsArray Gets the array of tags for the measurements passed in filter */
/*The filter contains a list of measurement names*/
/*then with these measurement names a list of tags is obtained*/
func (dbc *DatabaseCfg) GetIfxMeasurementTagsArray(filter string) ([]string, error) {
	var err error
	var tags []string
	var msmts []*IfxMeasurementCfg
	var namesfilter string
	namesfilter = "`name` IN ('" + strings.Replace(filter, ",", "','", -1) + "')"
	//Get the list of measurement tags
	if err = dbc.x.Where(namesfilter).Find(&msmts); err != nil {
		log.Warnf("Failed to get IfxMeasurementCfg data filtered with (%s). Error: %v", filter, err)
		return nil, err
	}
	for _, msmt := range msmts {
		for _, tag := range msmt.Tags {
			//Don't add duplicates
			if !strings.Contains(","+strings.Join(tags, ",")+",", ","+tag+",") {
				tags = append(tags, tag)
			}
		}
	}
	log.Infof("Got Measurement Tags Successfully: %+v", tags)
	return tags, nil
}

/*AddIfxMeasurementCfg for adding new devices*/
func (dbc *DatabaseCfg) AddIfxMeasurementCfg(dev *IfxMeasurementCfg) (int64, error) {
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
	log.Infof("Added new InfluxMeasurement Successfully with id %d ", dev.ID)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelIfxMeasurementCfg for deleting influx databases from ID*/
func (dbc *DatabaseCfg) DelIfxMeasurementCfg(id int64) (int64, error) {
	var affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Where("id='" + strconv.FormatInt(id, 10) + "'").Delete(&IfxMeasurementCfg{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Deleted Successfully influx measurements with ID %d [ %d Devices Affected  ]", id, affected)
	dbc.addChanges(affected)
	return affected, nil
}

/*UpdateIfxMeasurementCfg for adding new influxdb*/
func (dbc *DatabaseCfg) UpdateIfxMeasurementCfg(id int64, dev *IfxMeasurementCfg) (int64, error) {
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

	affected, err = session.Where("id='" + strconv.FormatInt(id, 10) + "'").UseBool().AllCols().Update(dev)
	if err != nil {
		session.Rollback()
		return 0, err
	}
	err = session.Commit()
	if err != nil {
		return 0, err
	}

	log.Infof("Updated Influx Measurement Successfully with id %d and data:%+v, affected", id, dev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*GetIfxMeasurementCfgAffectOnDel for deleting devices from ID*/
func (dbc *DatabaseCfg) GetIfxMeasurementCfgAffectOnDel(id int64) ([]*DbObjAction, error) {
	var devices []*AlertIDCfg
	var obj []*DbObjAction
	if err := dbc.x.Where("kapacitorid='" + strconv.FormatInt(id, 10) + "'").Find(&devices); err != nil {
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
