package config

import (
	"fmt"
)

/***************************
	OperationCfg instructions
	-GetOperationCfgCfgByID(struct)
	-GetOperationCfgMap (map - for internal config use
	-GetOperationCfgArray(Array - for web ui use )
	-AddOperationCfg
	-DelOperationCfg
	-UpdateOperationCfg
	-GetOperationCfgAffectOnDel
***********************************/

/*GetOperationCfgByID get operation data by id*/
func (dbc *DatabaseCfg) GetOperationCfgByID(id string) (OperationCfg, error) {
	cfgarray, err := dbc.GetOperationCfgArray("id='" + id + "'")
	if err != nil {
		return OperationCfg{}, err
	}
	if len(cfgarray) > 1 {
		return OperationCfg{}, fmt.Errorf("Error %d results on get OperationCfg by id %s", len(cfgarray), id)
	}
	if len(cfgarray) == 0 {
		return OperationCfg{}, fmt.Errorf("Error no values have been returned with this id %s in the config table", id)
	}
	return *cfgarray[0], nil
}

/*GetOperationCfgMap  return data in map format*/
func (dbc *DatabaseCfg) GetOperationCfgMap(filter string) (map[string]*OperationCfg, error) {
	cfgarray, err := dbc.GetOperationCfgArray(filter)
	cfgmap := make(map[string]*OperationCfg)
	for _, val := range cfgarray {
		cfgmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return cfgmap, err
}

/*GetOperationCfgArray generate an array of operations with all its information */
func (dbc *DatabaseCfg) GetOperationCfgArray(filter string) ([]*OperationCfg, error) {
	var err error
	var operations []*OperationCfg
	//Get Only data for selected operations
	if len(filter) > 0 {
		if err = dbc.x.Where(filter).Find(&operations); err != nil {
			log.Warnf("Fail to get OperationCfg data filtered with %s : %v\n", filter, err)
			return nil, err
		}
	} else {
		if err = dbc.x.Find(&operations); err != nil {
			log.Warnf("Fail to get OperationCfg data: %v\n", err)
			return nil, err
		}
	}
	return operations, nil
}

/*AddOperationCfg for adding new operations*/
func (dbc *DatabaseCfg) AddOperationCfg(dev *OperationCfg) (int64, error) {
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
	log.Infof("Added new operation Successfully with id %s ", dev.ID)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelOperationCfg for deleting operations from ID*/
func (dbc *DatabaseCfg) DelOperationCfg(id string) (int64, error) {
	var affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Where("id='" + id + "'").Delete(&OperationCfg{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Deleted Successfully operation with ID %s", id)
	dbc.addChanges(affected)
	return affected, nil
}

/*UpdateOperationCfg for updating operation*/
func (dbc *DatabaseCfg) UpdateOperationCfg(id string, dev *OperationCfg) (int64, error) {
	var affecteddev, affected int64
	var err error
	session := dbc.x.NewSession()
	defer session.Close()
	if id != dev.ID { //ID has been changed
		affecteddev, err = session.Where("operationid='" + id + "'").Cols("operationid").Update(&AlertIDCfg{OperationID: dev.ID})
		if err != nil {
			session.Rollback()
			return 0, fmt.Errorf("Error on Update AlertIDCfg, on update operationid (old) %s with (new): %s, error: %s", id, dev.ID, err)
		}
		log.Infof("Updated OperationID to %d alerts ", affecteddev)
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

	log.Infof("Updated operation Successfully with id %s and data:%+v", id, dev)
	dbc.addChanges(affected + affecteddev)
	return affected, nil
}

/*GetOperationCfgAffectOnDel for deleting operations from ID*/
func (dbc *DatabaseCfg) GetOperationCfgAffectOnDel(id string) ([]*DbObjAction, error) {
	var alerts []*AlertIDCfg
	var obj []*DbObjAction
	if err := dbc.x.Where("operationid='" + id + "'").Find(&alerts); err != nil {
		log.Warnf("Error on Get operation id %d for alerts , error: %s", id, err)
		return nil, err
	}

	for _, val := range alerts {
		obj = append(obj, &DbObjAction{
			Type:     "alertidcfg",
			TypeDesc: "",
			ObID:     val.ID,
			Action:   "Change alert to Other operation",
		})

	}
	return obj, nil
}
