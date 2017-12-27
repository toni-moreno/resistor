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

export function tableCellParser (data: any, column: string) {
  if (column === "parseMeasurements") {
    console.log("DATA",data);
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
