export const DeviceStatComponentConfig: any =
  {
    'name' : 'Device Stats',
    'table-columns' : [
      { 'title': 'ID', 'name': 'ID' },
      { 'title': 'DeviceID', 'name': 'DeviceID' },
      { 'title': 'AlertID', 'name': 'AlertID' },
      { 'title': 'OrderID', 'name': 'OrderID' },
      { 'title': 'ProductID', 'name': 'ProductID' },
      { 'title': 'ExceptionID', 'name': 'ExceptionID' },
      { 'title': 'Active', 'name': 'Active' },
      { 'title': 'BaseLine', 'name': 'BaseLine' },
      { 'title': 'FilterTagKey', 'name': 'FilterTagKey' },
      { 'title': 'FilterTagValue', 'name': 'FilterTagValue' },
      { 'title': 'Description', 'name': 'Description' },
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
