# v 0.4.0  (unreleased )
### New features.
* Changes for defining alerts depending on product and some errors fixed.
    * 'InfluxMeasurement' field filled in with data depending on selected product.
    * 'BaselineID' field filled in with data depending on selected product.
    * 'GroupID' field filled in with data depending on selected product and label changed to 'AlertGroup'.
    * 'InfluxDB' field filled in with data depending on selected measurement.
    * 'DeviceIDTag' field filled in with data depending on selected product and label changed to 'ProductTag'.
* 'Alerts' label changed to 'Alert Definitions'.
* Changes for defining device stats (alert exceptions) based on product.
* Changes for Alerting Endpoints refactoring. Now the configuration for httppost and logging is done with several form fields. The configuration for slack still is done with a form field in JSON format.
* The size of TplData textarea has been increased.
* The Import Data button for Influx DB Servers is disabled on create mode.

### fixes
* Label for 'ThresholdType' field changed to 'TrendType' and field only visible when TriggerType selected is 'Trend'.
* Fixed errors on showing 'IsCustomExpression' field and related for alert component.
* Fixed errors on showing empty threshold fields when value was 0 for alert component.
* Fixed errors on some tooltips for alert component.
* Fixed error on logging endpoint when the directory with the logfile does not exist.
* Fixed error on product modify. When 'unselecting' one measurement the related tags were not 'unselected' from taglist fields.
* Fixed error on renaming an alert. The related kapacitor task with the new name was created, but the related kapacitor task with the old name was not deleted.
* Fixed error when creating an alert with lambda expression.

### breaking changes
* 'deviceid_tag' column changed to 'producttag' on 'alert_id_cfg' table.
* 'groupid' column changed to 'alertgroup' on 'alert_id_cfg' table.

# v 0.3.0  (unreleased )
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

# v 0.2.0  (unreleased )
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
