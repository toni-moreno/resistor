export const IfxDBComponentConfig: any =
  {
    'name' : 'Influx Databases',
    'table-columns' : [
      { 'title': 'Influx Server', 'name': 'IfxServer' },
      { 'title': 'DB Name', 'name': 'Name' },
      { 'title': 'Retention', 'name': 'Retention' },
      { 'title': 'Measurements', 'name': 'Measurements', 'transform':"parseMeasurements" },
    ]
  };

export const TableRole : string = 'viewdelete';
export const OverrideRoleActions : Array<Object> = [
  {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
  {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'},
]

export function TableCellParser (data: any, column: string) {
  if (column === "parseMeasurements") {
    if (data){
      var test: any = '<ul class="list-unstyled">';
      for (var i of data) {
          test +="<li>"
          test +="<span>"+i.Name+"</span>";
          test += "</li>";
      }
      test += "</ul>"
      return test
    }
  }
  return ""
}
