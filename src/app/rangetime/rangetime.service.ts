import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class RangeTimeService {

    constructor(private http: HttpService,) {
    }

    jsonParser(key,value) {
        return value;
    }

    addRangeTimeItem(dev) {
        return this.http.post('/api/cfg/rangetimes',JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());

    }

    editRangeTimeItem(dev, id) {
        return this.http.put('/api/cfg/rangetimes/'+id,JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());
    }


    getRangeTimeItem(filter_s: string) {
        return this.http.get('/api/cfg/rangetimes')
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getRangeTimeItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/cfg/rangetimes/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteRangeTimeItem(id : string){
      return this.http.get('/api/cfg/rangetimes/checkondel/'+id)
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

    deleteRangeTimeItem(id : string) {
        // return an observable
        return this.http.delete('/api/cfg/rangetimes/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };
}
