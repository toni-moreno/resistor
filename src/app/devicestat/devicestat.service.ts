import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class DeviceStatService {

    constructor(private http: HttpService,) {
    }

    addDeviceStatItem(dev) {
        return this.http.post('/api/cfg/devicestat',JSON.stringify(dev,function (key,value) {
            if ( key == 'Active')
                return ( value === "true" || value === true);
            return value;
        }))
        .map( (responseData) => responseData.json());

    }

    editDeviceStatItem(dev, id) {
        return this.http.put('/api/cfg/devicestat/'+id,JSON.stringify(dev,function (key,value) {
            if ( key == 'Active')
                return ( value === "true" || value === true);
            return value;        }))
        .map( (responseData) => responseData.json());
    }


    getDeviceStatItem(filter_s: string) {
        return this.http.get('/api/cfg/devicestat')
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getDeviceStatItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/cfg/devicestat/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteDeviceStatItem(id : string){
      return this.http.get('/api/cfg/devicestat/checkondel/'+id)
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

    testDeviceStatItem(instance) {
      // return an observable
      return this.http.post('/api/cfg/devicestat/ping/',JSON.stringify(instance,function (key,value) {
          return value;
      }))
      .map(
        (responseData) => responseData.json()
      );
    };

    deleteDeviceStatItem(id : string) {
        // return an observable
        return this.http.delete('/api/cfg/devicestat/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };
}
