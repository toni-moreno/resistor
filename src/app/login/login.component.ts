
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { HttpService} from '../core/http.service'

@Component({
  selector: 'login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})

export class LoginComponent {
  constructor(public router: Router, public http: HttpService) {
  }
  ifErrors: any;

  login(event, username, password) {
    event.preventDefault();
    let body = JSON.stringify({ username, password });
    this.http.post('/login', body, null, true)
      .subscribe(
        response => {
          this.router.navigate(['home']);
        },
        error => {
          this.ifErrors = error['_body'];
          console.log(error.text());
        }
      );
  }
}
