import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../../core/http.service';

declare var _:any;

@Injectable()
export class ExportServiceCfg {

    constructor(public http: HttpService) {
        console.log('Task Service created.', http);
    }

    exportFastRecursive(type : string, id : string) {
        // return an observable
        return this.http.get('/api/cfg/export/'+type+'/'+id)
        .map((res) => {
        //return new Blob([res.arrayBuffer()],{type: "application/octet-stream" })
        return [new Blob([res['_body']],{type: "application/json"}),res.json()];
        })
    }

    bulkExport(values) {
      return this.http.post('/api/cfg/bulkexport',values)
      .map((res) => {
          console.log(res);
          return [new Blob([res['_body']],{type: "application/json"}),res.json()];
      })
    }

    exportRecursive(type : string, id : string, values) {
        console.log(values);
        // return an observable
        return this.http.post('/api/cfg/export/'+type+'/'+id, values)
        .map((res) => {
            return [new Blob([res['_body']],{type: "application/json"}),res.json()];
        })
    }
}


/*this.exportServiceCfg.exportRecursive(item.exportType, item.row.ID).subscribe(
data => {
  console.log(data);
  saveAs(data[0], data[1].Info.FileName || item.row.ID + ".json");
},
err => console.error(err),
() => console.log("DONE"),
)
console.log(item);
*/
