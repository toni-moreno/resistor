export const ProductComponentConfig: any =
  {
    'name' : 'Product',
    'table-columns' : [
      { 'title': 'ID', 'name': 'ID' },
      { 'title': 'Base Lines', 'name': 'BaseLines' },
      { 'title': 'AlertGroups', 'name': 'AlertGroups' },
      { 'title': 'FieldEvalPeriods', 'name': 'FieldResolutions' },
      { 'title': 'Measurements', 'name': 'Measurements' },
      { 'title': 'Product Tag', 'name': 'ProductTag' },
      { 'title': 'CommonTags', 'name': 'CommonTags' },
      { 'title': 'ExtraTags', 'name': 'ExtraTags' },
      { 'title': 'Imported', 'name': 'Imported','transform':'datetime' },
    ],
    'slug' : 'productcfg'
  };
  export const TableRole : string = 'fulledit';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'export', 'type':'icon', 'icon' : 'glyphicon glyphicon-download-alt text-info', 'tooltip': 'Export item'},
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'edit', 'type':'icon', 'icon' : 'glyphicon glyphicon-edit text-warning', 'tooltip': 'Edit item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'},
  ]
