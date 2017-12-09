import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class IfxMeasurementService {

    constructor(private http: HttpService,) {
    }

    addIfxMeasurementItem(dev) {
        return this.http.post('/api/cfg/ifxmeasurement',JSON.stringify(dev,function (key,value) {
             if ( key == 'Tags'||
                  key == 'Fields' ) {
                  return value.split(',');
                }

                return value;
        }))
        .map( (responseData) => responseData.json());

    }

    editIfxMeasurementItem(dev, id) {
        return this.http.put('/api/cfg/ifxmeasurement/'+id,JSON.stringify(dev,function (key,value) {
            if ( key == 'Tags'||
                 key == 'Fields' ) {
            return value.split(',');
          }
            return value;
        }))
        .map( (responseData) => responseData.json());
    }


    getIfxMeasurementItem(filter_s: string) {
        return this.http.get('/api/cfg/ifxmeasurement')
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getIfxMeasurementItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/cfg/ifxmeasurement/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteIfxMeasurementItem(id : string){
      return this.http.get('/api/cfg/ifxmeasurement/checkondel/'+id)
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

    deleteIfxMeasurementItem(id : string) {
        // return an observable
        return this.http.delete('/api/cfg/ifxmeasurement/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };
}
