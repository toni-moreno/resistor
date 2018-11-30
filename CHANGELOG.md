# CHANGELOG.md

## v 0.6.5  (30/11/2018)

### New features.

* New bulk actions added for lists on several components.
* Added new button on Export Data to do a Bulk Export.
  * This export also includes the export of InfluxDB catalog.
* Tables from InfluxDB catalog can be imported from a JSON file.
* Imported datetime added in all tables which can be imported.
* The message shown to the user when trying to create a duplicated alert has been modified.
* Added more information into the body for email endpoint.
* Tooltip for Trend Sign on Alert Definition modified.
* New validator for WeekDays field on Range Time component.
* Order on left menu modified to put Influx Catalog under External Server Config.

### fixes

* Fix in resinjector UDF to use Range Times with start and end times in different days.

### breaking changes


## v 0.6.4  (22/11/2018)

### New features.

* Changes to use TEXT IDs for InfluxDB catalog tables.

### fixes


### breaking changes

* -- Execute the following steps in order to mantain your data:

  * -- 1 -- stop resistor
  * -- 2 -- dump your database
  * -- 3 -- rename modified tables to backup temp tables
    * ALTER TABLE ifx_db_cfg RENAME TO temp_ifx_db_cfg;
    * ALTER TABLE ifx_measurement_cfg RENAME TO temp_ifx_measurement_cfg;
    * ALTER TABLE ifx_db_meas_rel RENAME TO temp_ifx_db_meas_rel;
    * ALTER TABLE alert_id_cfg RENAME TO temp_alert_id_cfg;
  * -- 4 -- add new id columns
    * ALTER TABLE temp_ifx_db_cfg ADD COLUMN id_new TEXT;
    * ALTER TABLE temp_ifx_measurement_cfg ADD COLUMN id_new TEXT;
    * ALTER TABLE temp_ifx_db_meas_rel ADD COLUMN ifxdbid_new TEXT;
    * ALTER TABLE temp_ifx_db_meas_rel ADD COLUMN ifxmeasid_new TEXT;
    * ALTER TABLE temp_alert_id_cfg ADD COLUMN influxdb_new TEXT;

  * -- 5 -- set value on new id columns
    * set sql_mode=PIPES_AS_CONCAT; -- Only for MySQL
    * UPDATE temp_ifx_db_cfg SET id_new = ifxserver || '-' || name;
    * UPDATE temp_ifx_measurement_cfg SET id_new = (SELECT ifxserver || '-' || temp_ifx_db_cfg.name || '-' || temp_ifx_measurement_cfg.name FROM temp_ifx_db_cfg, temp_ifx_db_meas_rel WHERE temp_ifx_db_cfg.id = ifxdbid AND temp_ifx_measurement_cfg.id = temp_ifx_db_meas_rel.ifxmeasid);
    * UPDATE temp_ifx_db_meas_rel SET ifxdbid_new = (SELECT id_new FROM temp_ifx_db_cfg WHERE temp_ifx_db_cfg.id = ifxdbid);
    * UPDATE temp_ifx_db_meas_rel SET ifxmeasid_new = (SELECT id_new FROM temp_ifx_measurement_cfg WHERE temp_ifx_measurement_cfg.id = ifxmeasid);
    * UPDATE temp_alert_id_cfg SET influxdb_new = (SELECT id_new FROM temp_ifx_db_cfg WHERE temp_ifx_db_cfg.id = influxdb);
  * -- 6 -- DROP INDEXES;
    * DROP INDEX UQE_ifx_db_cfg_ifxdb ON temp_ifx_db_cfg; -- Only for MySQL
    * DROP INDEX UQE_alert_id_cfg_id ON temp_alert_id_cfg; -- Only for MySQL
    * DROP INDEX UQE_ifx_db_cfg_ifxdb; -- Only for SQLite
    * DROP INDEX UQE_alert_id_cfg_id; -- Only for SQLite
  * -- 7 -- start resistor
  * -- 8 -- modified tables will be created again
  * -- 9 -- stop resistor
  * -- 10 -- copy backup content to modified tables
    * INSERT INTO ifx_db_cfg (id, name, ifxserver, retention, description) SELECT id_new, name, ifxserver, retention, description FROM temp_ifx_db_cfg;
    * INSERT INTO ifx_measurement_cfg (id, name, tags, fields, description) SELECT id_new, name, tags, fields, description FROM temp_ifx_measurement_cfg;
    * INSERT INTO ifx_db_meas_rel (ifxdbid, ifxmeasid, ifxmeasname) SELECT ifxdbid_new, ifxmeasid_new, ifxmeasname FROM temp_ifx_db_meas_rel;
    * INSERT INTO alert_id_cfg (id, active, baselineid, productid, alertgroup, numalertid, description, influxdb, influxrp, influxmeasurement, tagdescription, influxfilter, triggertype, intervalcheck, alertfrequency, alertnotify, operationid, field, iscustomexpression, fielddesc, extradata, statfunc, critdirection, shift, trendtype, trendsign, fieldtype, rate, fieldresolution, th_crit_def, th_crit_ex1, th_crit_ex2, th_crit_rangetime_id, th_warn_def, th_warn_ex1, th_warn_ex2, th_warn_rangetime_id, th_info_def, th_info_ex1, th_info_ex2, th_info_rangetime_id, grafana_server, grafana_dash_label, grafana_panel_id, producttag, deviceid_label, extra_tag, extra_label, alert_extra_text, kapacitorid, modified, servers_wo_last_deployment, last_deployment_time) SELECT id, active, baselineid, productid, alertgroup, numalertid, description, influxdb_new, influxrp, influxmeasurement, tagdescription, influxfilter, triggertype, intervalcheck, alertfrequency, alertnotify, operationid, field, iscustomexpression, fielddesc, extradata, statfunc, critdirection, shift, trendtype, trendsign, fieldtype, rate, fieldresolution, th_crit_def, th_crit_ex1, th_crit_ex2, th_crit_rangetime_id, th_warn_def, th_warn_ex1, th_warn_ex2, th_warn_rangetime_id, th_info_def, th_info_ex1, th_info_ex2, th_info_rangetime_id, grafana_server, grafana_dash_label, grafana_panel_id, producttag, deviceid_label, extra_tag, extra_label, alert_extra_text, kapacitorid, modified, servers_wo_last_deployment, last_deployment_time FROM temp_alert_id_cfg;
  * -- 11 -- DROP TEMP TABLES
    * DROP TABLE temp_ifx_db_cfg;
    * DROP TABLE temp_ifx_measurement_cfg;
    * DROP TABLE temp_ifx_db_meas_rel;
    * DROP TABLE temp_alert_id_cfg;
  * -- 12 -- start resistor


