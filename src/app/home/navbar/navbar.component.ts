import { Component, ViewChild,ViewContainerRef, Input, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'res-navbar',
  templateUrl: './navbar.component.html'
})

export class NavbarComponent {

  @Input() version : any;
  @Output() public toogleMenu:EventEmitter<any> = new EventEmitter<any>();
  @Output() public logout:EventEmitter<any> = new EventEmitter<any>();
  @Output() public showModal:EventEmitter<any> = new EventEmitter<any>();
  @Output() public link:EventEmitter<any> = new EventEmitter<any>();

  constructor() { }

  toggleClicked() {
    this.toogleMenu.emit();
  }

  logOutClicked() {
    this.logout.emit();
  }

  shownAboutModal() {
    this.showModal.emit();
  }

  linkClicked() {
    this.link.emit();
  }

}
