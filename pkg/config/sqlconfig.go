package config

//Real Time Filtering by device/alertid/or other tags

type DeviceStatCfg struct {
	ID             int64  `xorm:"'id' pk autoincr"`
	Order          int64  `xorm:"order"`
	DeviceID       string `xorm:"deviceid" binding:"Required"`
	AlertID        string `xorm:"alertid" binding:"Required"`
	Exception      int64  `xorm:"exception"`
	Active         bool   `xorm:"active"`
	BaseLine       string `xorm:"baseline"`
	FilterTagKey   string `xorm:"filterTagKey"`
	FilterTagValue string `xorm:"filterTagValue"`
	Description    string `xorm:"description"`
}

type ProductCfg struct {
	ID          string   `xorm:"'id' unique" binding:"Required"`
	CommonTags  []string `xorm:"commontags"`
	Description string   `xorm:"description"`
}

type KapacitorCfg struct {
	ID          string `xorm:"'id' unique" binding:"Required"`
	URL         string `xorm:"URL" binding:"Required"`
	Description string `xorm:"description"`
}

type RangeTimeCfg struct {
	ID          string `xorm:"'id' unique" binding:"Required"`
	MaxHour     int    `xorm:"'max_hour' default 23"`
	MinHour     int    `xorm:"'min_hour' default 0"`
	WeeKDays    string `xorm:"'weekdays' default '0123456'"`
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

type TemplateCfg struct {
	ID            string `xorm:"'id' unique" binding:"Required"`
	TrigerType    string `xorm:"trigertype" binding:"Required;In(DEADMAN,THRESHOLD,TREND)"` //deadman
	StatFunc      string `xorm:"statfunc"`
	CritDirection string `xorm:"critdirection"`
	ThresholdType string `xorm:"thresholdtype"` //Absolute/Relative
	TplData       string `xorm:"tpldata"`
	Description   string `xorm:"description"`
}

type OutHTTPCfg struct {
	ID          string   `xorm:"'id' unique" binding:"Required"`
	Url         string   `xorm:"url" binding:"Required"`
	Headers     []string `xorm:"headers"`
	AlertTpl    string   `xorm:"alert_tpl"`
	Description string   `xorm:"description"`
}

// SnmpDevMGroups Mgroups defined on each SnmpDevice
type AlertHTTPOutRel struct {
	AlertID   string `xorm:"alert_id"`
	HTTPOutID string `xorm:"http_out_id"`
}

type AlertIdCfg struct {
	//Alert ID data
	ID          string `xorm:"'id' unique" binding:"Required"` //Autogenerated with next 4 ids (IDBaseLine-IdProduct-IdGroup-IdNumAlert)
	BaselineID  string `xorm:"baselineid" binding:"Required"`
	ProductID   string `xorm:"productid" binding:"Required"` //FK - > Product_devices
	GroupID     string `xorm:"groupid" binding:"Required"`
	NumAlertID  int    `xorm:"numalertid" binding:"Required"`
	Description string `xorm:"description"`
	//Alert Origin data
	InfluxDB          string `xorm:"influxDB" binding:"Required"`
	InfluxRP          string `xorm:"influxRP" binding:"Required"`
	InfluxMeasurement string `xorm:"influxmeasurement" binding:"Required"`
	TagDescription    string `xorm:"tagdescription"`
	InfluxFilter      string `xorm:"influxfilter"`
	TrigerType        string `xorm:"trigertype" binding:"Required;In(DEADMAN,THRESHOLD,TREND)"` //deadman|
	IntervalCheck     string `xorm:"intervalcheck" binding:"Required"`
	OperationID       string `xorm:"operationid"`
	Field             string `xorm:"field" binding:"Required"`
	StatFunc          string `xorm:"statfunc" binding:"Required"`
	CritDirection     string `xorm:"critdirection" binding:"Required"`
	Shift             int64  `xorm:"shift"`
	ThresholdType     string `xorm:"thresholdtype"` //Absolute/Relative
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
	//Where to deploy this rule
	KapacitorID string `xorm:"kapacitorid" binding:"Required"`

	OutHTTP []string `xorm:"-"` //relation between alertIDcfgs
}

// SQLConfig read from DB

type SQLConfig struct {
	DeviceStat map[int64]*DeviceStatCfg
	RangeTime  map[string]*RangeTimeCfg
	Product    map[string]*ProductCfg
	Kapacitor  map[string]*KapacitorCfg
	AlertID    map[string]*AlertIdCfg
	Template   map[string]*TemplateCfg
	OutHTTP    map[string]*OutHTTPCfg
}

func Init(cfg *SQLConfig) error {

	log.Debug("--------------------Initializing Config-------------------")

	log.Debug("-----------------------END Config metrics----------------------")
	return nil
}
