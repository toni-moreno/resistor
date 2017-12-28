export const AlertComponentConfig: any =
  {
    'name' : 'Alerts',
    'table-columns' : [
      { 'title': 'ID', 'name': 'ID' },
      { 'title': 'BaselineID', 'name': 'BaselineID' },
      { 'title': 'ProductID', 'name': 'ProductID' },
      { 'title': 'GroupID', 'name': 'GroupID' },
      { 'title': 'NumAlertID', 'name': 'NumAlertID' },
      { 'title': 'InfluxDB', 'name': 'InfluxDB' },
      { 'title': 'InfluxRP', 'name': 'InfluxRP' },
      { 'title': 'InfluxMeasurement', 'name': 'InfluxMeasurement' },
      { 'title': 'TagDescription', 'name': 'TagDescription' },
      { 'title': 'InfluxFilter', 'name': 'InfluxFilter' },
      { 'title': 'TrigerType', 'name': 'TrigerType' },
      { 'title': 'IntervalCheck', 'name': 'IntervalCheck' },
      { 'title': 'OperationID', 'name': 'OperationID' },
      { 'title': 'Field', 'name': 'Field' },
      { 'title': 'StatFunc', 'name': 'StatFunc' },
      { 'title': 'CritDirection', 'name': 'CritDirection' },
      { 'title': 'Shift', 'name': 'Shift' },
      { 'title': 'ThresholdType', 'name': 'ThresholdType' },
      { 'title': 'ThCritDef', 'name': 'ThCritDef' },
      { 'title': 'ThCritEx1', 'name': 'ThCritDef' },
      { 'title': 'ThCritEx2', 'name': 'ThCritDef' },
      { 'title': 'ThCritRangeTimeID', 'name': 'ThCritDef' },
      { 'title': 'ThWarnDef', 'name': 'ThCritDef' },
      { 'title': 'ThWarnEx1', 'name': 'ThCritDef' },
      { 'title': 'ThWarnEx2', 'name': 'ThCritDef' },
      { 'title': 'ThWarnRangeTimeID', 'name': 'ThCritDef' },
      { 'title': 'ThWarnEx2', 'name': 'ThCritDef' },
      { 'title': 'ThInfoDef', 'name': 'ThCritDef' },
      { 'title': 'ThInfoEx1', 'name': 'ThCritDef' },
      { 'title': 'ThInfoEx2', 'name': 'ThCritDef' },
      { 'title': 'ThInfoRangeTimeID', 'name': 'ThCritDef' },
      { 'title': 'GrafanaServer', 'name': 'ThCritDef' },
      { 'title': 'GrafanaDashLabel', 'name': 'ThCritDef' },
      { 'title': 'GrafanaDashPanelID', 'name': 'ThCritDef' },
      { 'title': 'DeviceIDTag', 'name': 'ThCritDef' },
      { 'title': 'DeviceIDLabel', 'name': 'ThCritDef' },
      { 'title': 'ExtraTag', 'name': 'ThCritDef' },
      { 'title': 'ExtraLabel', 'name': 'ThCritDef' },
      { 'title': 'KapacitorID', 'name': 'ThCritDef' },
      { 'title': 'OutHTTP', 'name': 'ThCritDef' }
    ]
  };

  export const TableRole : string = 'fulledit';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'export', 'type':'icon', 'icon' : 'glyphicon glyphicon-download-alt text-info', 'tooltip': 'Export item'},
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'edit', 'type':'icon', 'icon' : 'glyphicon glyphicon-edit text-warning', 'tooltip': 'Edit item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'},
  ]
