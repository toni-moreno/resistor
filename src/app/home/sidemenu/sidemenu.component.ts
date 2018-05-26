import { Component, ViewChild,ViewContainerRef, Input, Output, EventEmitter } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { WindowRef } from '../../common/windowref';


@Component({
  selector: 'res-sidemenu',
  templateUrl: './sidemenu.component.html'
})

export class SideMenuComponent {

  item_type: string;

  @Input() menuItems : any;
  @Input() mode : any;

  @Output() public clickMenu:EventEmitter<any> = new EventEmitter<any>();
  @Output() public clickButton:EventEmitter<any> = new EventEmitter<any>();
  @Output() public showModal:EventEmitter<any> = new EventEmitter<any>();
  @Output() public link:EventEmitter<any> = new EventEmitter<any>();

  nativeWindow: any

  constructor(){ }

  linkClicked() {
    this.link.emit()
  }

  expandMenu(i : any) : boolean{
    this.menuItems[i].expanded = !this.menuItems[i].expanded;
    return this.menuItems[i].expanded;
  }

  shownAboutModal() {
    this.showModal.emit();
  }

  clickedMenu(menuItem : any) : void {
    this.item_type = menuItem.selector;
    this.clickMenu.emit(menuItem);
  }

  clickedButton(menuItem : any) : void {
    console.log("CLICKED",menuItem);
    this.clickButton.emit(menuItem);
  }
}
