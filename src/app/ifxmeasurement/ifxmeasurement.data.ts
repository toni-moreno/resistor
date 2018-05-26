export const IfxMeasurementComponentConfig: any =
  {
    'name' : 'Influx Measurements',
    'table-columns' : [
      { 'title': 'Measurement Name', 'name': 'Name' },
      { 'title': 'Tags', 'name': 'Tags' },
      { 'title': 'Fields', 'name': 'Fields' },
    ],
    'slug' : 'ifxmeasurementcfg'
  };

  export const TableRole : string = 'viewdelete';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'},
  ]
