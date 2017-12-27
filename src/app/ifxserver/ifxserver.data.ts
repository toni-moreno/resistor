export const IfxServerComponentConfig: any =
  {
    'name' : 'Influx DB Servers',
    'table-columns' : [
      { 'title': 'ID', 'name': 'ID' },
      { 'title': 'Connection URL', 'name': 'URL' },
      { 'title': 'Admin User', 'name': 'AdminUser' },
      { 'title': 'AdminPasswd', 'name': 'AdminPasswd' },
      { 'title': 'Description', 'name': 'Description' }
    ]
  };

  export const TableRole : string = 'fulledit';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'export', 'type':'icon', 'icon' : 'glyphicon glyphicon-download-alt text-info', 'tooltip': 'Export item'},
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'edit', 'type':'icon', 'icon' : 'glyphicon glyphicon-edit text-warning', 'tooltip': 'Edit item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'},
    {'name':'importcatalog', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-flash text-info', 'tooltip': 'Import catalog'}
  ]
