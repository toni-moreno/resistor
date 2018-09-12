import { Component, ChangeDetectionStrategy, ViewChild, OnInit } from '@angular/core';
import { FormBuilder, Validators} from '@angular/forms';
import { FormArray, FormGroup, FormControl} from '@angular/forms';

import { AlertEventService } from './alertevent.service';
import { ValidationService } from '../common/custom-validation/validation.service'
import { ExportServiceCfg } from '../common/dataservice/export.service'
import { ExportFileModal } from '../common/dataservice/export-file-modal';

import { GenericModal } from '../common/custom-modal/generic-modal';
import { Observable } from 'rxjs/Rx';

import { TableListComponent } from '../common/table-list.component';
import { AlertEventComponentConfig, TableRole, OverrideRoleActions } from './alertevent.data';

declare var _:any;

@Component({
  selector: 'alertevent-component',
  providers: [AlertEventService, ValidationService],
  templateUrl: './alertevent.component.html',
  styleUrls: ['../../css/component-styles.css']
})

export class AlertEventComponent implements OnInit {
  @ViewChild('viewModal') public viewModal: GenericModal;
  @ViewChild('viewModalDelete') public viewModalDelete: GenericModal;
  @ViewChild('listTableComponent') public listTableComponent: TableListComponent;
  @ViewChild('exportFileModal') public exportFileModal : ExportFileModal;


  public editmode: string; //list , create, modify
  public componentList: Array<any>;
  public filter: string;
  public sampleComponentForm: any;
  public counterItems : number = null;
  public counterErrors: any = [];
  public defaultConfig : any = AlertEventComponentConfig;
  public tableRole : any = TableRole;
  public overrideRoleActions: any = OverrideRoleActions;
  public selectedArray : any = [];

  public data : Array<any>;
  public isRequesting : boolean;

  private builder;
  private oldID : string;

  ngOnInit() {
    this.editmode = 'list';
    this.reloadData();
  }

  constructor(public alertEventService: AlertEventService, public exportServiceCfg : ExportServiceCfg, builder: FormBuilder) {
    this.builder = builder;
  }

  createStaticForm() {
    this.sampleComponentForm = this.builder.group({
      ID: [this.sampleComponentForm ? this.sampleComponentForm.value.ID: '', Validators.required],
      AlertID: [this.sampleComponentForm ? this.sampleComponentForm.value.AlertID: '', Validators.required],
      Message: [this.sampleComponentForm ? this.sampleComponentForm.value.Message: '', Validators.required],
      Details: [this.sampleComponentForm ? this.sampleComponentForm.value.Details: '', Validators.required],
      Time: [this.sampleComponentForm ? this.sampleComponentForm.value.Time: '', Validators.required],
      Duration: [this.sampleComponentForm ? this.sampleComponentForm.value.Duration: '', Validators.required],
      Level: [this.sampleComponentForm ? this.sampleComponentForm.value.Level: '', Validators.required],
      PreviousLevel: [this.sampleComponentForm ? this.sampleComponentForm.value.PreviousLevel: '', Validators.required]
    });
  }

  reloadData(action? : any) {
    this.alertEventService.getAlertEventWithParams(action)
      .subscribe(
      data => {
        this.isRequesting = false;
        this.componentList = data;
        this.data = data;
        this.editmode = "list";
      },
      err => console.error(err),
      () => console.log('DONE')
      );
  }

  customActions(action : any) {
    switch (action.option) {
      case 'view':
        this.viewItem(action.event);
      break;
      case 'edit':
        this.editSampleItem(action.event);
      break;
      case 'remove':
        this.removeItem(action.event);
      break;
      case 'tableaction':
        this.applyAction(action.event, action.data);
      break;
      case 'reloaddata':
        this.reloadData(action);
      break;
    }
  }

  applyAction(action : any, data? : Array<any>) : void {
    this.selectedArray = data || [];
    switch(action.action) {
       case "RemoveAllSelected": {
          this.removeAllSelectedItems(this.selectedArray);
          break;
       }
       default: {
          break;
       }
    }
  }

  viewItem(id) {
    this.viewModal.parseObject(id);
  }

  exportItem(item : any) : void {
    this.exportFileModal.initExportModal(item);
  }

  removeAllSelectedItems(myArray) {
    let obsArray = [];
    this.counterItems = 0;
    this.isRequesting = true;
    let idList = "";
    for (let i in myArray) {
      idList = idList + "'" + myArray[i].ID + "',";
    }
    if (idList.length > 0) idList = idList.slice(0, -1);
    this.deleteSampleItem(idList,false);
  }

  removeItem(row) {
    let id = row.ID;
    this.alertEventService.checkOnDeleteAlertEventItem(id)
      .subscribe(
        data => {
        this.viewModalDelete.parseObject(data)
      },
      err => console.error(err),
      () => { }
      );
  }

  editSampleItem(row) {
    let id = row.ID;
    this.alertEventService.getAlertEventItemById(id)
      .subscribe(data => {
        this.sampleComponentForm = {};
        this.sampleComponentForm.value = data;
        this.oldID = data.ID
        this.createStaticForm();
        this.editmode = "modify";
      },
      err => console.error(err)
      );
 	}

  deleteSampleItem(id, recursive?) {
    if (!recursive) {
    this.alertEventService.deleteAlertEventItem(id)
      .subscribe(data => { },
      err => console.error(err),
      () => { this.viewModalDelete.hide(); this.reloadData() }
      );
    } else {
      return this.alertEventService.deleteAlertEventItem(id)
      .do(
        (test) =>  { this.counterItems++; console.log(this.counterItems)},
        (err) => { this.counterErrors.push({'ID': id, 'error' : err})}
      );
    }
  }

  cancelEdit() {
    this.editmode = "list";
    this.reloadData();
  }

  createMultiselectArray(tempArray) : any {
    let myarray = [];
    for (let entry of tempArray) {
      myarray.push({ 'id': entry.ID, 'name': entry.ID, 'extraData': entry.Description });
    }
    return myarray;
  }

}
