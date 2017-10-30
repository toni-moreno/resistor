import { Injectable } from '@angular/core';
import { HttpAPI } from '../common/httpAPI'
import { Observable } from 'rxjs/Observable';

import { TableData } from './sample-component.data'

declare var _:any;

@Injectable()
export class SampleComponentService {

    constructor(public httpAPI: HttpAPI) {
    }

    addSampleItem(dev) {
        return this.httpAPI.post('/api/cfg/influxservers',JSON.stringify(dev,function (key,value) {
                if ( key == 'Port'  ||
                key == 'Timeout' ) {
                  return parseInt(value);
                }
                return value;
        }))
        .map( (responseData) => responseData.json());

    }

    editSampleItem(dev, id) {
        return this.httpAPI.put('/api/cfg/influxservers/'+id,JSON.stringify(dev,function (key,value) {
            if ( key == 'Port'  ||
            key == 'Timeout' ) {
              return parseInt(value);
            }
            return value;

        }))
        .map( (responseData) => responseData.json());
    }


    getSampleItem(filter_s: string) {
        // return an observable
        let data : Array<any> = TableData;
        let test = Observable.of(data);

        return test;

        /*return this.httpAPI.get('/api/cfg/influxservers')
        .map( (responseData) => {
            return responseData.json();
        })
        */
    }

    getSampleItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.httpAPI.get('/api/cfg/influxservers/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteSampleItem(id : string){
      return this.httpAPI.get('/api/cfg/influxservers/checkondel/'+id)
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

    testSampleItem(influxserver) {
      // return an observable
      return this.httpAPI.post('/api/cfg/influxservers/ping/',JSON.stringify(influxserver,function (key,value) {
          if ( key == 'Port'  ||
          key == 'Timeout' ) {
            return parseInt(value);
          }
          return value;
      }))
      .map(
        (responseData) => responseData.json()
      );
    };

    deleteSampleItem(id : string) {
        // return an observable
        return this.httpAPI.delete('/api/cfg/influxservers/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };
}