## v 0.6.3  (09/11/2018)

### New features.

* Changes on Templates component:
  * The user is informed on frontend if an error occurs deploying templates.
* Changes on Alerts Definition component:
  * The user is informed on frontend if an error occurs deploying tasks.
* Changes on Alert Events and Alert Events History components:
  * The user need to be signed in again to get the list of alert events or the detail of an alert event.
  * New information about Alert Events and Alert Events History accessible by curl.
    * curl http://localhost:6090/api/rt/alertevent/list/
    * curl http://localhost:6090/api/rt/alertevent/byid/numeric_id
    * curl http://localhost:6090/api/rt/alerteventhist/list/
    * curl http://localhost:6090/api/rt/alerteventhist/byid/numeric_id
* Changes on alerting to Slack Endpoints.
  * Fields of type Time are displayed in the local timezone of the person seeing the message.

### fixes

* Changes on Alerts Definition component:
  * New variable ID_ALERT for Kapacitor tasks. Needed for resInjector UDF.
  * Fixed error when an alert is updated and there isn't network connection.
* Changes on Alert Events and Alert Events History components:
  * FirstEventTime can not be empty.
* Changes on alerting to Slack Endpoints.
  * Fixed errors with malformed dashboard URL.
* Changes on resInjector UDF:
  * The UDF does not stop if ProductTag does not arrive in Point.

### breaking changes

