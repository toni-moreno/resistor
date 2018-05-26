import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { HttpService } from '../core/http.service';

declare var _:any;

@Injectable()
export class LoginService {

    constructor(private http: HttpService) {
    }

    login(data) {
        return  this.http.post('/login', data, null, true)
        .map( (responseData) => responseData.json());
    }
}
