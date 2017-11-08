import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class OutHTTPService {

    constructor(private http: HttpService,) {
    }

    addOutHTTPItem(dev) {
        return this.http.post('/api/cfg/outhttp',JSON.stringify(dev,function (key,value) {
             if ( key == 'Headers') {
                  return value.split(',');
                }

                return value;
        }))
        .map( (responseData) => responseData.json());

    }

    editOutHTTPItem(dev, id) {
        return this.http.put('/api/cfg/outhttp/'+id,JSON.stringify(dev,function (key,value) {
            if ( key == 'Headers') {
                 return value.split(',');
               }
            return value;
        }))
        .map( (responseData) => responseData.json());
    }


    getOutHTTPItem(filter_s: string) {
        return this.http.get('/api/cfg/outhttp')
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getOutHTTPItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/cfg/outhttp/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteOutHTTPItem(id : string){
      return this.http.get('/api/cfg/outhttp/checkondel/'+id)
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

    deleteOutHTTPItem(id : string) {
        // return an observable
        return this.http.delete('/api/cfg/outhttp/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };
}