* This version must be installed before the deployment of new version of templates.


## v 0.6.2  (31/10/2018)

### New features.


### fixes

* Changes on Alert Events and Alert Events History components:
  * AlertID column shows the ID of the Resistor Alert again, not the ID received from the Kapacitor AlertNode.

### breaking changes


## v 0.6.1  (30/10/2018)

### New features.

* New parameter 'correlationidtemplate' on [alerting] section of config file.
  * Template for constructing a unique ID for a given alert. The ID will be: taskName + "|" + correlation_id_template
  * Example: correlationidtemplate="{{range $key, $value := .Tags}} {{ $key }}:{{ $value }}.{{end}}"
* Changes on Alert Definition component: 
  * Added new link to the Dashboard URL.
  * UID Tag field filled in with tags from measurement.
  * Unsubscriptions done if another item from menu is clicked when a get response is pending.
* Icon for unselect on single-select fields modified to trash-icon.

### fixes

* Changes to fix errors on Dashboard URL construction.
* Changes to fix warnings with disabled fields on Alert Definition component.
* Changes to fix error when column for multiselection is shown.

### breaking changes


## v 0.6.0  (23/10/2018)

### New features.

* New information about Resistor version accessible by curl.
  * curl http://localhost:6090/api/rt/agent/info/version/
* New information about Alert Events and Alert Events History accessible by curl.
  * curl http://localhost:6090/api/rt/alertevent/groupbylevel/
  * curl http://localhost:6090/api/rt/alerteventhist/groupbylevel/
* New endpoint type email.
* Information of alerts sent to Slack improved.

### fixes

### breaking changes


## v 0.5.11  (19/10/2018)

### New features.

* New component Operation added.
  * This component is used to define the operation instructions for the operator related to the alerts.
* Field OperationID on Alert Definition changed to select type.
* New json fields included into the data sent to the HTTP Post endpoints.
  * resistor-operationid: operationid related to the alert.
  * resistor-operationurl: operationurl related to the alert.
  * resistor-dashboardurl: dashboardurl related to the alert.
* New field IDTag on Alert.
  * This field and its value are included into the data sent to the HTTP Post endpoints (json fields resistor-id-tag-name and resistor-id-tag-value).
  * If the field is empty, ProductTag field is used to sent data to the HTTP Post endpoints.  
* The user is informed on frontend if an error occurs when deploying a task on kapacitor server.
* The user is informed on frontend if an error occurs when deploying a template on kapacitor server.
* The deployment of tasks and templates is always done when the user press the 'Deploy Item' button on lists.
* Blocker added to Alert Event and Alert Event History components.
* Added an unselect button for single select fields.
* Added the possibility of Enable/Disable for multiple items in Alert Definition, Alerting Endpoints and Device Stats components.
* Errors 400 and 422 shown on frontend.
* Some logs improved.
* About modified.
* Labels changed on Alert Definition and Product components.
* Status buttons row on Alert Events and Alert Events History components improved.
* Evaluation Period field on Alert Definition component changed to select type and filled in depending on values from selected product.
* When deleting one operation, the related alerts are not deleted.

### fixes

* Changes to fix #88
  * New parameter resistorurl added on config file
* Changes to fix error when sending to slack endpoint with empty SSL fields and InsecureSkipVerify=true
* Changes to fix error if AlertNotify field on Alert Definition is empty.
* Changes to fix error if Endpoint field on Alert Definition is empty.
* Changes to fix error if there are no items when selecting a component on Export Data.
* Changes to fix error when entering components with column filtering if a previous column filtering has been done.

### breaking changes

* Changes to fix error when defining a template in MySQL
  * Column tpldata on table template_cfg modified to MEDIUMTEXT
    * Execute the following sql:
      * ALTER TABLE template_cfg MODIFY COLUMN tpldata MEDIUMTEXT;
* Changes to fix error when importing a measurement with a very large number of fields in MySQL
  * Column fields on table ifx_measurement_cfg modified to MEDIUMTEXT
    * Execute the following sql:
      * ALTER TABLE ifx_measurement_cfg MODIFY COLUMN fields MEDIUMTEXT;
