export const KapacitorComponentConfig: any =
  {
    'name' : 'Kapacitor Backends',
    'table-columns' : [
      { 'title': 'ID', 'name': 'ID' },
      { 'title': 'URL', 'name': 'URL' },
      { 'title': 'Imported', 'name': 'Imported','transform':'datetime' },
    ],
    'slug' : 'kapacitorcfg'
  };
  export const TableRole : string = 'fulledit';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'export', 'type':'icon', 'icon' : 'glyphicon glyphicon-download-alt text-info', 'tooltip': 'Export item'},
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'test-connection', 'type':'icon', 'icon' : 'glyphicon glyphicon-flash text-info', 'tooltip': 'Test connection'},
    {'name':'edit', 'type':'icon', 'icon' : 'glyphicon glyphicon-edit text-warning', 'tooltip': 'Edit item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'}
  ]
