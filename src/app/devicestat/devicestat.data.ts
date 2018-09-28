export const DeviceStatComponentConfig: any =
  {
    'name' : 'Device Stats',
    'table-columns' : [
      { 'title': 'ID', 'name': 'ID','tooltip':'Unique identifier for the Device Stat'},
      { 'title': 'ProductID', 'name': 'ProductID','tooltip':'ID of the product to associate this exception' },
      { 'title': 'BaseLine', 'name': 'BaseLine','tooltip':'Line used for filtering' },
      { 'title': 'AlertID', 'name': 'AlertID','tooltip':'AlertID with format line-product-alertgroup-nnn. Regular expressions accepted.' },
      { 'title': 'DeviceID', 'name': 'DeviceID','tooltip':'Id of the Device or * for generic rules.' },
      { 'title': 'OrderID', 'name': 'OrderID','tooltip':'OrderID for application of rules' },
      { 'title': 'ExceptionID', 'name': 'ExceptionID','tooltip':'ID of the exception to apply (-1: alerts NOT sent, 0: default values for the alerts, 1: Ex1 values for the alerts, 2: Ex2 values for the alerts)' },
      { 'title': 'Active', 'name': 'Active','tooltip':'Indicates if this exception must be considered or not' },
      { 'title': 'FilterTagKey', 'name': 'FilterTagKey','tooltip':'Name of the tag used for filtering' },
      { 'title': 'FilterTagValue', 'name': 'FilterTagValue','tooltip':'Value of the tag used for filtering. Regular expressions accepted.' },
      { 'title': 'Description', 'name': 'Description','tooltip':'Description of the Device Stat' },
    ],
    'slug' : 'devicestatcfg'
  };
  export const TableRole : string = 'fulledit';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'export', 'type':'icon', 'icon' : 'glyphicon glyphicon-download-alt text-info', 'tooltip': 'Export item'},
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'edit', 'type':'icon', 'icon' : 'glyphicon glyphicon-edit text-warning', 'tooltip': 'Edit item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'},
  ]
