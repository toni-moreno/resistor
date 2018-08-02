export const AlertEventComponentConfig: any =
  {
    'name' : 'Alert Events',
    'table-columns' : [
      {'title':'UID','name':'UID'},
      {'title':'ID','name':'ID'},
      {'title':'Message','name':'Message'},
      {'title':'Time','name':'Time'},
      {'title':'Duration','name':'Duration','transform':'ns2s'},
      {'title':'Level','name':'Level','transform':'color'},
      //{'title':'PreviousLevel','name':'PreviousLevel'}
    ],
    'slug' : 'alerteventcfg'
  };
  export const TableRole : string = 'viewonly';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
  ]
