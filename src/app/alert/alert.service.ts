import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class AlertService {

    constructor(private http: HttpService,) {
    }

    jsonParser(key,value) {
        if ( key == 'Active' || key == 'IsCustomExpression' || key == 'Rate' ) {
            return ( value === "true" || value === true);
        }
        if ( (key == 'ExtraData' || key == 'AlertNotify')  && (value == null || value.length == 0) ) {
            return 0;
        }
        if ( key == 'AlertNotify' ) {
            return parseInt(value);
        }
        if ( key == 'Endpoint' ) {
            return String(value).split(',');
        }
        return value;
    }

    addAlertItem(dev) {
        return this.http.post('/api/cfg/alertid',JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());

    }

    deployAlertItem(dev) {
        return this.http.post('/api/cfg/alertid/deploy',JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());

    }

    editAlertItem(dev, id) {
        return this.http.put('/api/cfg/alertid/'+id,JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());
    }


    getAlertItem(filter_s: string) {
        return this.http.get('/api/cfg/alertid')
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getAlertItemById(id : string) {
        return this.http.get('/api/cfg/alertid/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    getAlertItemByProductId(productid : string) {
        return this.http.get('/api/cfg/alertid/byproductid/'+productid)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteAlertItem(id : string){
      return this.http.get('/api/cfg/alertid/checkondel/'+id)
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

    deleteAlertItem(id : string) {
        // return an observable
        return this.http.delete('/api/cfg/alertid/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };
}
