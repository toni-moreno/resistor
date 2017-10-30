import { Component, ViewChild,ViewContainerRef } from '@angular/core';
import { NgSwitch, NgSwitchCase, NgSwitchDefault } from '@angular/common';
import { Router } from '@angular/router';
import { HttpAPI} from '../common/httpAPI';
import { Observable } from 'rxjs/Observable';
import { BlockUIService } from '../common/blockui/blockui-service';
import { BlockUIComponent } from '../common/blockui/blockui-component'; // error
import { ImportFileModal } from '../common/dataservice/import-file-modal';
import { HomeService } from './home.service';
import { AboutModal } from './about-modal'
import { WindowRef } from '../common/windowref';
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
  api: string;
  item_type: string;
  version: RInfo;
  configurationItems : Array<any> = [
  {'title': 'ConfItem 1', 'selector' : 'conf1'},
  {'title': 'ConfItem 2', 'selector' : 'conf2'},
  {'title': 'ConfItem 3', 'selector' : 'conf3'},
  {'title': 'ConfItem 4', 'selector' : 'conf4'},
  {'title': 'ConfItem 5', 'selector' : 'conf5'},
  {'title': 'ConfItem 6', 'selector' : 'conf6'},
  {'title': 'ConfItem 7', 'selector' : 'conf7'},
  {'title': 'ConfItem 8', 'selector' : 'conf8'},
  ];

  runtimeItems : Array<any> = [
  {'title': 'Agent status', 'selector' : 'runtime1'},
  ];

  mode : boolean = false;
  userIn : boolean = false;

  elapsedReload: string = '';
  lastReload: Date;

  constructor(private winRef: WindowRef,public router: Router, public httpAPI: HttpAPI, private _blocker: BlockUIService, public homeService: HomeService) {
    this.nativeWindow = winRef.nativeWindow;
    this.getFooterInfo();
    this.item_type= "influxservers";
  }

  link(url: string) {
    this.nativeWindow.open(url);
  }


  logout() {
    this.homeService.userLogout()
    .subscribe(
    response => {
      this.router.navigate(['/login']);
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

  clickMenu(selected : string) : void {
    this.item_type = "";
    this.item_type = selected;
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
