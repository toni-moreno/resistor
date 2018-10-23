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
	Description    string `xorm:"description text"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (DeviceStatCfg) TableName() string {
	return "device_stat_cfg"
}

// IfxServerCfg Influx server  config
type IfxServerCfg struct {
	ID          string `xorm:"'id' unique" binding:"Required"`
	URL         string `xorm:"URL" binding:"Required"`
	AdminUser   string `xorm:"adminuser"`
	AdminPasswd string `xorm:"adminpasswod"`
	Description string `xorm:"description text"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (IfxServerCfg) TableName() string {
	return "ifx_server_cfg"
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
	Description  string           `xorm:"description text"`
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
	Tags        []string `xorm:"tags text" binding:"Required"`
	Fields      []string `xorm:"fields mediumtext" binding:"Required"`
	Description string   `xorm:"description text"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (IfxMeasurementCfg) TableName() string {
	return "ifx_measurement_cfg"
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
	ID               string   `xorm:"'id' unique" binding:"Required"`
	ProductTag       string   `xorm:"producttag" binding:"Required"`
	CommonTags       []string `xorm:"commontags"`
	ExtraTags        []string `xorm:"extratags"`
	BaseLines        []string `xorm:"baselines"`
	Measurements     []string `xorm:"measurements"`
	AlertGroups      []string `xorm:"alertgroups"`
	FieldResolutions []string `xorm:"fieldresolutions"`
	Description      string   `xorm:"description text"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (ProductCfg) TableName() string {
	return "product_cfg"
}

// ProductGroupCfg Product Group Catalog Config type
type ProductGroupCfg struct {
	ID          string   `xorm:"'id' unique" binding:"Required"`
	Products    []string `xorm:"products"`
	Description string   `xorm:"description text"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (ProductGroupCfg) TableName() string {
	return "product_group_cfg"
}

// KapacitorCfg Kapacitor URL's config type
type KapacitorCfg struct {
	ID          string `xorm:"'id' unique" binding:"Required"`
	URL         string `xorm:"URL" binding:"Required"`
	Description string `xorm:"description text"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (KapacitorCfg) TableName() string {
	return "kapacitor_cfg"
}

// OperationCfg Operation instructions
type OperationCfg struct {
	ID          string `xorm:"'id' unique" binding:"Required"`
	URL         string `xorm:"URL" binding:"Required"`
	Description string `xorm:"description text"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (OperationCfg) TableName() string {
	return "operation_cfg"
}

// RangeTimeCfg Range or periods Times definition config type
type RangeTimeCfg struct {
	ID          string `xorm:"'id' unique" binding:"Required"`
	MaxHour     int    `xorm:"'max_hour' default 23"`
	MinHour     int    `xorm:"'min_hour' default 0"`
	WeekDays    string `xorm:"'weekdays' default '0123456'"`
	Description string `xorm:"description text"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (RangeTimeCfg) TableName() string {
	return "range_time_cfg"
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
	TriggerType             string    `xorm:"triggertype" binding:"Required;In(DEADMAN,THRESHOLD,TREND)"` //deadman
	StatFunc                string    `xorm:"statfunc"`
	CritDirection           string    `xorm:"critdirection"`
	TrendType               string    `xorm:"trendtype"` //Absolute/Relative
	TrendSign               string    `xorm:"trendsign"` //Positive/Negative
	FieldType               string    `xorm:"fieldtype"` //Counter/Gauge
	TplData                 string    `xorm:"tpldata mediumtext"`
	Description             string    `xorm:"description text"`
	Modified                time.Time `xorm:"modified"`
	ServersWOLastDeployment []string  `xorm:"servers_wo_last_deployment"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (TemplateCfg) TableName() string {
	return "template_cfg"
}

// EndpointCfg Alert Destination HTTP based backends config
type EndpointCfg struct {
	ID                 string   `xorm:"'id' unique" binding:"Required"`
	Type               string   `xorm:"type"`
	Description        string   `xorm:"description text"`
	URL                string   `xorm:"url"`
	Headers            []string `xorm:"headers"`
	BasicAuthUsername  string   `xorm:"basicauthusername"`
	BasicAuthPassword  string   `xorm:"basicauthpassword"`
	LogFile            string   `xorm:"logfile"`
	LogLevel           string   `xorm:"loglevel"`
	Enabled            bool     `xorm:"enabled"`
	Channel            string   `xorm:"channel"`
	SlackUsername      string   `xorm:"slackusername"`
	IconEmoji          string   `xorm:"iconemoji"`
	SslCa              string   `xorm:"sslca"`
	SslCert            string   `xorm:"sslcert"`
	SslKey             string   `xorm:"sslkey"`
	InsecureSkipVerify bool     `xorm:"insecureskipverify"`
	Host               string   `xorm:"host"`
	Port               int      `xorm:"port"`
	Username           string   `xorm:"username"`
	Password           string   `xorm:"password"`
	From               string   `xorm:"from"`
	To                 []string `xorm:"to"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (EndpointCfg) TableName() string {
	return "endpoint_cfg"
}

// AlertEndpointRel  relation between Alerts and Endpoint's type
type AlertEndpointRel struct {
	AlertID    string `xorm:"alert_id"`
	EndpointID string `xorm:"endpoint_id"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (AlertEndpointRel) TableName() string {
	return "alert_endpoint_rel"
}

// AlertIDCfg Alert Definition Config type
type AlertIDCfg struct {
	//Alert ID data
	ID          string `xorm:"'id' unique" binding:"Required"` //Autogenerated with next 4 ids (IDBaseLine-IdProduct-IdGroup-IdNumAlert)
	Active      bool   `xorm:"active"`
	BaselineID  string `xorm:"baselineid" binding:"Required"`
	ProductID   string `xorm:"productid" binding:"Required"` //FK - > Product_devices
	AlertGroup  string `xorm:"alertgroup" binding:"Required"`
	NumAlertID  int    `xorm:"numalertid" binding:"Required"`
	Description string `xorm:"description text"`
	//Alert Origin data
	InfluxDB           int64  `xorm:"influxdb" binding:"Required"`
	InfluxRP           string `xorm:"influxrp" binding:"Required"`
	InfluxMeasurement  string `xorm:"influxmeasurement" binding:"Required"`
	TagDescription     string `xorm:"tagdescription"`
	InfluxFilter       string `xorm:"influxfilter"`
	TriggerType        string `xorm:"triggertype" binding:"Required;In(DEADMAN,THRESHOLD,TREND)"` //deadman|
	IntervalCheck      string `xorm:"intervalcheck" binding:"Required"`
	AlertFrequency     string `xorm:"alertfrequency"`
	AlertNotify        int    `xorm:"alertnotify"`
	OperationID        string `xorm:"operationid"`
	Field              string `xorm:"field" binding:"Required"`
	IsCustomExpression bool   `xorm:"iscustomexpression"`
	FieldDesc          string `xorm:"fielddesc"`
	ExtraData          int64  `xorm:"extradata"`
	StatFunc           string `xorm:"statfunc"`
	CritDirection      string `xorm:"critdirection"`
	Shift              string `xorm:"shift"`
	TrendType          string `xorm:"trendtype"` //Absolute/Relative
	TrendSign          string `xorm:"trendsign"` //Positive/Negative
	FieldType          string `xorm:"fieldtype"` //Counter/Gauge
	Rate               bool   `xorm:"rate"`
	FieldResolution    string `xorm:"fieldresolution"`
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
	ProductTag         string `xorm:"producttag"`
	DeviceIDLabel      string `xorm:"deviceid_label"`
	ExtraTag           string `xorm:"extra_tag"`
	ExtraLabel         string `xorm:"extra_label"`
	AlertExtraText     string `xorm:"alert_extra_text"`
	IDTag              string `xorm:"idtag"`
	//Where to deploy this rule
	KapacitorID string `xorm:"kapacitorid" binding:"Required"`

	Endpoint                []string  `xorm:"-"` //relation with endpointcfgs
	Modified                time.Time `xorm:"modified"`
	ServersWOLastDeployment []string  `xorm:"servers_wo_last_deployment"`
	LastDeploymentTime      time.Time `xorm:"last_deployment_time"`
}

//AlertIDCfgJSON structure with data info in json
type AlertIDCfgJSON struct {
	//Alert ID data
	ID           string `json:"alertid"` //Autogenerated with next 4 ids (IDBaseLine-IdProduct-IdGroup-IdNumAlert)
	BaselineID   string `json:"baselineid"`
	ProductID    string `json:"productid"` //FK - > Product_devices
	ProductGroup string `json:"productgroup,omitempty"`
	AlertGroup   string `json:"alertgroup"`
	//Alert Origin data
	InfluxDBName      string `json:"influxdbname"`
	InfluxRP          string `json:"influxrp"`
	InfluxMeasurement string `json:"influxmeasurement"`
	InfluxFilter      string `json:"influxfilter,omitempty"`
	TriggerType       string `json:"triggertype"` //deadman|
	IntervalCheck     string `json:"intervalcheck"`
	AlertFrequency    string `json:"alertfrequency,omitempty"`
	OperationID       string `json:"operationid,omitempty"`
	Field             string `json:"field"`
	StatFunc          string `json:"statfunc,omitempty"`
	CritDirection     string `json:"critdirection,omitempty"`
	Shift             string `json:"shift,omitempty"`
	TrendType         string `json:"trendtype,omitempty"` //Absolute/Relative
	TrendSign         string `json:"trendsign,omitempty"` //Positive/Negative
	FieldType         string `json:"fieldtype,omitempty"` //Counter/Gauge
	FieldResolution   string `json:"fieldresolution,omitempty"`
	//thresholds
	//APPLIED
	ThCrit float64 `json:"th_crit"`
	ThWarn float64 `json:"th_warn"`
	ThInfo float64 `json:"th_info"`

	AlertExtraText string `json:"alert_extra_text,omitempty"`
	IDTag          string `json:"idtag,omitempty"`
	//Where to deploy this rule
	KapacitorID string `json:"kapacitorid"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (AlertIDCfg) TableName() string {
	return "alert_id_cfg"
}

// AlertEventHist is a structure that contains relevant data about an alert event.
// The structure is intended to be JSON encoded, providing a consistent data format.
type AlertEventHist struct {
	ID              int64         `xorm:"'id' pk"`
	CorrelationID   string        `xorm:"correlationid"`
	AlertID         string        `xorm:"alertid"`
	ProductID       string        `xorm:"productid"`
	ProductTagValue string        `xorm:"producttagvalue"`
	Field           string        `xorm:"field"`
	Message         string        `xorm:"message"`
	Details         string        `xorm:"details text"`
	FirstEventTime  time.Time     `xorm:"firsteventtime"`
	EventTime       time.Time     `xorm:"eventtime"`
	Duration        time.Duration `xorm:"duration"`
	Level           string        `xorm:"level"`
	PreviousLevel   string        `xorm:"previousLevel"`
	Tags            []string      `xorm:"tags"`
	Value           float64       `xorm:"value"`
	MonExc          string        `xorm:"mon_exc"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (AlertEventHist) TableName() string {
	return "alert_event_hist"
}

// AlertEvent is a structure that contains relevant data about an alert event.
// The structure is intended to be JSON encoded, providing a consistent data format.
type AlertEvent struct {
	ID              int64         `xorm:"'id' pk autoincr"`
	CorrelationID   string        `xorm:"correlationid"`
	AlertID         string        `xorm:"alertid"`
	ProductID       string        `xorm:"productid"`
	ProductTagValue string        `xorm:"producttagvalue"`
	Field           string        `xorm:"field"`
	Message         string        `xorm:"message"`
	Details         string        `xorm:"details text"`
	FirstEventTime  time.Time     `xorm:"firsteventtime"`
	EventTime       time.Time     `xorm:"eventtime"`
	Duration        time.Duration `xorm:"duration"`
	Level           string        `xorm:"level"`
	PreviousLevel   string        `xorm:"previousLevel"`
	Tags            []string      `xorm:"tags"`
	Value           float64       `xorm:"value"`
	MonExc          string        `xorm:"mon_exc"`
}

// TableName go-xorm way to set the Table name to something different to "alert_h_t_t_p_out_rel"
func (AlertEvent) TableName() string {
	return "alert_event"
}

// AlertEventsSummary is a structure that contains summary data about alert events.
type AlertEventsSummary struct {
	Level string `json:"level"`
	Num   int    `json:"num"`
}

// DBConfig read from DB
type DBConfig struct {
	DeviceStat        map[int64]*DeviceStatCfg
	Operation         map[string]*OperationCfg
	RangeTime         map[string]*RangeTimeCfg
	Product           map[string]*ProductCfg
	Kapacitor         map[string]*KapacitorCfg
	AlertID           map[string]*AlertIDCfg
	AlertEventHistMap map[int64]*AlertEventHist
	AlertEventMap     map[int64]*AlertEvent
	Template          map[string]*TemplateCfg
	Endpoint          map[string]*EndpointCfg
}

// Init initialices the DB
func Init(cfg *DBConfig) error {

	log.Debug("--------------------Initializing Config-------------------")

	log.Debug("-----------------------END Config metrics----------------------")
	return nil
}
