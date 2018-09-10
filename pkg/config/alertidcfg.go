package config

import (
	"fmt"
)

/***************************
	AlertIDCfg DB backends
	-GetAlertIDCfgCfgByID(struct)
	-GetAlertIDCfgMap (map - for interna config use
	-GetAlertIDCfgArray(Array - for web ui use )
	-AddAlertIDCfg
	-DelAlertIDCfg
	-UpdateAlertIDCfg
  -GetAlertIDCfgAffectOnDel
***********************************/

/*GetAlertIDCfgByID get device data by id*/
func (dbc *DatabaseCfg) GetAlertIDCfgByID(id string) (AlertIDCfg, error) {
	cfgarray, err := dbc.GetAlertIDCfgArray("id='" + id + "'")
	if err != nil {
		return AlertIDCfg{}, err
	}
	if len(cfgarray) > 1 {
		return AlertIDCfg{}, fmt.Errorf("Error %d results on get AlertIDCfgArray by id %s", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return AlertIDCfg{}, fmt.Errorf("Error no values have been returned with this id %s in the config table", id)
	}
	return *cfgarray[0], nil
}

/*GetAlertIDCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetAlertIDCfgMap(filter string) (map[string]*AlertIDCfg, error) {
	cfgarray, err := dbc.GetAlertIDCfgArray(filter)
	cfgmap := make(map[string]*AlertIDCfg)
	for _, val := range cfgarray {
		cfgmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetAlertIDCfgArray generate an array of devices with all its information */
func (dbc *DatabaseCfg) GetAlertIDCfgArray(filter string) ([]*AlertIDCfg, error) {
	var err error
	var devices []*AlertIDCfg
	//Get Only data for selected devices
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&devices); err != nil {
			log.Warnf("Fail to get AlertIDCfg data filtered with %s : %v\n", filter, err)
			return nil, err
		}
	} else {
		if err = dbc.x.Find(&devices); err != nil {
			log.Warnf("Fail to get AlertIDCfg data: %v\n", err)
			return nil, err
		}
	}

	//Asign Endpoints to AlertID's
	var alerthttp []*AlertEndpointRel
	if err = dbc.x.Find(&alerthttp); err != nil {
		log.Warnf("Fail to get AlertIDCfg and its Endpoint relationship data: %v\n", err)
		return devices, err
	}

	//Load Measurements and metrics relationship
	//We assign field metric ID to each measurement
	for _, v := range devices {
		for _, r := range alerthttp {
			if r.AlertID == v.ID {
				v.Endpoint = append(v.Endpoint, r.EndpointID)
			}
		}
	}

	return devices, nil
}

/*AddAlertIDCfg for adding new devices*/
func (dbc *DatabaseCfg) AddAlertIDCfg(dev *AlertIDCfg) (int64, error) {
	var err error
	var affected, newo int64
	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Insert(dev)
	if err != nil {
		session.Rollback()
		return 0, err
	}
	//save new alert endpoint relations
	for _, o := range dev.Endpoint {

		ostruct := AlertEndpointRel{
			AlertID:    dev.ID,
			EndpointID: o,
		}
		newo, err = session.Insert(&ostruct)
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
	log.Infof("Added new Alert ID Successfully with id %s [ %d Endpoint] ", dev.ID, newo)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelAlertIDCfg for deleting alerts from ID*/
func (dbc *DatabaseCfg) DelAlertIDCfg(id string) (int64, error) {
	var affectedouts, affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	//first deleting references in AlertEndpointRel
	affectedouts, err = session.Where("alert_id='" + id + "'").Delete(&AlertEndpointRel{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on delete AlertEndpointRel with id: %s, error: %s", id, err)
	}

	affected, err = session.Where("id='" + id + "'").Delete(&AlertIDCfg{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Deleted Successfully Alert with ID %s [ %d Affected alerts   ] [%d affected endpoints]", id, affected, affectedouts)
	dbc.addChanges(affected + affectedouts)
	return affected, nil
}

/*UpdateAlertIDCfg for updating AlertIDCfg*/
func (dbc *DatabaseCfg) UpdateAlertIDCfg(id string, dev *AlertIDCfg) (int64, error) {
	var affected int64
	var err error
	session := dbc.x.NewSession()
	defer session.Close()
	//first deleting references in AlertEndpointRel
	_, err = session.Where("alert_id='" + id + "'").Delete(&AlertEndpointRel{})
	if err != nil {
		session.Rollback()
		return 0, fmt.Errorf("Error on delete AlertEndpointRel with id: %s, error: %s", id, err)
	}
	//save again alert endpoint relations
	for _, o := range dev.Endpoint {

		ostruct := AlertEndpointRel{
			AlertID:    dev.ID,
			EndpointID: o,
		}
		_, err = session.Insert(&ostruct)
		if err != nil {
			session.Rollback()
			return 0, err
		}
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

	log.Infof("Updated Alert ID Config Successfully with id %s and data:%+v, %d affected", id, dev, affected)
	dbc.addChanges(affected)
	return affected, nil
}

/*GetAlertIDCfgAffectOnDel NO config tables affected when deleting AlertIDCfg */
func (dbc *DatabaseCfg) GetAlertIDCfgAffectOnDel(id string) ([]*DbObjAction, error) {
	//var devices []*AlertIDCfg
	var obj []*DbObjAction

	return obj, nil
}
