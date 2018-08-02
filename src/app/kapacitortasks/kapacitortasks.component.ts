import { Component, ChangeDetectionStrategy, ViewChild, OnInit } from '@angular/core';
import { FormBuilder, Validators} from '@angular/forms';
import { FormArray, FormGroup, FormControl} from '@angular/forms';

import { KapacitorTasksService } from './kapacitortasks.service';
import { ValidationService } from '../common/custom-validation/validation.service'
import { ExportServiceCfg } from '../common/dataservice/export.service'
import { ExportFileModal } from '../common/dataservice/export-file-modal';

import { GenericModal } from '../common/custom-modal/generic-modal';
import { Observable } from 'rxjs/Rx';

import { TableListComponent } from '../common/table-list.component';
import { KapacitorTasksComponentRt, TableRole, OverrideRoleActions } from './kapacitortasks.data';

declare var _:any;

@Component({
  selector: 'kapacitortasks-component',
  providers: [KapacitorTasksService, ValidationService],
  templateUrl: './kapacitortasks.component.html',
  styleUrls: ['../../css/component-styles.css']
})

export class KapacitorTasksComponent implements OnInit {
  @ViewChild('viewModal') public viewModal: GenericModal;
  @ViewChild('viewModalDelete') public viewModalDelete: GenericModal;
  @ViewChild('listTableComponent') public listTableComponent: TableListComponent;
  @ViewChild('exportFileModal') public exportFileModal : ExportFileModal;


  public editmode: string; //list , create, modify
  public componentList: Array<any>;
  public filter: string;
  public sampleComponentForm: any;
  public alertHandler : any = null;
  public counterItems : number = null;
  public counterErrors: any = [];
  public defaultConfig : any = KapacitorTasksComponentRt;
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

  constructor(public kapacitorTasksService: KapacitorTasksService, public exportServiceCfg : ExportServiceCfg, builder: FormBuilder) {
    this.builder = builder;
  }

  createStaticForm() {
    this.sampleComponentForm = this.builder.group({
      ID: [this.sampleComponentForm ? this.sampleComponentForm.value.ID : '', Validators.required],
      Description: [this.sampleComponentForm ? this.sampleComponentForm.value.Description : '']
    });
  }

  reloadData() {
    // now it's a simple subscription to the observable
    this.alertHandler = null;
    this.kapacitorTasksService.getKapacitorTasksItem(null)
      .subscribe(
      data => {
        console.log(data);
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
      case 'new' :
        this.newItem()
      break;
      case 'export' :
        this.exportItem(action.event);
      break;
      case 'view':
        this.viewItem(action.event);
      break;
        case 'tableaction':
        this.applyAction(action.event, action.data);
      break;
    }
  }


  applyAction(action : any, data? : Array<any>) : void {
    this.selectedArray = data || [];
    switch(action.action) {
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

  newItem() {
    //No hidden fields, so create fixed Form
    this.alertHandler =  null;
    this.createStaticForm();
    this.editmode = "create";
  }

  cancelEdit() {
    this.editmode = "list";
    this.reloadData();
  }

  genericForkJoin(obsArray: any) {
    Observable.forkJoin(obsArray)
              .subscribe(
                data => {
                  this.selectedArray = [];
                  this.reloadData()
                },
                err => console.error(err),
              );
  }

  createMultiselectArray(tempArray) : any {
    let myarray = [];
    for (let entry of tempArray) {
      myarray.push({ 'id': entry.ID, 'name': entry.ID, 'extraData': entry.Description });
    }
    return myarray;
  }

}
