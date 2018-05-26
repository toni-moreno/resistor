export const ProductGroupComponentConfig: any =
  {
    'name' : 'Product Groups',
    'table-columns' : [
      { 'title': 'ID', 'name': 'ID' },
      { 'title': 'Products', 'name': 'Products' },
    ],
    'slug' : 'productgroupcfg'
  };
  export const TableRole : string = 'fulledit';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'export', 'type':'icon', 'icon' : 'glyphicon glyphicon-download-alt text-info', 'tooltip': 'Export item'},
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'edit', 'type':'icon', 'icon' : 'glyphicon glyphicon-edit text-warning', 'tooltip': 'Edit item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'},
  ]
