package config

import (
	"fmt"
	"strconv"
)

/***************************
	IfxDBCfg DB backends
	-GetIfxDBCfgCfgByID(struct)
	-GetIfxDBCfgMap (map - for interna config use
	-GetIfxDBCfgArray(Array - for web ui use )
	-AddIfxDBCfg
	-DelIfxDBCfg
	-UpdateIfxDBCfg
  -GetIfxDBCfgAffectOnDel
***********************************/

/*GetIfxDBCfgByID get device data by id*/
func (dbc *DatabaseCfg) GetIfxDBCfgByID(id int64) (IfxDBCfg, error) {
	cfgarray, err := dbc.GetIfxDBCfgArray("id='" + strconv.FormatInt(id, 10) + "'")
	if err != nil {
		return IfxDBCfg{}, err
	}
	if len(cfgarray) > 1 {
		return IfxDBCfg{}, fmt.Errorf("Error %d results on get IfxDBCfg by id %d", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return IfxDBCfg{}, fmt.Errorf("Error no values have been returned with this id %d in the influx config table", id)
	}
	return *cfgarray[0], nil
}

/*GetIfxDBCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetIfxDBCfgMap(filter string) (map[int64]*IfxDBCfg, error) {
	cfgarray, err := dbc.GetIfxDBCfgArray(filter)
	cfgmap := make(map[int64]*IfxDBCfg)
	for _, val := range cfgarray {
		cfgmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetIfxDBCfgArray generate an array of devices with all its information */
func (dbc *DatabaseCfg) GetIfxDBCfgArray(filter string) ([]*IfxDBCfg, error) {
	var err error
	var devices []*IfxDBCfg
	//Get Only data for selected devices
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&devices); err != nil {
			log.Warnf("Fail to get IfxDBCfg  data filteter with %s : %v\n", filter, err)
			return nil, err
		}
	} else {
		if err = dbc.x.Find(&devices); err != nil {
			log.Warnf("Fail to get influxcfg   data: %v\n", err)
			return nil, err
		}
	}

	for _, mVal := range devices {

		//Measurements for each DB
		var dbmeas []*IfxDBMeasRel
		if err = dbc.x.Where("ifxdbid ==" + strconv.FormatInt(mVal.ID, 10)).Find(&dbmeas); err != nil {
			log.Warnf("Fail to get MGroup Measurements relationship  data: %v\n", err)
		}

		for _, mgm := range dbmeas {
			mVal.Measurements = append(mVal.Measurements, mgm.IfxMeasID)
		}
	}
	return devices, nil
}

/*AddIfxDBCfg for adding new devices*/
func (dbc *DatabaseCfg) AddIfxDBCfg(dev IfxDBCfg) (int64, error) {
	var err error
	var affected int64
	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Insert(dev)
	if err != nil {
		session.Rollback()
		return 0, err
	}

	//Measurement Fields
	for _, meas := range dev.Measurements {
		mstruct := IfxDBMeasRel{
			IfxDBID:   dev.ID,
			IfxMeasID: meas,
		}
		_, err = session.Insert(&mstruct)
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
	log.Infof("Added new Kapacitor backend Successfully with id %s ", dev.ID)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelIfxDBCfg for deleting influx databases from ID*/
func (dbc *DatabaseCfg) DelIfxDBCfg(id int64) (int64, error) {
	var affecteddev, affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affecteddev, err = session.Where("ifxdbid ==" + strconv.FormatInt(id, 10)).Delete(&IfxDBMeasRel{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on Delete Metric with id on delete IfxDBCfg with id: %d, error: %s", id, err)
	}

	affected, err = session.Where("id='" + strconv.FormatInt(id, 10) + "'").Delete(&IfxDBCfg{})
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

// AddOrUpdateIfxDBCfg this method insert data if not previouosly exist the tuple ifxServer.Name or update it if already exist
func (dbc *DatabaseCfg) AddOrUpdateIfxDBCfg(dev IfxDBCfg) (int64, error) {
	log.Debugf("ADD OR UPDATE %+v", dev)
	//check if exist
	m, err := dbc.GetIfxDBCfgArray("name == '" + dev.Name + "' AND ifxserver == '" + dev.IfxServer + "'")
	if err != nil {
		return 0, err
	}
	switch len(m) {
	case 1:
		log.Debugf("Updating InfluxDB %+v", m)
		return dbc.UpdateIfxDBCfg(m[0].ID, dev)
	case 0:
		log.Debugf("Adding new InfluxDB %+v", dev)
		return dbc.AddIfxDBCfg(dev)
	default:
		log.Errorf("There is some error when searching for db %+v , found %d", dev, len(m))
		return 0, fmt.Errorf("There is some error when searching for db %+v , found %d", dev, len(m))
	}

}

/*UpdateIfxDBCfg for adding new influxdb*/
func (dbc *DatabaseCfg) UpdateIfxDBCfg(id int64, dev IfxDBCfg) (int64, error) {
	var affecteddev, affected int64
	var err error

	//first get ID

	session := dbc.x.NewSession()
	defer session.Close()

	//first check if id < 1 => search for the current ID for the unique IfxServer.Name

	affected, err = session.Where("id='" + strconv.FormatInt(id, 10) + "'").UseBool().AllCols().Update(dev)
	if err != nil {
		session.Rollback()
		return 0, err
	}

	//Delete relations
	affecteddev, err = session.Where("ifxdbid ==" + strconv.FormatInt(id, 10)).Delete(&IfxDBMeasRel{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on Delete Metric with id on delete IfxDBCfg with id: %d, error: %s", id, err)
	}
	//Add New Relations
	//Measurement Fields
	for _, meas := range dev.Measurements {
		mstruct := IfxDBMeasRel{
			IfxDBID:   id,
			IfxMeasID: meas,
		}
		_, err = session.Insert(&mstruct)
		if err != nil {
			session.Rollback()
			return 0, err
		}
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}

	log.Infof("Updated KapacitorID Config Successfully with id %s and data:%+v, affected", id, dev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*GetIfxDBCfgAffectOnDel for deleting devices from ID*/
func (dbc *DatabaseCfg) GetIfxDBCfgAffectOnDel(id int64) ([]*DbObjAction, error) {
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
