import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class TemplateService {

    constructor(private http: HttpService,) {
    }

    jsonParser(key,value) {
        return value
    }

    addTemplateItem(dev) {
        return this.http.post('/api/cfg/template',JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());

    }

    deployTemplateItem(dev) {
        return this.http.post('/api/cfg/template/deploy',JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());

    }

    editTemplateItem(dev, id) {
        return this.http.put('/api/cfg/template/'+id,JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());
    }


    getTemplateItem(filter_s: string) {
        return this.http.get('/api/cfg/template')
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getTemplateItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/cfg/template/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteTemplateItem(id : string){
      return this.http.get('/api/cfg/template/checkondel/'+id)
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

    deleteTemplateItem(id : string) {
        // return an observable
        return this.http.delete('/api/cfg/template/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };
}
