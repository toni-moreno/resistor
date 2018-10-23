package config

import (
	"fmt"
	"strconv"
)

/***************************
	AlertEvent Alert events
	-GetAlertEventByID(struct)
	-GetAlertEventMap (map - for interna config use
	-GetAlertEventArray(Array - for web ui use )
***********************************/

/*GetAlertEventByID get device data by id*/
func (dbc *DatabaseCfg) GetAlertEventByID(id int64) (AlertEvent, error) {
	alevtarray, err := dbc.GetAlertEventArray("id='" + strconv.FormatInt(id, 10) + "'")
	if err != nil {
		return AlertEvent{}, err
	}
	if len(alevtarray) > 1 {
		return AlertEvent{}, fmt.Errorf("Error %d results on get AlertEventArray by id %d", len(alevtarray), id)
	}
	if len(alevtarray) == 0 {
		return AlertEvent{}, fmt.Errorf("Error no values have been returned with this id %d in the config table", id)
	}
	return *alevtarray[0], nil
}

/*GetAlertEventMap  return data in map format*/
func (dbc *DatabaseCfg) GetAlertEventMap(filter string) (map[int64]*AlertEvent, error) {
	alevtarray, err := dbc.GetAlertEventArray(filter)
	alevtmap := make(map[int64]*AlertEvent)
	for _, val := range alevtarray {
		alevtmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return alevtmap, err
}

/*GetAlertEventArray generate an array of alert events with all its information */
func (dbc *DatabaseCfg) GetAlertEventArray(filter string) ([]*AlertEvent, error) {
	log.Debugf("Getting AlertEvent data filtered with %s", filter)
	var err error
	var alevtarray []*AlertEvent
	if err = dbc.x.Where(filter).Desc("id").Find(&alevtarray); err != nil {
		log.Warnf("Fail to get AlertEvent data filtered with %s : %v\n", filter, err)
		return nil, err
	}
	return alevtarray, nil
}

/*GetAlertEventsByLevelArray generates an array of alert events summary */
func (dbc *DatabaseCfg) GetAlertEventsByLevelArray(filter string) ([]*AlertEventsSummary, error) {
	log.Debugf("Getting AlertEventsSummary data filtered with %s", filter)
	var err error
	var alevtsummarray []*AlertEventsSummary
	sqlquery := "SELECT level, count(*) as num FROM alert_event"
	if len(filter) > 0 {
		sqlquery = sqlquery + " WHERE " + filter
	}
	sqlquery = sqlquery + " GROUP BY level"
	if err = dbc.x.SQL(sqlquery).Find(&alevtsummarray); err != nil {
		log.Warnf("Error getting AlertEventsSummary data filtered with %s: %v\n", filter, err)
		return nil, err
	}
	return alevtsummarray, nil
}

/*GetAlertEventArrayWithParams generate an array of alert events with all its information */
func (dbc *DatabaseCfg) GetAlertEventArrayWithParams(filter string, page int64, itemsPerPage int64, maxSize int64, sortColumn string, sortDir string) ([]*AlertEvent, error) {
	log.Debugf("Getting AlertEvent data filtered with filter: %s, page: %d, itemsPerPage: %d, maxSize: %d, sortColumn: %s, sortDir: %s", filter, page, itemsPerPage, maxSize, sortColumn, sortDir)
	var err error
	var alevts []*AlertEvent
	/*
		SELECT * FROM alert_event
		WHERE filter
		ORDER BY sortColumn
		sortDir (ASC/DESC)
		LIMIT itemsPerPage * maxSize
		OFFSET itemsPerPage * (page - 1)
	*/
	sqlquery := "SELECT * FROM alert_event"
	if len(filter) > 0 {
		sqlquery = sqlquery + " WHERE " + filter
	}
	if len(sortColumn) > 0 {
		sqlquery = sqlquery + " ORDER BY " + sortColumn
	}
	if len(sortDir) > 0 {
		sqlquery = sqlquery + " " + sortDir
	}
	if itemsPerPage > 0 && maxSize > 0 {
		limit := itemsPerPage * maxSize
		sqlquery = sqlquery + " LIMIT " + strconv.FormatInt(limit, 10)
	}
	if itemsPerPage > 0 && maxSize > 0 && page > 0 {
		offset := itemsPerPage * maxSize * ((page - 1) / maxSize)
		sqlquery = sqlquery + " OFFSET " + strconv.FormatInt(offset, 10)
	}

	log.Debugf("Getting AlertEvent data filtered with sqlquery: %s", sqlquery)

	if err = dbc.x.SQL(sqlquery).Find(&alevts); err != nil {
		log.Warnf("Fail to get AlertEvent data filtered with sqlquery: %s. Error : %s", sqlquery, err)
		return nil, err
	}
	return alevts, nil
}

/*AddAlertEvent for adding new alert events*/
func (dbc *DatabaseCfg) AddAlertEvent(dev *AlertEvent) (int64, error) {
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
	log.Infof("Added new Alert event Successfully with id %d", dev.ID)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelAlertEvent for deleting alert events from list of IDs*/
func (dbc *DatabaseCfg) DelAlertEvent(idlist string) (int64, error) {
	var affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Where("id IN (" + idlist + ")").Delete(&AlertEvent{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}
	log.Infof("Deleted Successfully Alert events with IDs %s [ %d Affected alert events]", idlist, affected)
	dbc.addChanges(affected)
	return affected, nil
}

/*GetAlertEventAffectOnDel NO config tables affected when deleting alert events*/
func (dbc *DatabaseCfg) GetAlertEventAffectOnDel(id string) ([]*DbObjAction, error) {
	var obj []*DbObjAction
	return obj, nil
}