* Changes to avoid error when importing a measurement with a very large number of tags in MySQL
  * Column tags on table ifx_measurement_cfg modified to TEXT
    * Execute the following sql:
      * ALTER TABLE ifx_measurement_cfg MODIFY COLUMN tags TEXT;
* Changes to fix error when inserting an alert event in MySQL
  * Column details on tables alert_event and alert_event_hist modified to TEXT
    * Execute the following sqls:
      * ALTER TABLE alert_event MODIFY COLUMN details TEXT;
      * ALTER TABLE alert_event_hist MODIFY COLUMN details TEXT;
* Changes to avoid error when inserting an entity with a very large description in MySQL
  * Column description on all tables modified to TEXT
    * Execute the following sqls:
      * ALTER TABLE alert_id_cfg MODIFY COLUMN description TEXT;
      * ALTER TABLE device_stat_cfg MODIFY COLUMN description TEXT;
      * ALTER TABLE endpoint_cfg MODIFY COLUMN description TEXT;
      * ALTER TABLE ifx_db_cfg MODIFY COLUMN description TEXT;
      * ALTER TABLE ifx_measurement_cfg MODIFY COLUMN description TEXT;
      * ALTER TABLE ifx_server_cfg MODIFY COLUMN description TEXT;
      * ALTER TABLE kapacitor_cfg MODIFY COLUMN description TEXT;
      * ALTER TABLE product_cfg MODIFY COLUMN description TEXT;
      * ALTER TABLE product_group_cfg MODIFY COLUMN description TEXT;
      * ALTER TABLE range_time_cfg MODIFY COLUMN description TEXT;
      * ALTER TABLE template_cfg MODIFY COLUMN description TEXT;

## v 0.5.10  (08/10/2018)

### New features.

### fixes

* Changes to fix error with empty ID of Alerts.

### breaking changes


## v 0.5.9  (05/10/2018)

### New features.

### fixes

* Changes to fix not found error when Kapacitor post alert to Resistor.
* Changes to fix error when editing consecutively alerts with different products.

### breaking changes


## v 0.5.8  (03/10/2018)
### New features.
* Changes on Templates component:
    * New field 'FieldType'. Now you can have templates of type COUNTER or GAUGE.
* Changes on Alert Definition component:
    * New field 'FieldType'. Now you can have alerts for fields of type COUNTER or GAUGE.
    * New field 'AlertFrequency'. Indicates the interval used to emit alerts.
    * New field 'AlertNotify'. Indicates the number of cycles of AlertFrequency to force an alert is sent to endpoint.
    * New field 'FieldResolution'. Indicates the time unit for the derivative node used in alerts for counter fields.
    * New field 'Rate'. Used in alerts for counter fields. If true, 1s is used as unit time; else you can choose the unit with 'FieldResolution' field.
* Changes on Product component:
    * New field 'FieldResolutions' to indicate the list of possible time units for the derivative node used in counter templates.
* Changes on Alert Events component:
    * New column 'FirstEventTime' added.
* Changes on Alert Events History component:
    * New column 'FirstEventTime' added.

### fixes
* Changes on Import Data component:
    * Fixed error when importing catalog from InfluxDB on MySQL database.
	* If an error occurs when importing catalog, the spinner component is hidden.
* Changes on lists:
    * Fixed error: sometimes the 'undefined' word was shown.

### breaking changes
* New field fieldtype on table template_cfg.
    * Execute the following sql to mantain the previous behaviour:
        * UPDATE template_cfg SET fieldtype = 'GAUGE';
* New field fieldtype on table alert_id_cfg.
    * Execute the following sql to mantain the previous behaviour:
        * UPDATE alert_id_cfg SET fieldtype = 'GAUGE';
* New field alertfrequency on table alert_id_cfg.
    * Execute the following sql to mantain the previous behaviour:
        * UPDATE alert_id_cfg SET alertfrequency = '1m';


