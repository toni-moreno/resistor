//Resistor IP
var RESISTOR_IP string

//Resistor Port
var RESISTOR_PORT string

//Instruction ID
var ID_INSTRUCTION string

//Alert extra text
var ALERT_EXTRA_TEXT string

//LB-LE sets if its a base line or extended one.
var ID_LINE string

//Product ID
var ID_PRODUCT string

//Group of alerts ID
var ID_GROUP string

//Numerical ID
var ID_NUMALERT int

//Alert Task ID
var ID_ALERT string

//Internal Task ID
var ID string

//Resistor POST URL
var http_post_url = 'http://' + RESISTOR_IP + ':' + RESISTOR_PORT + '/api/rt/kapfilter/alert/'

//InfluxDB server where stores the info
var INFLUX_BD string

//InfluxDB Retention Policy
var INFLUX_RP string

//InfluxDB Measurement
var INFLUX_MEAS string

//InfluxDB GROUP BY query
var influx_agrup = [*]

//Extra filters
var INFLUX_FILTER = lambda: TRUE

//StateChanges duration forcing alert to be sent on endpoint
var STATECHANGES_DURATION = 5m

//Interval of window node
var INTERVAL_CHECK duration

//Interval to emit values into next node
var EVERY = 1m
@EXTRADATA@@UNITDATA@
//Field to eval
var FIELD lambda

//Field desc to display on message
var FIELD_DESC= ''

//Default - CRIT Threshold
var TH_CRIT_DEF float

//Default - WARN Threshold
var TH_WARN_DEF float

//Default - INFO Threshold
var TH_INFO_DEF float

//Exception 1 - CRIT Threshold
var TH_CRIT_EX1 float

//Exception 1 - WARN Threshold
var TH_WARN_EX1 float

//Exception 1 - INFO Threshold
var TH_INFO_EX1 float

//Exception 2 - CRIT Threshold
var TH_CRIT_EX2 float

//Exception 2 - WARN Threshold
var TH_WARN_EX2 float

//Exception 3 - INFO Threshold
var TH_INFO_EX2 float

//Minimum Hour to eval the alert - CRIT
var TH_CRIT_MIN_HOUR int

//Minimum Hour to eval the alert - WARN
var TH_WARN_MIN_HOUR int

//Minimum Hour to eval the alert - INFO
var TH_INFO_MIN_HOUR int

//Maximum Hour to eval the alert - CRIT
var TH_CRIT_MAX_HOUR int

//Maximum Hour to eval the alert - WARN
var TH_WARN_MAX_HOUR int

//Maximum Hour to eval the alert - INFO
var TH_INFO_MAX_HOUR int

//Week day to eval the alert (sunday = 0) - CRIT
var DAY_WEEK_CRIT string

//Week day to eval the alert (sunday = 0) - WARN
var DAY_WEEK_WARN string

//Week day to eval the alert (sunday = 0) - INFO
var DAY_WEEK_INFO string

//Kapacitor extraTags - idTag
var idTag = 'alertID'

//Kapacitor extraTags - levelTag
var levelTag = 'level'

//Kapacitor extraFields messageField
var messageField = 'message'

//Kapacitor extraFields messageField
var durationField = 'duration'

//Grafana Server URL
var GRAFANA_SERVER = ''

//Grafana Dashboard Label 
var GRAFANA_DASH_LABEL = ''

//Grafana Dashboard Panel ID
var GRAFANA_DASH_PANELID = ''

//Device ID Tag - Core product key
var DEVICEID_TAG string 

//Device ID Label
var DEVICEID_LABEL = ''

//Extra Label
var EXTRA_LABEL = ''

//Extra Tag Key
var EXTRA_TAG = ''

//Message - Message content
var message = 'Alert  ('+ID_ALERT+') in '+DEVICEID_TAG+': {{index .Tags "'+DEVICEID_TAG+'"}} ['+EXTRA_LABEL+' : {{- index .Tags "'+EXTRA_TAG+'" -}}]. Status: {{ .Level }} ' 

//Message - Details content
var details = '''
<h3>'''+ID_ALERT+'''</h3>

{{block "alert" .}}
  <p><b> [ {{.Name}} ] - [ '''+ FIELD_DESC + ''' ] : 
  {{if eq .Level "OK" }}
  <span style="color:green;">{{ .Level }} value : {{ index .Fields "value" }} </span>
  {{else if eq .Level  "CRITICAL" }}
  <span style="color:red;">{{ .Level }} value: {{ index .Fields "value" }} </span>
  {{else if eq .Level  "WARNING" }}
  <span style="color:orange;">{{ .Level }} value: {{ index .Fields "value" }} </span>
  {{else if eq .Level  "INFO" }}
  <span style="color:blue;">{{ .Level }} value: {{ index .Fields "value" }} </span>
  {{end}}
{{end}}
<h3>TAGS</h3>
{{block "taglist" .}}
  {{range $key, $value := .Tags}}
  <li><strong>{{ $key }}</strong>: {{ $value }}</li>
  {{end}}
{{end}}
<h3>FIELDS</h3>
{{block "fieldlist" .}}
   {{range $key, $value := .Fields}}
    <li><strong>{{ $key }}</strong>: {{ $value }}</li>
   {{end}} 
{{end}}
<hr>
<h3> Extra Info </h3>
Time: {{ .Time }} 
Time Sec: {{ .Time.Second }} 
<br>
'''+ ALERT_EXTRA_TEXT + '''<br>
<b>Action Code: </b> ''' + ID_INSTRUCTION

//TICKSCRIPT:
//================
stream
    |from()
        .database(INFLUX_BD)
        .retentionPolicy(INFLUX_RP)
        .measurement(INFLUX_MEAS)
        .groupBy(influx_agrup)
      .where(INFLUX_FILTER)
    |eval(FIELD)
        .quiet()
        .as('value')
        .keep('value')@DERIVNODE@@WINDOW@
    |@FUNCTION@('value'@EXTRANODE@)
      .as('value')
    @resInjector()
        .alertId(ID_ALERT)
        .productID(ID_PRODUCT)
        .searchByTag(DEVICEID_TAG)
        .setLine(ID_LINE)
        .timeCrit(DAY_WEEK_CRIT, TH_CRIT_MAX_HOUR, TH_CRIT_MIN_HOUR)
        .timeWarn(DAY_WEEK_WARN, TH_WARN_MAX_HOUR, TH_WARN_MIN_HOUR)
        .timeInfo(DAY_WEEK_INFO, TH_INFO_MAX_HOUR, TH_INFO_MIN_HOUR)
    |alert()
        .info(lambda: if("check_info", float("value") @DIRECTION@ if("mon_exc" == 0, TH_INFO_DEF, if("mon_exc" == 1, TH_INFO_EX1, if("mon_exc" == 2, TH_INFO_EX2, 0.0))),FALSE))
        .warn(lambda: if("check_warn", float("value") @DIRECTION@ if("mon_exc" == 0, TH_WARN_DEF, if("mon_exc" == 1, TH_WARN_EX1, if("mon_exc" == 2, TH_WARN_EX2, 0.0))),FALSE))
        .crit(lambda: if("check_crit", float("value") @DIRECTION@ if("mon_exc" == 0, TH_CRIT_DEF, if("mon_exc" == 1, TH_CRIT_EX1, if("mon_exc" == 2, TH_CRIT_EX2, 0.0))),FALSE))
        .stateChangesOnly(STATECHANGES_DURATION)
        .id(ID)
        .idTag(idTag)
        .levelTag(levelTag)
        .messageField(messageField)
        .durationField(durationField)
        .message(message)
        .details(details)
        .post(http_post_url)