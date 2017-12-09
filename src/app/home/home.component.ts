import { Component, ViewChild,ViewContainerRef } from '@angular/core';
import { NgSwitch, NgSwitchCase, NgSwitchDefault } from '@angular/common';
import { Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { BlockUIService } from '../common/blockui/blockui-service';
import { BlockUIComponent } from '../common/blockui/blockui-component';
import { ImportFileModal } from '../common/dataservice/import-file-modal';
import { HomeService } from './home.service';
import { AboutModal } from './about-modal'
import { WindowRef } from '../common/windowref';

//Menu Components  to load them dynamically
import { IfxServerComponent } from '../ifxserver/ifxserver.component';
import { IfxDBComponent } from '../ifxdb/ifxdb.component';
import { IfxMeasurementComponent } from '../ifxmeasurement/ifxmeasurement.component';
import { KapacitorComponent } from '../kapacitor/kapacitor.component';
import { RangeTimeComponent } from '../rangetime/rangetime.component';
import { ProductComponent } from '../product/product.component';
import { TemplateComponent } from '../template/template.component';
import { OutHTTPComponent } from '../outhttp/outhttp.component';
import { AlertComponent } from '../alert/alert.component';
import { DeviceStatComponent } from '../devicestat/devicestat.component'
import { NavbarComponent } from './navbar/navbar.component'
import { SideMenuComponent } from './sidemenu/sidemenu.component'


declare var _:any;

@Component({
  selector: 'home',
  templateUrl: './home.component.html',
  styleUrls: [ './home.component.css' ],
  providers: [BlockUIService, HomeService]
})

export class HomeComponent {

  @ViewChild('blocker', { read: ViewContainerRef }) container: ViewContainerRef;
  @ViewChild('importFileModal') public importFileModal : ImportFileModal;
  @ViewChild('aboutModal') public aboutModal : AboutModal;
  @ViewChild('RuntimeComponent') public rt : any;

  nativeWindow: any
  response: string;
  item_type: string;
  version: RInfo;
  menuItems : Array<any> = [
  {'groupName' : 'Runtime', 'icon': 'glyphicon glyphicon-play', 'expanded': true, 'items':
    [
      {'title': 'Agent status', 'selector' : 'runtime', 'component': null}
    ]
  },
  {'groupName' : 'Influx Catalog', 'icon': 'glyphicon glyphicon-play', 'expanded': true, 'items':
  [
    {'title': 'Influx DB Servers ', 'selector' : 'ifxserver-component', 'component': IfxServerComponent},
    {'title': 'Influx Databases', 'selector' : 'ifdb-component', 'component': IfxDBComponent},
    {'title': 'Influx Measurements', 'selector' : 'ifmeasurement-component', 'component': IfxMeasurementComponent}
  ]
},
  {'groupName' : 'Configuration', 'icon': 'glyphicon glyphicon-cog', 'expanded': true, 'items':
    [
      {'title': 'Kapacitor Backends', 'selector' : 'kapacitor-component', 'component': KapacitorComponent},
      {'title': 'RangeTime', 'selector' : 'rangetime-component', 'component': RangeTimeComponent},
      {'title': 'Product', 'selector' : 'product-component', 'component': ProductComponent},
      {'title': 'Template', 'selector' : 'template-component', 'component': TemplateComponent},
      {'title': 'OutHTTP', 'selector' : 'outhttp-component', 'component': OutHTTPComponent},
      {'title': 'Alert Definition', 'selector' : 'alert-component', 'component': AlertComponent},
      {'title': 'Device Stats', 'selector' : 'devicestat-component', 'component': DeviceStatComponent},
    ]
  }];


  componentList = KapacitorComponent;

  mode : boolean = false;
  userIn : boolean = false;

  elapsedReload: string = '';
  lastReload: Date;

  constructor(private winRef: WindowRef,public router: Router, private _blocker: BlockUIService, public homeService: HomeService) {
    this.nativeWindow = winRef.nativeWindow;
    this.getFooterInfo();
    this.item_type= "kapacitor-component";
  }

  link(url: string) {
    this.nativeWindow.open(url);
  }

  expandMenu(i : any) : boolean{
    this.menuItems[i].expanded = !this.menuItems[i].expanded;
    return this.menuItems[i].expanded;
  }

  logout() {
    this.homeService.userLogout()
    .subscribe(
    response => {
      this.router.navigate(['/sign-in']);
    },
    error => {
      alert(error.text());
      console.log(error.text());
    }
    );
  }
  changeModeMenu() {
    this.mode = !this.mode
  }

  clickMenu(menuItem : any) : void {
    this.componentList = menuItem.component;
  }

  showImportModal() {
    this.importFileModal.initImport();
  }

  showExportBulkModal() {
    //this.exportBulkFileModal.initExportModal(null, false);
  }

  showAboutModal() {
    this.aboutModal.showModal(this.version);
  }

  reloadConfig() {
    this._blocker.start(this.container, "Reloading Conf. Please wait...");
    if (this.rt) this.rt.updateRuntimeInfo(null,null,false);
    this.homeService.reloadConfig()
    .subscribe(
    response => {
      this.lastReload = new Date();
      this.elapsedReload = response;
      this._blocker.stop();
    },
    error => {
      this._blocker.stop();
      alert(error.text());
      console.log(error.text());
    }
    );
  }

  getFooterInfo() {
    this.homeService.getInfo()
    .subscribe(data => {
      this.version = data;
      this.userIn = true;
    },
    err => console.error(err),
    () =>  {}
    );
  }
}
