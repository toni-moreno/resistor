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

//Extra filters
var INFLUX_FILTER = lambda: TRUE

//StateChanges duration forcing alert to be sent on endpoint
var STATECHANGES_DURATION = 5m

//Interval of window node
var INTERVAL_CHECK duration

//Field to eval
var FIELD lambda

//Interval to emit values into next node
var every = 1m

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
var DEVICEID_TAG = '*' 

//InfluxDB GROUP BY query
var influx_agrup = [DEVICEID_TAG]

//Device ID Label
var DEVICEID_LABEL = ''

//Extra Label
var EXTRA_LABEL = ''

//Extra Tag Key
var EXTRA_TAG = ''

//Message - Message content
var message = 'Alert ('+ID_ALERT+'). '+DEVICEID_TAG+': {{index .Tags "'+DEVICEID_TAG+'"}} ['+EXTRA_LABEL+' : {{- index .Tags "'+EXTRA_TAG+'" -}}]. Status: {{ .Level }}' 

//Message - Details content
var details = '''
<h3>'''+ID_ALERT+'''</h3>

{{block "alert" .}}
  <p><b> [ {{.Name}} ] : 
  {{if eq .Level "OK" }}
    <span style="color:green;">{{ .Level }} value : {{ index .Fields "emitted" }} </span>
  {{else if eq .Level  "CRITICAL" }}
    <span style="color:red;">{{ .Level }} value: {{ index .Fields "emitted" }} </span>
  {{else if eq .Level  "WARNING" }}
    <span style="color:orange;">{{ .Level }} value: {{ index .Fields "emitted" }} </span>
  {{else if eq .Level  "INFO" }}
    <span style="color:blue;">{{ .Level }} value: {{ index .Fields "emitted" }} </span>
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
        .keep('value')
    |stats(INTERVAL_CHECK)
        .align()
    |derivative('emitted')
        .unit(INTERVAL_CHECK)
        .nonNegative()
    @resInjector()
        .alertId(ID_ALERT)
        .productID(ID_PRODUCT)
        .searchByTag(DEVICEID_TAG)
        .setLine(ID_LINE)
    |alert()
        .info(lambda: if("check_info", float("emitted") < if("mon_exc" == 0, 1.0, if("mon_exc" == 1, 1.0, if("mon_exc" == 2, 1.0, 0.0))),FALSE))
        .warn(lambda: if("check_warn", float("emitted") < if("mon_exc" == 0, 1.0, if("mon_exc" == 1, 1.0, if("mon_exc" == 2, 1.0, 0.0))),FALSE))
        .crit(lambda: if("check_crit", float("emitted") < if("mon_exc" == 0, 1.0, if("mon_exc" == 1, 1.0, if("mon_exc" == 2, 1.0, 0.0))),FALSE))
        .stateChangesOnly(STATECHANGES_DURATION)
        .id(ID)
        .idTag(idTag)
        .levelTag(levelTag)
        .messageField(messageField)
        .durationField(durationField)
        .message(message)
        .details(details)
        .post(http_post_url)