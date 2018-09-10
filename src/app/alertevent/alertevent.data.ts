export const AlertEventComponentConfig: any =
  {
    'name' : 'Alert Events',
    'table-columns' : [
      {'title':'ID','name':'ID'},
      {'title':'AlertID','name':'AlertID'},
      {'title':'ProductID','name':'ProductID'},
      {'title':'ProductTagValue','name':'ProductTagValue'},
      {'title':'Tags','name':'Tags'},
      {'title':'Field','name':'Field'},
      {'title':'Value','name':'Value','transform':'decimal'},
      {'title':'Time','name':'Time','transform':'datetime'},
      {'title':'Duration','name':'Duration','transform':'ns2s'},
      {'title':'Level','name':'Level','transform':'color'},
      //{'title':'PreviousLevel','name':'PreviousLevel'}
    ],
    'slug' : 'alerteventhist'
  };
  export const TableRole : string = 'viewdelete';
  export const FilterColumn : string = 'ProductID';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'},
  ]
