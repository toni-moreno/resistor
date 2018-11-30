export const OperationComponentConfig: any =
  {
    'name' : 'Operations',
    'table-columns' : [
      { 'title': 'ID', 'name': 'ID' },
      { 'title': 'URL', 'name': 'URL' },
      { 'title': 'Imported', 'name': 'Imported','transform':'datetime' },
    ],
    'slug' : 'operationcfg'
  };
  export const TableRole : string = 'fulledit';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'export', 'type':'icon', 'icon' : 'glyphicon glyphicon-download-alt text-info', 'tooltip': 'Export item'},
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'edit', 'type':'icon', 'icon' : 'glyphicon glyphicon-edit text-warning', 'tooltip': 'Edit item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'}
  ]