# v 0.5.7  (28/09/2018)
### New features.
* Changes on Alert Events component:
    * Component 'Alert Events' has been divided into components 'Alert Events' and 'Alert Events History'.
    * On 'Alert Events' component only the last alert event of each CorrelationID is shown.
	* On 'Alert Events History' component all the previous alert events of each CorrelationID are shown.
	* When a new alert event arrives it's added to the alert_event table and the previous alert events related are moved to the alert_event_hist table.
	* Also a clean process is executed periodically to move to history table the alerts with status OK.
	* Buttons for filtering by Level have been added.
    * Set Alert Events component as initial component when login.
* New parameter 'cleanperiod' on new section [alerting] of config file.
    * Period used to move alert events with status OK from alert_event to alert_event_hist.
    * Example: cleanperiod = "3m"
* Changes on Alerting Endpoints component:
    * SlackEnabled field has been changed to Enabled and it's used for all endpoints.
    * 'Triggered by' information added to message for Slack.
* Changes on Alert Definition component:
    * InfluxFilter field: Link to lambda expressions explanation added.
* Changes on Device Stats component:
    * BaseLine placed before AlertID.
* Parameter 'proxyurl' on config file moved from [http] section to new section [endpoints].
    * Example: proxyurl = "http://proxyIP:proxyPort"

### fixes
* Changes on Alert Definition component:
    * IntervalCheck field: Fixed the regular expression to check the data has a valid format.

### breaking changes
* Field slackenabled on table endpoint_cfg changed to enabled.
    * Execute the following sqls to mantain the previous behaviour:
        * UPDATE endpoint_cfg SET enabled = 1;
        * UPDATE endpoint_cfg SET enabled = slackenabled WHERE type = 'slack';
* New table alert_event created. This table is related with alert_event_hist table and it has the id column as PK and autoincrement.
    * Execute the following sql to mantain the previous behaviour:
        * UPDATE SQLITE_SEQUENCE SET seq = (SELECT seq FROM SQLITE_SEQUENCE WHERE name = 'alert_event_hist') WHERE name = 'alert_event';


# v 0.5.6  (20/09/2018)
### New features.
* New parameter 'proxyurl' added to config file to use it on Alerting Endpoints, if needed.
    * Example: proxyurl = "http://proxyIP:proxyPort"
* Changes on Alert Definition component:
    * InfluxFilter field: New information added on tooltip and placeholder added.
    * Columns on list reordered.
* Changes on Alert Events component:
    * New column on list: MonExc (Monitoring Exception applied).
    * Columns on list reordered.
* Changes on Influx DB Servers component:
    * Column AdminPasswd removed from list.
* Changes on Product component:
    * Measurements field: Sorted alphabetically.
    * Tags field: Sorted alphabetically.
* Changes on Templates component:
    * Columns on list reordered.

### fixes
* Changes on Templates component:
    * TrendType field: Only shown when TriggerType is Trend.

### breaking changes


# v 0.5.5  (17/09/2018)
### New features.

### fixes
* Management of errors on communication with endpoints improved.

### breaking changes


# v 0.5.4  (14/09/2018)
### New features.
* Vars and ExecutionStats of kapacitor tasks changed to json strings.

### fixes

### breaking changes


# v 0.5.3  (14/09/2018)
### New features.
* ProductGroup added in TaskAlertInfo sent to logfile and httppost.

### fixes
* Fixed error with headers sent to httppost.

### breaking changes


# v 0.5.2  (14/09/2018)
### New features.
* Thresholds with their value added in TaskAlertInfo sent to logfile and httppost.

### fixes

### breaking changes


# v 0.5.1  (14/09/2018)
### New features.
* Add correlationid to taskAlertInfo sent to logfile and httppost.

### fixes
* Fixed error when showing Alert Definition List.

### breaking changes


# v 0.5.0  (14/09/2018)
### New features.
* 'Enable/Disable Edit' button on filtering section moved to left, text changed to 'Show/Hide multiselect', text removed and shown as tooltip.
* First and Last buttons shown on pagination section.
* Selector fields modified to include the possibility of adding custom items on list.
* Changes on Device Stats component:
    * Selector for fields: ProductID, AlertID, DeviceID, BaseLine and FilterTagKey modified to include the possibility of adding custom items on list.
    * Some tooltips have been modified.
