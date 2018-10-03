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
	-GetIfxDBCfgArrayByMeasName(Array - for web ui use )
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
		return IfxDBCfg{}, fmt.Errorf("Error no values have been returned with this id %d in the config table", id)
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
		if err = dbc.x.Where("ifxdbid =" + strconv.FormatInt(mVal.ID, 10)).Find(&dbmeas); err != nil {
			log.Warnf("Fail to get MGroup Measurements relationship  data: %v\n", err)
		}

		for _, m := range dbmeas {
			mVal.Measurements = append(mVal.Measurements, &ItemComponent{ID: m.IfxMeasID, Name: m.IfxMeasName})
		}

		/*results, err := dbc.x.Query("select rel.ifxdbid as dbid , rel.ifxmeasid as measid , meas.name as measname  from  ifx_db_meas_rel as rel , ifx_measurement_cfg as meas  where rel.ifxmeasid  = meas.ID and rel.ifxmeasid = " + strconv.FormatInt(mVal.ID, 10))
		if err != nil {
			log.Warnf("Fail to Query DB to Measurement : %s", err)
		}

		for _, item := range results {
			id, err := strconv.ParseInt(string(item["measid"]), 10, 64)
			if err != nil {
				log.Warnf("Fail to parse int from select  data: %v\n", err)
			}
			lab := string(item["measname"])

			mVal.Measurements = append(mVal.Measurements, &ItemComponent{ID: id, Label: lab})
		}*/
	}
	return devices, nil
}

/*GetIfxDBCfgArrayByMeasName Gets an array of Influx databases with their measurements */
func (dbc *DatabaseCfg) GetIfxDBCfgArrayByMeasName(filter string) ([]*IfxDBCfg, error) {
	var err error
	var devices []*IfxDBCfg
	sqlquery := "select distinct id, name, ifxserver, retention, description from ifx_db_cfg, ifx_db_meas_rel where ifx_db_meas_rel.ifxdbid = ifx_db_cfg.id "
	if len(filter) > 0 {
		sqlquery = sqlquery + " and ifxmeasname = '" + filter + "' "
	}
	sqlquery = sqlquery + " order by name, id"
	//Get Only data for selected devices
	if err = dbc.x.SQL(sqlquery).Find(&devices); err != nil {
		log.Warnf("Fail to get IfxDBCfg data filtered with %s : %v\n", filter, err)
		return nil, err
	}

	for _, mVal := range devices {

		//Measurements for each DB
		var dbmeas []*IfxDBMeasRel
		if err = dbc.x.Where("ifxdbid =" + strconv.FormatInt(mVal.ID, 10)).And("ifxmeasname = '" + filter + "'").Find(&dbmeas); err != nil {
			log.Warnf("Fail to get MGroup Measurements relationship data: %v\n", err)
		}

		for _, m := range dbmeas {
			mVal.Measurements = append(mVal.Measurements, &ItemComponent{ID: m.IfxMeasID, Name: m.IfxMeasName})
		}

	}
	return devices, nil
}

/*AddIfxDBCfg for adding new devices*/
func (dbc *DatabaseCfg) AddIfxDBCfg(dev *IfxDBCfg) (int64, error) {
	var err error
	var affected int64
	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Insert(dev)
	if err != nil {
		session.Rollback()
		return 0, err
	}

	//log.Debugf("IROW1 %d  | IROW2 %d ", irow, irow2)

	//Measurement Fields
	for _, meas := range dev.Measurements {
		mstruct := IfxDBMeasRel{
			IfxDBID:     dev.ID,
			IfxMeasID:   meas.ID,
			IfxMeasName: meas.Name,
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
	log.Infof("Added new Influx DB backend Config Successfully with id %d ", dev.ID)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelIfxDBCfg for deleting influx databases from ID*/
func (dbc *DatabaseCfg) DelIfxDBCfg(id int64) (int64, error) {
	var affecteddev, affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affecteddev, err = session.Where("ifxdbid =" + strconv.FormatInt(id, 10)).Delete(&IfxDBMeasRel{})
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
func (dbc *DatabaseCfg) AddOrUpdateIfxDBCfg(dev *IfxDBCfg) (int64, error) {
	log.Debugf("AddOrUpdateIfxDBCfg. ADD OR UPDATE %+v", dev)
	//check if exist
	m, err := dbc.GetIfxDBCfgArray("name = '" + dev.Name + "' AND ifxserver = '" + dev.IfxServer + "'")
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
func (dbc *DatabaseCfg) UpdateIfxDBCfg(nid int64, new *IfxDBCfg) (int64, error) {
	var affecteddev, affected int64
	var err error

	m, err := dbc.GetIfxDBCfgArray("name = '" + new.Name + "' AND ifxserver = '" + new.IfxServer + "'")
	if err != nil {
		return 0, err
	}
	old := m[0]
	//first get ID

	session := dbc.x.NewSession()
	defer session.Close()
	//Delete first

	//first check if id < 1 => search for the current ID for the unique IfxServer.Name

	affected, err = session.Where("id='" + strconv.FormatInt(old.ID, 10) + "'").UseBool().AllCols().Update(new)
	if err != nil {
		session.Rollback()
		return 0, err
	}

	//Delete current relations
	affecteddev, err = session.Where("ifxdbid =" + strconv.FormatInt(old.ID, 10)).Delete(&IfxDBMeasRel{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on Delete Metric with id on delete IfxDBCfg with id: %d, error: %s", old.ID, err)
	}
	//Delete old measurements new ones has been initialized witn  new ID's
	for _, meas := range old.Measurements {
		dbc.DelIfxMeasurementCfg(meas.ID)
	}
	//Add New Relations
	//Measurement Fields
	for _, meas := range new.Measurements {
		mstruct := IfxDBMeasRel{
			IfxDBID:     old.ID,
			IfxMeasID:   meas.ID,
			IfxMeasName: meas.Name,
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

	log.Infof("Updated Influx DB Config Successfully with id %d and data:%+v, affected", old.ID, new)
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
