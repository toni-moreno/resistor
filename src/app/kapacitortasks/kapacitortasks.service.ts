import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class KapacitorTasksService {

    constructor(private http: HttpService,) {
    }

    jsonParser(key,value) {
        return value;
    }

    getKapacitorTasksItem(filter_s: string) {
        return this.http.get('/api/rt/kapacitor/tasks')
        .map( (responseData) => {
            return responseData.json();
        })
    }

    getKapacitorTasksItemById(id : string) {
        // return an observable
        console.log("ID: ",id);
        return this.http.get('/api/rt/kapacitor/tasks/'+id)
        .map( (responseData) =>
            responseData.json()
    )};

}