* Changes on Alert Events component:
    * New columns on list: ProductID, ProductTagValue, Field, Tags and Value.
    * Datetime format applied on column Time.
    * Alert Events sorted by ID desc on init.
    * Added the possibility of single and multiple deletion.
    * Management of filter column modified to add filter column field on each column.
    * New field included indicating the time of the last refresh.
    * New button included to refresh the list.
* Changes on Kapacitor Tasks component:
    * Set Kapacitor Tasks component as initial component when login.
    * Alerts not deployed on kapacitor shown on kapacitor tasks list.
    * New column NumErrors included on list.
    * Datetime format applied on columns Created, Modified and LastEnabled.
    * Management of filter column modified to add filter column field on each column.
    * Columns reordered.
    * New field included indicating the time of the last refresh.
    * New button included to refresh the list.
* Changes on Alerting Endpoints component:
    * Configuration of HTTP Post endpoint modified to take into account Headers and BasicAuth form fields.
    * Form fields added to configure Slack endpoint.
    * More information sent to logfile and httppost.
* Changes on Alert Definition component:
    * 0 not allowed on NumAlertID.
    * LastDeploymentTime added to Alerts list.

### fixes
* Fixed error when deleting a range time related to an alert.

### breaking changes
* Table alert_event_cfg changed to alert_event_hist. Some fields of the table also changed.
* Execute the following sql if you want to copy data from alert_event_cfg to alert_event_hist:
    * insert into alert_event_hist (id, alertid, message, details, eventtime, duration, level, previousLevel) select uid, id, message, details, eventtime, duration, level, previousLevel from alert_event_cfg;
* Execute the following sql to drop old table alert_event_cfg:
    * DROP TABLE alert_event_cfg;
* Table out_http_cfg changed to endpoint_cfg.
* Table alert_http_out_rel changed to alert_endpoint_rel.
* Execute the following sqls if you want to copy data from old tables to new tables:
    * insert into endpoint_cfg (id, type, description, url, headers, basicauthusername, basicauthpassword, logfile, loglevel, slackenabled, channel, slackusername, iconemoji, sslca, sslcert, sslkey, insecureskipverify) select id, type, description, url, headers, basicauthusername, basicauthpassword, logfile, loglevel, slackenabled, channel, slackusername, iconemoji, sslca, sslcert, sslkey, insecureskipverify from out_http_cfg;
    * insert into alert_endpoint_rel (alert_id, endpoint_id) select alert_id, http_out_id from alert_http_out_rel;
* Execute the following sql to drop old tables out_http_cfg and alert_http_out_rel:
    * DROP TABLE out_http_cfg;
    * DROP TABLE alert_http_out_rel;

# v 0.4.0  (23/08/2018)
### New features.
* Changes for defining alerts depending on product and some errors fixed.
    * 'InfluxMeasurement' field filled in with data depending on selected product.
    * 'BaselineID' field filled in with data depending on selected product.
    * 'GroupID' field filled in with data depending on selected product and label changed to 'AlertGroup'.
    * 'InfluxDB' field filled in with data depending on selected measurement.
    * 'DeviceIDTag' field filled in with data depending on selected product and label changed to 'ProductTag'.
* 'Alerts' label changed to 'Alert Definitions'.
* Changes for defining device stats (alert exceptions) based on product.
* Changes for Alerting Endpoints refactoring.
    * Now the configuration for httppost and logging is done with several form fields.
    * The configuration for slack still is done with a form field in JSON format.
* The size of TplData textarea has been increased.
* The Import Data button for Influx DB Servers is disabled on create mode.
* If an alert is created or modified with Active=true, the related kapacitor task is created or modified disabled, then it's enabled.
    * This is done in order to Kapacitor applies new values to task.
* Changed log level message when a kapacitor task is not found: from error to debug. This message is also logged when a kapacitor task is been created.

