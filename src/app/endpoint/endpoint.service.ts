import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class EndpointService {

    constructor(private http: HttpService,) {
    }

    jsonParser(key,value) {
        if ( key == 'Headers' || key == 'To' ) {
            if(typeof value === 'string') return value.split(',');
        }
        if ( key == 'Enabled' || key == 'InsecureSkipVerify' ) {
            return ( value === "true" || value === true);
        }
        return value;
    }

    addEndpointItem(dev) {
        return this.http.post('/api/cfg/endpoint',JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());

    }

    editEndpointItem(dev, id) {
        return this.http.put('/api/cfg/endpoint/'+id,JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());
    }


    getEndpointItem(filter_s: string) {
        return this.http.get('/api/cfg/endpoint')
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getEndpointItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/cfg/endpoint/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteEndpointItem(id : string){
      return this.http.get('/api/cfg/endpoint/checkondel/'+id)
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

    deleteEndpointItem(id : string) {
        // return an observable
        return this.http.delete('/api/cfg/endpoint/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };
}
