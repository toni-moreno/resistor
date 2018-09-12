export const KapacitorTasksComponentRt: any =
  {
    'name' : 'Kapacitor Tasks',
    'table-columns' : [
      {'title':'ServerID','name':'ServerID',filtering: {filterString: '', placeholder: 'Filter by ServerID'} },
      {'title':'URL','name':'URL',filtering: {filterString: '', placeholder: 'Filter by URL'} },
      {'title':'TaskID','name':'ID',filtering: {filterString: '', placeholder: 'Filter by TaskID'} },
      {'title':'Type','name':'Type',filtering: {filterString: '', placeholder: 'Filter by Type'} },
      {'title':'DBRPs','name':'DBRPs',filtering: {filterString: '', placeholder: 'Filter by DBRPs'} },
      {'title':'Status','name':'Status',filtering: {filterString: '', placeholder: 'Filter by Status'} },
      {'title':'Executing','name':'Executing',filtering: {filterString: '', placeholder: 'Filter by Executing'} },
      {'title':'Error','name':'Error','transform':'imgwtooltip',filtering: {filterString: '', placeholder: 'Filter by Error'} },
      {'title':'NumErrors','name':'NumErrors',filtering: {filterString: '', placeholder: 'Filter by NumErrors'} },
      {'title':'Created','name':'Created','transform':'datetime' },
      {'title':'Modified','name':'Modified','transform':'datetime' },
      {'title':'LastEnabled','name':'LastEnabled','transform':'datetime' },
      /*  
      {'title':'TICKscript','name':'TICKscript' },
      {'title':'Vars','name':'Vars' },
      {'title':'Dot','name':'Dot' },
{'title':'ExecutionStats','name':'ExecutionStats' },
*/
],

    'slug' : 'kapacitortasksrt'
  };
  export const TableRole : string = 'viewonly';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
  ]