### fixes
* Label for 'ThresholdType' field changed to 'TrendType' and field only visible when TriggerType selected is 'Trend'.
* Fixed errors on showing 'IsCustomExpression' field and related for alert component.
* Fixed errors on showing empty threshold fields when value was 0 for alert component.
* Fixed errors on tooltips for alert component.
* Fixed error on logging endpoint when the directory with the logfile does not exist.
* Fixed error on product modify. When 'unselecting' one measurement the related tags were not 'unselected' from taglist fields.
* Fixed error on renaming an alert. The related kapacitor task with the new name was created, but the related kapacitor task with the old name was not deleted.
* Fixed error when creating an alert with lambda expression.

### breaking changes
* 'deviceid_tag' column changed to 'producttag' on 'alert_id_cfg' table.
* 'groupid' column changed to 'alertgroup' on 'alert_id_cfg' table.
* Execute the following sql to update your table:
    * UPDATE alert_id_cfg SET producttag = deviceid_tag, alertgroup = groupid
* 'trigertype' column changed to 'triggertype' on 'alert_id_cfg' and 'template_cfg' tables.
* 'thresholdtype' column changed to 'trendtype' on 'alert_id_cfg' and 'template_cfg' tables.
* Execute the following sqls to update your tables:
    * UPDATE alert_id_cfg SET triggertype = trigertype, trendtype = thresholdtype
    * UPDATE template_cfg SET triggertype = trigertype, trendtype = thresholdtype


# v 0.3.0  (03/08/2018)
### New features.
* Added udf config info.
* The data to check on alerts can be a measurement field or a lambda expression.
* Field ExtraData added on alerts to use on MovingAverage and Percentile functions.
* Field Active added on alerts to activate or deactivate the related kapacitor task.
* Field InfluxFilter added on alerts to include a lambda expression to filter data from Influx.
* Fields on Alert Definition screen reordered and/or restyled.
* New component Kapacitor Tasks added on Runtime.
* Columns on Alert Events formatted.
* TimeLogs added on resinjector UDF.
* Port changed from 8090 to 6090.

### fixes
* Function MOVAVG changed to MOVINGAVERAGE on Templates and Alerts definition.
* Literals for Alerting Endpoints modified.
* Modal window modified to show html info.
* ng-table modified to show column tooltips correctly.
* Literal 'ThresholdType' changed to 'TrendType'.

### breaking changes
* Changes on Product component to improve configuration of products.
* 'commontags' column changed to 'products' on 'product_group_cfg' table.

# v 0.2.0  (31/07/2018)
### New features.
* Added resinjector sample file and options to build and package the UDF module.
* Added variable period to reload DB data on resinjector sample file.
* Added new resinjector deb package.
* Added rpm packaging files for resinjector UDF module.
* OutHttp modified to fit with new requirements.
* Change DevStats.ProductID field from single_select to input field.
* First version of logging and slack endpoints.
* resinjector UDF: new property method productID added and changes on logs.
* New component AlertEvent added.

### fixes
* Fixes for #23, #24, #28, #32
* Fix alertcfg export, added devicestat on import.
* Fix devicestat export, ensure string on device ID.
* Changes for importing devicestats and templates.
* OrderBy added to Device Stats.
* Fix resinjector service start CHDIR spawn error.
* Fix config file names on some docs and scripts.
* Little changes on texts for Product Group.
* Getting resistor_port for kapacitor vars.
* New pipe for left padding of AlertID with leading zeroes to get 001, 002, etc.
* Deploy templates and tasks on kapacitor when importing from file.
* Alert extra text field added on alert screen.

### breaking changes
* WeeKDays changed to WeekDays on range_time_cfg.
* Changes for TemplateCfg and DeviceStatCfg.
* Columns 'order' and 'exception' from table 'device_stat_cfg' changed to 'orderid' and 'exceptionid'.
* New field TrendSign and other changes for naming convention on templates.

# v 0.1.0  (never released!!)
### New features.
* Added main code skel
* Added new HTTP wrapper, and UI components
* Updated WebUI to Angular4
* Added Resinjector component
* WebUI skel and improvements.

### fixes

### breaking changes
