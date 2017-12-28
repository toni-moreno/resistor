export const ProductComponentConfig: any =
  {
    'name' : 'Product',
    'table-columns' : [
      { 'title': 'ID', 'name': 'ID' },
      { 'title': 'ID Tag Name', 'name': 'IDTagName' },
      { 'title': 'CommonTags', 'name': 'CommonTags' },
      { 'title': 'Base Lines', 'name': 'BaseLines' },
    ]
  };
  export const TableRole : string = 'fulledit';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'export', 'type':'icon', 'icon' : 'glyphicon glyphicon-download-alt text-info', 'tooltip': 'Export item'},
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'edit', 'type':'icon', 'icon' : 'glyphicon glyphicon-edit text-warning', 'tooltip': 'Edit item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'},
  ]
