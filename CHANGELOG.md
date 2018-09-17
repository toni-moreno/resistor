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
