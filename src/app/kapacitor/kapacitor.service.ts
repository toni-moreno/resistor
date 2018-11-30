import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class KapacitorService {

    constructor(private http: HttpService,) {
    }

    jsonParser(key,value) {
        return value;
    }

    addKapacitorItem(dev) {
        return this.http.post('/api/cfg/kapacitor',JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());

    }

    editKapacitorItem(dev, id) {
        return this.http.put('/api/cfg/kapacitor/'+id,JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());
    }


    getKapacitorItem(filter_s: string) {
        return this.http.get('/api/cfg/kapacitor')
        .map( (responseData) => {
            return responseData.json();
        })
        .map((kapacitors) => {
            console.log("MAP SERVICE",kapacitors);
            let result = [];
            if (kapacitors) {
                _.forEach(kapacitors,function(value,key){
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

    getKapacitorItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/cfg/kapacitor/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteKapacitorItem(id : string){
      return this.http.get('/api/cfg/kapacitor/checkondel/'+id)
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

    testKapacitorItem(instance) {
      // return an observable
      return this.http.post('/api/cfg/kapacitor/ping/',JSON.stringify(instance,function (key,value) {
          return value;
      }))
      .map(
        (responseData) => responseData.json()
      );
    };

    deleteKapacitorItem(id : string) {
        // return an observable
        return this.http.delete('/api/cfg/kapacitor/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };
}
