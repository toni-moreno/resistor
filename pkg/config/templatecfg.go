package config

import (
	"fmt"
)

/***************************
	TemplateCfg DB backends
	-GetTemplateCfgCfgByID(struct)
	-GetTemplateCfgMap (map - for interna config use
	-GetTemplateCfgArray(Array - for web ui use )
	-AddTemplateCfg
	-DelTemplateCfg
	-UpdateTemplateCfg
  -GetTemplateCfgAffectOnDel
***********************************/

/*GetTemplateCfgByID get template data by id*/
func (dbc *DatabaseCfg) GetTemplateCfgByID(id string) (TemplateCfg, error) {
	cfgarray, err := dbc.GetTemplateCfgArray("id='" + id + "'")
	if err != nil {
		return TemplateCfg{}, err
	}
	if len(cfgarray) > 1 {
		return TemplateCfg{}, fmt.Errorf("Error %d results on get TemplateCfg by id %s", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return TemplateCfg{}, fmt.Errorf("Error no values have been returned with this id %s in the influx config table", id)
	}
	return *cfgarray[0], nil
}

/*GetTemplateCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetTemplateCfgMap(filter string) (map[string]*TemplateCfg, error) {
	cfgarray, err := dbc.GetTemplateCfgArray(filter)
	cfgmap := make(map[string]*TemplateCfg)
	for _, val := range cfgarray {
		cfgmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetTemplateCfgArray generate an array of templates with all its information */
func (dbc *DatabaseCfg) GetTemplateCfgArray(filter string) ([]*TemplateCfg, error) {
	var err error
	var templates []*TemplateCfg
	//Get Only data for selected templates
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&templates); err != nil {
			log.Warnf("Fail to get TemplateCfg data filtered with %s : %v\n", filter, err)
			return nil, err
		}
	} else {
		if err = dbc.x.Find(&templates); err != nil {
			log.Warnf("Fail to get TemplateCfg data: %v\n", err)
			return nil, err
		}
	}
	return templates, nil
}

/*AddTemplateCfg for adding new templates*/
func (dbc *DatabaseCfg) AddTemplateCfg(dev *TemplateCfg) (int64, error) {
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
	log.Infof("Added new Template Successfully with id %s ", dev.ID)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelTemplateCfg for deleting influx databases from ID*/
func (dbc *DatabaseCfg) DelTemplateCfg(id string) (int64, error) {
	var affecteddev, affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Where("id='" + id + "'").Delete(&TemplateCfg{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}

	//log.Infof("Deleted Successfully Template with ID %s [ %d templates Affected  ]", id, affecteddev)
	log.Infof("Deleted Successfully Template with ID %s [ %d templates Affected  ]", id, affected)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*UpdateTemplateCfg for adding new influxdb*/
func (dbc *DatabaseCfg) UpdateTemplateCfg(id string, dev *TemplateCfg) (int64, error) {
	var affecteddev, affected int64
	var err error
	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Where("id='" + id + "'").UseBool().AllCols().Update(dev)
	if err != nil {
		session.Rollback()
		return 0, err
	}
	err = session.Commit()
	if err != nil {
		return 0, err
	}

	log.Infof("Updated Template Config Successfully with id %s and data:%+v, affected", id, dev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*GetTemplateCfgAffectOnDel for deleting templates from ID*/
func (dbc *DatabaseCfg) GetTemplateCfgAffectOnDel(id string) ([]*DbObjAction, error) {
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
