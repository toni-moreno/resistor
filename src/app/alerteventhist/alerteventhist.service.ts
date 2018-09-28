import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class AlertEventHistService {

    constructor(private http: HttpService,) {
    }

    jsonParser(key,value) {
        return value;
    }

    editAlertEventHistItem(dev, id) {
        return this.http.put('/api/rt/alerteventhist/'+id,JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());
    }

    getAlertEventHistWithParams(action? : any) {
        let params : string = "";
        //Add empty parameter to avoid error on backend
        params = params + "&empty=";
        if (action != null) {
            params = params + "&page=" + action.page;
            params = params + "&itemsPerPage=" + action.itemsPerPage;
            params = params + "&maxSize=" + action.maxSize;
            params = params + "&sortColumn=" + action.sortColumn;
            params = params + "&sortDir=" + action.sortDir;
        }
        return this.http.get('/api/rt/alerteventhist/withparams/'+params)
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getAlertEventHistItem(filter_s: string) {
        return this.http.get('/api/rt/alerteventhist')
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getAlertEventHistItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/rt/alerteventhist/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteAlertEventHistItem(id : string){
        return this.http.get('/api/rt/alerteventhist/checkondel/'+id)
        .map( (responseData) =>
         responseData.json()
        ).map((deleteobject) => {
            console.log("MAP SERVICE",deleteobject);
            let result : any = {'ID' : id};
            _.forEach(deleteobject,function(value,key){
                result[value.TypeDesc] = [];
            });
            _.forEach(deleteobject,function(value,key){
                result[value.TypeDesc].Description=value.Action;
                result[value.TypeDesc].push(value.ObID);
            });
            return result;
        });
      };
  
    deleteAlertEventHistItem(id : string) {
        // return an observable
        return this.http.delete('/api/rt/alerteventhist/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };

}
