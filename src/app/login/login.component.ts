
import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { LoginService } from './login.service'

@Component({
  selector: 'login',
  providers: [LoginService],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})

export class LoginComponent {
  constructor(public router: Router, public loginService: LoginService) {
  }
  ifErrors: any;

  login(event, username, password) {
    event.preventDefault();
    let body = JSON.stringify({ username, password });
    this.loginService.login(body)
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
