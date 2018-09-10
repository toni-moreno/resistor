import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class AlertEventService {

    constructor(private http: HttpService,) {
    }

    jsonParser(key,value) {
        return value;
    }

    editAlertEventItem(dev, id) {
        return this.http.put('/api/rt/alertevent/'+id,JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());
    }

    getAlertEventWithParams(action? : any) {
        let params : string = "";
        //Add empty parameter to avoid error on backend
        params = params + "&empty=";
        if (action != null) {
            if (action.filterString.length > 0) {
                params = params + "&filterColumn=" + action.filterColumn;
                params = params + "&filterString=" + action.filterString;
            }
            params = params + "&page=" + action.page;
            params = params + "&itemsPerPage=" + action.itemsPerPage;
            params = params + "&maxSize=" + action.maxSize;
            params = params + "&sortColumn=" + action.sortColumn;
            params = params + "&sortDir=" + action.sortDir;
        }
        return this.http.get('/api/rt/alertevent/withparams/'+params)
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getAlertEventItem(filter_s: string) {
        return this.http.get('/api/rt/alertevent')
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getAlertEventItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/rt/alertevent/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteAlertEventItem(id : string){
        return this.http.get('/api/rt/alertevent/checkondel/'+id)
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
  
    deleteAlertEventItem(id : string) {
        // return an observable
        return this.http.delete('/api/rt/alertevent/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };

}
