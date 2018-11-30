import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class OperationService {

    constructor(private http: HttpService,) {
    }

    jsonParser(key,value) {
        return value;
    }

    addOperationItem(dev) {
        return this.http.post('/api/cfg/operation',JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());

    }

    editOperationItem(dev, id) {
        return this.http.put('/api/cfg/operation/'+id,JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());
    }


    getOperationItem(filter_s: string) {
        return this.http.get('/api/cfg/operation')
        .map( (responseData) => {
            return responseData.json();
        })
        .map((operations) => {
            console.log("MAP SERVICE",operations);
            let result = [];
            if (operations) {
                _.forEach(operations,function(value,key){
                    console.log("FOREACH LOOP",value,value.ID);
                    if(filter_s && filter_s.length > 0 ) {
                        console.log("maching: "+value.ID+ "filter: "+filter_s);
                        var re = new RegExp(filter_s, 'gi');
                        if (value.ID.match(re)){
                            result.push(value);
                        }
                        console.log(value.ID.match(re));
                    } else {
                        result.push(value);
                    }
                });
            }
            return result;
        });
    }

    getOperationItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/cfg/operation/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteOperationItem(id : string){
      return this.http.get('/api/cfg/operation/checkondel/'+id)
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

    deleteOperationItem(id : string) {
        // return an observable
        return this.http.delete('/api/cfg/operation/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };
}
