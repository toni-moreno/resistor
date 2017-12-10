import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class ProductGroupService {

    constructor(private http: HttpService,) {
    }

    jsonParser(key,value) {
        if ( key == 'ProductGroups') {
            return String(value).split(',');
          }
        return value
    }

    addProductGroupItem(dev) {
        return this.http.post('/api/cfg/productgroup',JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());

    }

    editProductGroupItem(dev, id) {
         return this.http.put('/api/cfg/productgroup/'+id,JSON.stringify(dev,this.jsonParser))
        .map( (responseData) => responseData.json());
    }


    getProductGroupItem(filter_s: string) {
        return this.http.get('/api/cfg/productgroup')
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getProductGroupItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/cfg/productgroup/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

    checkOnDeleteProductGroupItem(id : string){
      return this.http.get('/api/cfg/productgroup/checkondel/'+id)
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

    deleteProductGroupItem(id : string) {
        // return an observable
        return this.http.delete('/api/cfg/productgroup/'+id)
        .map( (responseData) =>
         responseData.json()
        );
    };
}
