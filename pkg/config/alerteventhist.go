package config

import (
	"fmt"
	"strconv"
	"time"
)

/***************************
	AlertEventHist Alert events
	-GetAlertEventHistByID(struct)
	-GetAlertEventHistMap (map - for interna config use
	-GetAlertEventHistArray(Array - for web ui use )
***********************************/

/*GetAlertEventHistByID get device data by id*/
func (dbc *DatabaseCfg) GetAlertEventHistByID(id int64) (AlertEventHist, error) {
	alevtarray, err := dbc.GetAlertEventHistArray("id='" + strconv.FormatInt(id, 10) + "'")
	if err != nil {
		return AlertEventHist{}, err
	}
	if len(alevtarray) > 1 {
		return AlertEventHist{}, fmt.Errorf("Error %d results on get AlertEventHistArray by id %d", len(alevtarray), id)
	}
	if len(alevtarray) == 0 {
		return AlertEventHist{}, fmt.Errorf("Error no values have been returned with this id %d in the config table", id)
	}
	return *alevtarray[0], nil
}

/*GetAlertEventHistMap  return data in map format*/
func (dbc *DatabaseCfg) GetAlertEventHistMap(filter string) (map[int64]*AlertEventHist, error) {
	alevtarray, err := dbc.GetAlertEventHistArray(filter)
	alevtmap := make(map[int64]*AlertEventHist)
	for _, val := range alevtarray {
		alevtmap[val.ID] = val
		log.Debugf("%+v", *val)
	}
	return alevtmap, err
}

/*GetAlertEventHistArray generate an array of Alert Events History with all their information */
func (dbc *DatabaseCfg) GetAlertEventHistArray(filter string) ([]*AlertEventHist, error) {
	log.Debugf("Getting AlertEventHist data filtered with %s", filter)
	var err error
	var alevthistarray []*AlertEventHist
	if err = dbc.x.Where(filter).Desc("id").Find(&alevthistarray); err != nil {
		log.Warnf("Fail to get AlertEventHist data filtered with %s : %v\n", filter, err)
		return nil, err
	}
	return alevthistarray, nil
}

/*GetAlertEventsHistByLevelArray generates an array of alert events summary */
func (dbc *DatabaseCfg) GetAlertEventsHistByLevelArray(filter string) ([]*AlertEventsSummary, error) {
	log.Debugf("Getting AlertEventsHistSummary data filtered with %s", filter)
	var err error
	var alevtsummarray []*AlertEventsSummary
	sqlquery := "SELECT level, count(*) as num FROM alert_event_hist"
	if len(filter) > 0 {
		sqlquery = sqlquery + " WHERE " + filter
	}
	sqlquery = sqlquery + " GROUP BY level"
	if err = dbc.x.SQL(sqlquery).Find(&alevtsummarray); err != nil {
		log.Warnf("Error getting AlertEventsHistSummary data filtered with %s: %v\n", filter, err)
		return nil, err
	}
	return alevtsummarray, nil
}

/*GetAlertEventHistArrayWithParams generate an array of Alert Events History with all their information */
func (dbc *DatabaseCfg) GetAlertEventHistArrayWithParams(filter string, page int64, itemsPerPage int64, maxSize int64, sortColumn string, sortDir string) ([]*AlertEventHist, error) {
	log.Debugf("Getting AlertEventHist data filtered with filter: %s, page: %d, itemsPerPage: %d, maxSize: %d, sortColumn: %s, sortDir: %s", filter, page, itemsPerPage, maxSize, sortColumn, sortDir)
	var err error
	var alevts []*AlertEventHist
	/*
		SELECT * FROM alert_event_hist
		WHERE filter
		ORDER BY sortColumn
		sortDir (ASC/DESC)
		LIMIT itemsPerPage * maxSize
		OFFSET itemsPerPage * (page - 1)
	*/
	sqlquery := "SELECT * FROM alert_event_hist"
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

	log.Debugf("Getting AlertEventHist data filtered with sqlquery: %s", sqlquery)

	start := time.Now()
	if err = dbc.x.SQL(sqlquery).Find(&alevts); err != nil {
		log.Warnf("Fail to get AlertEventHist data filtered with sqlquery: %s. Error : %s", sqlquery, err)
		return nil, err
	}
	elapsed := time.Since(start)
	log.Debugf("TIMELOG: GetAlertEventHistArrayWithParams took %v", elapsed)
	return alevts, nil
}

/*AddAlertEventHist for adding new Alert Event History*/
func (dbc *DatabaseCfg) AddAlertEventHist(dev *AlertEventHist) (int64, error) {
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
	log.Infof("Added new Alert event history Successfully with id %d", dev.ID)
	dbc.addChanges(affected)
	return affected, nil
}

/*DelAlertEventHist for deleting alert events from list of IDs*/
func (dbc *DatabaseCfg) DelAlertEventHist(idlist string) (int64, error) {
	var affected int64
	var err error

	session := dbc.x.NewSession()
	defer session.Close()

	affected, err = session.Where("id IN (" + idlist + ")").Delete(&AlertEventHist{})
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

/*GetAlertEventHistAffectOnDel NO config tables affected when deleting alert events*/
func (dbc *DatabaseCfg) GetAlertEventHistAffectOnDel(id string) ([]*DbObjAction, error) {
	var obj []*DbObjAction
	return obj, nil
}
