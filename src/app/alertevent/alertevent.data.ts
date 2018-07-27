export const AlertEventComponentConfig: any =
  {
    'name' : 'Alert Events',
    'table-columns' : [
      {'title':'UID','name':'UID'},
      {'title':'ID','name':'ID'},
      {'title':'Message','name':'Message'},
      {'title':'Time','name':'Time'},
      {'title':'Duration','name':'Duration'},
      {'title':'Level','name':'Level'},
      //{'title':'PreviousLevel','name':'PreviousLevel'}
    ],
    'slug' : 'alerteventcfg'
  };
  export const TableRole : string = 'fulledit';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
  ]
