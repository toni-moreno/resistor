package config

import "time"

//Real Time Filtering by device/alertid/or other tags

// DeviceStatCfg current stats by device
type DeviceStatCfg struct {
	ID             int64  `xorm:"'id' pk autoincr"`
	OrderID        int64  `xorm:"orderid"`
	DeviceID       string `xorm:"deviceid" binding:"Required"`
	AlertID        string `xorm:"alertid" binding:"Required"`
	ProductID      string `xorm:"productid" binding:"Required"`
	ExceptionID    int64  `xorm:"exceptionid"`
	Active         bool   `xorm:"active"`
	BaseLine       string `xorm:"baseline"`
	FilterTagKey   string `xorm:"filterTagKey"`
	FilterTagValue string `xorm:"filterTagValue"`
	Description    string `xorm:"description"`
}

// IfxServerCfg Influx server  config
type IfxServerCfg struct {
	ID          string `xorm:"'id' unique" binding:"Required"`
	URL         string `xorm:"URL" binding:"Required"`
	AdminUser   string `xorm:"adminuser"`
	AdminPasswd string `xorm:"adminpasswod"`
	Description string `xorm:"description"`
}

// ItemComponent for ID's/Names
type ItemComponent struct {
	ID   int64
	Name string
}

// IfxDBCfg Influx Database definition
type IfxDBCfg struct {
	ID           int64            `xorm:"'id' pk autoincr"`
	Name         string           `xorm:"'name'  not null unique(ifxdb)" binding:"Required"`
	IfxServer    string           `xorm:"'ifxserver'  not null unique(ifxdb)" binding:"Required"`
	Retention    []string         `xorm:"retention" binding:"Required"`
	Description  string           `xorm:"description"`
	Measurements []*ItemComponent `xorm:"-"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (IfxDBCfg) TableName() string {
	return "ifx_db_cfg"
}

// IfxMeasurementCfg Measurement Definition
type IfxMeasurementCfg struct {
	ID          int64    `xorm:"'id'  pk autoincr " binding:"Required"`
	Name        string   `xorm:"'name' not  null"`
	Tags        []string `xorm:"tags" binding:"Required"`
	Fields      []string `xorm:"fields" binding:"Required"`
	Description string   `xorm:"description"`
}

// IfxDBMeasRel Relationship between Ifx DB's and Its measurements
type IfxDBMeasRel struct {
	IfxDBID     int64  `xorm:"ifxdbid" `
	IfxMeasID   int64  `xorm:"ifxmeasid"`
	IfxMeasName string `xorm:"ifxmeasname"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (IfxDBMeasRel) TableName() string {
	return "ifx_db_meas_rel"
}

// ProductCfg Product Catalog Config type
type ProductCfg struct {
	ID          string   `xorm:"'id' unique" binding:"Required"`
	IDTagName   string   `xorm:"idtagname" binding:"Required"` //Set the
	CommonTags  []string `xorm:"commontags"`
	BaseLines   []string `xorm:"baselines"`
	Description string   `xorm:"description"`
}

// ProductGroupCfg Product Group Catalog Config type
type ProductGroupCfg struct {
	ID          string   `xorm:"'id' unique" binding:"Required"`
	Products    []string `xorm:"commontags"`
	Description string   `xorm:"description"`
}

// KapacitorCfg Kapacitor URL's config type
type KapacitorCfg struct {
	ID          string `xorm:"'id' unique" binding:"Required"`
	URL         string `xorm:"URL" binding:"Required"`
	Description string `xorm:"description"`
}

// RangeTimeCfg Range or periods Times definition config type
type RangeTimeCfg struct {
	ID          string `xorm:"'id' unique" binding:"Required"`
	MaxHour     int    `xorm:"'max_hour' default 23"`
	MinHour     int    `xorm:"'min_hour' default 0"`
	WeekDays    string `xorm:"'weekdays' default '0123456'"`
	Description string `xorm:"description"`
}

//AlertID definition for
// N.A,
// None,
// Last,
// Max,
// Mean,
// Median,
// Min,
// MovingAverage,
// Percentile,
// Spread,
// Stddev,
// Sum

// TemplateCfg Templating data strucr
type TemplateCfg struct {
	ID                      string    `xorm:"'id' unique" binding:"Required"`
	TrigerType              string    `xorm:"trigertype" binding:"Required;In(DEADMAN,THRESHOLD,TREND)"` //deadman
	StatFunc                string    `xorm:"statfunc"`
	CritDirection           string    `xorm:"critdirection"`
	ThresholdType           string    `xorm:"thresholdtype"` //Absolute/Relative
	TrendSign               string    `xorm:"trendsign"`     //Positive/Negative
	TplData                 string    `xorm:"tpldata"`
	Description             string    `xorm:"description"`
	Modified                time.Time `xorm:"modified"`
	ServersWOLastDeployment []string  `xorm:"servers_wo_last_deployment"`
}

// OutHTTPCfg Alert Destination HTTP based backends config
type OutHTTPCfg struct {
	ID          string `xorm:"'id' unique" binding:"Required"`
	Type        string `xorm:"type"`
	JSONConfig  string `xorm:"json_config"`
	Description string `xorm:"description"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (OutHTTPCfg) TableName() string {
	return "out_http_cfg"
}

// AlertHTTPOutRel  relation between Alerts and HTTP Out's type
type AlertHTTPOutRel struct {
	AlertID   string `xorm:"alert_id"`
	HTTPOutID string `xorm:"http_out_id"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (AlertHTTPOutRel) TableName() string {
	return "alert_http_out_rel"
}

// AlertIDCfg Alert Definition Config type
type AlertIDCfg struct {
	//Alert ID data
	ID          string `xorm:"'id' unique" binding:"Required"` //Autogenerated with next 4 ids (IDBaseLine-IdProduct-IdGroup-IdNumAlert)
	Active      bool   `xorm:"active"`
	BaselineID  string `xorm:"baselineid" binding:"Required"`
	ProductID   string `xorm:"productid" binding:"Required"` //FK - > Product_devices
	GroupID     string `xorm:"groupid" binding:"Required"`
	NumAlertID  int    `xorm:"numalertid" binding:"Required"`
	Description string `xorm:"description"`
	//Alert Origin data
	InfluxDB           int64  `xorm:"influxdb" binding:"Required"`
	InfluxRP           string `xorm:"influxrp" binding:"Required"`
	InfluxMeasurement  string `xorm:"influxmeasurement" binding:"Required"`
	TagDescription     string `xorm:"tagdescription"`
	InfluxFilter       string `xorm:"influxfilter"`
	TrigerType         string `xorm:"trigertype" binding:"Required;In(DEADMAN,THRESHOLD,TREND)"` //deadman|
	IntervalCheck      string `xorm:"intervalcheck" binding:"Required"`
	OperationID        string `xorm:"operationid"`
	Field              string `xorm:"field" binding:"Required"`
	IsCustomExpression bool   `xorm:"iscustomexpression"`
	FieldDesc          string `xorm:"fielddesc"`
	ExtraData          int64  `xorm:"extradata"`
	StatFunc           string `xorm:"statfunc"`
	CritDirection      string `xorm:"critdirection"`
	Shift              string `xorm:"shift"`
	ThresholdType      string `xorm:"thresholdtype"` //Absolute/Relative
	TrendSign          string `xorm:"trendsign"`     //Positive/Negative
	//thresholds
	//CRITICAL
	ThCritDef         float64 `xorm:"th_crit_def"`
	ThCritEx1         float64 `xorm:"th_crit_ex1"`
	ThCritEx2         float64 `xorm:"th_crit_ex2"`
	ThCritRangeTimeID string  `xorm:"th_crit_rangetime_id"`

	//WARN
	ThWarnDef         float64 `xorm:"th_warn_def"`
	ThWarnEx1         float64 `xorm:"th_warn_ex1"`
	ThWarnEx2         float64 `xorm:"th_warn_ex2"`
	ThWarnRangeTimeID string  `xorm:"th_warn_rangetime_id"`

	//INFO
	ThInfoDef         float64 `xorm:"th_info_def"`
	ThInfoEx1         float64 `xorm:"th_info_ex1"`
	ThInfoEx2         float64 `xorm:"th_info_ex2"`
	ThInfoRangeTimeID string  `xorm:"th_info_rangetime_id"`

	//Grafana dashboard
	GrafanaServer      string `xorm:"grafana_server"`
	GrafanaDashLabel   string `xorm:"grafana_dash_label"`
	GrafanaDashPanelID string `xorm:"grafana_panel_id"`
	DeviceIDTag        string `xorm:"deviceid_tag"`
	DeviceIDLabel      string `xorm:"deviceid_label"`
	ExtraTag           string `xorm:"extra_tag"`
	ExtraLabel         string `xorm:"extra_label"`
	AlertExtraText     string `xorm:"alert_extra_text"`
	//Where to deploy this rule
	KapacitorID string `xorm:"kapacitorid" binding:"Required"`

	OutHTTP                 []string  `xorm:"-"` //relation with outhttpcfgs
	Modified                time.Time `xorm:"modified"`
	ServersWOLastDeployment []string  `xorm:"servers_wo_last_deployment"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (AlertIDCfg) TableName() string {
	return "alert_id_cfg"
}

// AlertEventCfg is a structure that contains relevant data about an alert event.
// The structure is intended to be JSON encoded, providing a consistent data format.
type AlertEventCfg struct {
	UID           int64         `xorm:"'uid' pk autoincr"`
	ID            string        `xorm:"id"`
	Message       string        `xorm:"message"`
	Details       string        `xorm:"details"`
	Time          time.Time     `xorm:"eventtime"`
	Duration      time.Duration `xorm:"duration"`
	Level         string        `xorm:"level"`
	PreviousLevel string        `xorm:"previousLevel"`
	//Data     models.Result `xorm:"data"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (AlertEventCfg) TableName() string {
	return "alert_event_cfg"
}

// DBConfig read from DB
type DBConfig struct {
	DeviceStat map[int64]*DeviceStatCfg
	RangeTime  map[string]*RangeTimeCfg
	Product    map[string]*ProductCfg
	Kapacitor  map[string]*KapacitorCfg
	AlertID    map[string]*AlertIDCfg
	AlertEvent map[int64]*AlertEventCfg
	Template   map[string]*TemplateCfg
	OutHTTP    map[string]*OutHTTPCfg
}

// Init initialices the DB
func Init(cfg *DBConfig) error {

	log.Debug("--------------------Initializing Config-------------------")

	log.Debug("-----------------------END Config metrics----------------------")
	return nil
}
