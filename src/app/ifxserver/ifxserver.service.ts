import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class IfxServerService {

    constructor(private http: HttpService,) {
    }

    addIfxServerItem(dev) {
        return this.http.post('/api/cfg/ifxserver',JSON.stringify(dev,function (key,value) {
             if ( key == 'CommonTags') {
                  return value.split(',');
                }

                return value;
        }))
        .map( (responseData) => responseData.json());

    }

    editIfxServerItem(dev, id) {
        return this.http.put('/api/cfg/ifxserver/'+id,JSON.stringify(dev,function (key,value) {
            if ( key == 'CommonTags') {
                 return value.split(',');
               }
            return value;
        }))
        .map( (responseData) => responseData.json());
    }


    getIfxServerItem(filter_s: string) {
        return this.http.get('/api/cfg/ifxserver')
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getIfxServerItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/cfg/ifxserver/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteIfxServerItem(id : string){
      return this.http.get('/api/cfg/ifxserver/checkondel/'+id)
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

    deleteIfxServerItem(id : string) {
        // return an observable
        return this.http.delete('/api/cfg/ifxserver/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };
}
