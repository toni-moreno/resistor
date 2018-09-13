export const AlertEventComponentConfig: any =
  {
    'name' : 'Alert Events',
    'table-columns' : [
      {'title':'ID','name':'ID',filtering: {filterString: '', placeholder: 'Filter by ID'},sort:'desc'},
      {'title':'AlertID','name':'AlertID',filtering: {filterString: '', placeholder: 'Filter by AlertID'}},
      {'title':'ProductID','name':'ProductID',filtering: {filterString: '', placeholder: 'Filter by ProductID'}},
      {'title':'ProductTag:Value','name':'ProductTagValue',filtering: {filterString: '', placeholder: 'Filter by ProductTagValue'}},
      {'title':'Tags','name':'Tags',filtering: {filterString: '', placeholder: 'Filter by Tags'}},
      {'title':'Field','name':'Field',filtering: {filterString: '', placeholder: 'Filter by Field'}},
      {'title':'Value','name':'Value','transform':'decimal'},
      {'title':'Time','name':'Time','transform':'datetime'},
      {'title':'Duration','name':'Duration','transform':'ns2s'},
      {'title':'Level','name':'Level','transform':'color',filtering: {filterString: '', placeholder: 'Filter by Level'}},
      //{'title':'PreviousLevel','name':'PreviousLevel'}
    ],
    'slug' : 'alerteventhist'
  };
  export const TableRole : string = 'viewdelete';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'},
  ]
