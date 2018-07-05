import { Component, ChangeDetectionStrategy, ViewChild, OnInit } from '@angular/core';
import { FormBuilder, Validators} from '@angular/forms';
import { FormArray, FormGroup, FormControl} from '@angular/forms';

import { DeviceStatService } from './devicestat.service';
import { AlertService } from '../alert/alert.service';
import { ProductService } from '../product/product.service';

import { ValidationService } from '../common/custom-validation/validation.service'
import { ExportServiceCfg } from '../common/dataservice/export.service'
import { ExportFileModal } from '../common/dataservice/export-file-modal';

import { GenericModal } from '../common/custom-modal/generic-modal';
import { Observable } from 'rxjs/Rx';

import { TableListComponent } from '../common/table-list.component';
import { DeviceStatComponentConfig, TableRole, OverrideRoleActions } from './devicestat.data';
import { IMultiSelectOption, IMultiSelectSettings, IMultiSelectTexts } from '../common/multiselect-dropdown';

declare var _:any;

@Component({
  selector: 'devicestat-component',
  providers: [DeviceStatService, AlertService, ProductService, ValidationService],
  templateUrl: './devicestat.component.html',
  styleUrls: ['../../css/component-styles.css']
})

export class DeviceStatComponent implements OnInit {
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
  public defaultConfig : any = DeviceStatComponentConfig;
  public tableRole : any = TableRole;
  public overrideRoleActions: any = OverrideRoleActions;
  public select_alert : IMultiSelectOption[] = [];
  public select_product : IMultiSelectOption[] = [];
  private single_select: IMultiSelectSettings = {singleSelect: true};

  public selectedArray : any = [];
  public  : any = [];


  public data : Array<any>;
  public isRequesting : boolean;

  private builder;
  private oldID : string;

  ngOnInit() {
    this.editmode = 'list';
    this.reloadData();
  }

  constructor(public devicestatService: DeviceStatService, public alertService: AlertService, public productService: ProductService, public exportServiceCfg: ExportServiceCfg, builder: FormBuilder) {
    this.builder = builder;
  }

  createStaticForm() {
    this.sampleComponentForm = this.builder.group({
      ID: [this.sampleComponentForm ? this.sampleComponentForm.value.ID : ''],
      OrderID: [this.sampleComponentForm ? this.sampleComponentForm.value.OrderID : '', Validators.required],
      DeviceID: [this.sampleComponentForm ? this.sampleComponentForm.value.DeviceID : '', Validators.required],
      AlertID: [this.sampleComponentForm ? this.sampleComponentForm.value.AlertID : '', Validators.required],
      ProductID: [this.sampleComponentForm ? this.sampleComponentForm.value.ProductID : '', Validators.required],
      ExceptionID: [this.sampleComponentForm ? this.sampleComponentForm.value.ExceptionID : '', Validators.required],
      Active: [this.sampleComponentForm ? this.sampleComponentForm.value.Active : '', Validators.required],
      BaseLine: [this.sampleComponentForm ? this.sampleComponentForm.value.BaseLine : '', Validators.required],
      FilterTagKey: [this.sampleComponentForm ? this.sampleComponentForm.value.FilterTagKey : ''],
      FilterTagValue: [this.sampleComponentForm ? this.sampleComponentForm.value.FilterTagValue : ''],
      Description: [this.sampleComponentForm ? this.sampleComponentForm.value.Description : '']
    });
  }

  reloadData() {
    // now it's a simple subscription to the observable
    this.alertHandler = null;
    this.devicestatService.getDeviceStatItem(null)
      .subscribe(
      data => {
        this.isRequesting = false;
        this.componentList = data
        this.data = data;
        this.editmode = "list";
      },
      err => console.error(err),
      () => console.log('DONE')
      );
  }

  customActions(action : any) {
    console.log(action);
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
      case 'edit':
        this.editSampleItem(action.event);
      break;
      case 'remove':
        this.removeItem(action.event);
      break;
      case 'tableaction':
        this.applyAction(action.event, action.data);
      break;
    }
  }

  applyAction(action : any, data? : Array<any>) : void {
    console.log(action);
    this.selectedArray = data || [];
    console.log(this.selectedArray);
    switch(action.action) {
       case "RemoveAllSelected": {
          this.removeAllSelectedItems(this.selectedArray);
          break;
       }
       case "ChangeProperty": {
          this.updateAllSelectedItems(this.selectedArray,action.field,action.value)
          break;
       }
       case "AppendProperty": {
         this.updateAllSelectedItems(this.selectedArray,action.field,action.value,true);
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
    for (let i in myArray) {
      console.log("Removing ",myArray[i].ID)
      this.deleteSampleItem(myArray[i].ID,true);
      obsArray.push(this.deleteSampleItem(myArray[i].ID,true));
    }
    this.genericForkJoin(obsArray);
    console.log(this.counterItems);
  }

  removeItem(row) {
    let id = row.ID;
    console.log('remove', id);
    this.devicestatService.checkOnDeleteDeviceStatItem(id)
      .subscribe(
        data => {
          console.log(data);
        this.viewModalDelete.parseObject(data)
      },
      err => console.error(err),
      () => { }
      );
  }
  newItem() {
    //No hidden fields, so create fixed Form
    this.getAlertItem();
    this.getProductItem();
    this.createStaticForm();
    this.editmode = "create";
  }

  editSampleItem(row) {
    let id = row.ID;
    this.getAlertItem();
    this.getProductItem();
    this.devicestatService.getDeviceStatItemById(id)
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
    this.devicestatService.deleteDeviceStatItem(id)
      .subscribe(data => { },
      err => console.error(err),
      () => { this.viewModalDelete.hide(); this.reloadData() }
      );
    } else {
      return this.devicestatService.deleteDeviceStatItem(id)
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

  saveSampleItem() {
    console.log("SAVE");
    if (this.sampleComponentForm.valid) {
      this.devicestatService.addDeviceStatItem(this.sampleComponentForm.value)
        .subscribe(data => { console.log(data) },
        err => {
          console.log(err);
        },
        () => { this.editmode = "list"; this.reloadData() }
        );
    }
  }

  updateAllSelectedItems(mySelectedArray,field,value, append?) {
    let obsArray = [];
    this.counterItems = 0;
    this.isRequesting = true;
    if (!append)
    for (let component of mySelectedArray) {
      component[field] = value;
      obsArray.push(this.updateSampleItem(true,component));
    } else {
      let tmpArray = [];
      if(!Array.isArray(value)) value = value.split(',');
      console.log(value);
      for (let component of mySelectedArray) {
        console.log(value);
        //check if there is some new object to append
        let newEntries = _.differenceWith(value,component[field],_.isEqual);
        tmpArray = newEntries.concat(component[field])
        console.log(tmpArray);
        component[field] = tmpArray;
        obsArray.push(this.updateSampleItem(true,component));
      }
    }
    this.genericForkJoin(obsArray);
    //Make sync calls and wait the result
    this.counterErrors = [];
  }

  updateSampleItem(recursive?, component?) {
    if(!recursive) {
      if (this.sampleComponentForm.valid) {
        var r = true;
        if (this.sampleComponentForm.value.ID != this.oldID) {
          r = confirm("Changing DeviceStat Instance ID from " + this.oldID + " to " + this.sampleComponentForm.value.ID + ". Proceed?");
        }
        if (r == true) {
          this.devicestatService.editDeviceStatItem(this.sampleComponentForm.value, this.oldID)
            .subscribe(data => { console.log(data) },
            err => console.error(err),
            () => { this.editmode = "list"; this.reloadData() }
            );
        }
      }
    } else {
      return this.devicestatService.editDeviceStatItem(component, component.ID)
      .do(
        (test) =>  { this.counterItems++ },
        (err) => { this.counterErrors.push({'ID': component['ID'], 'error' : err['_body']})}
      )
      .catch((err) => {
        return Observable.of({'ID': component.ID , 'error': err['_body']})
      })
    }
  }


  testSampleItemConnection() {
    this.devicestatService.testDeviceStatItem(this.sampleComponentForm.value)
    .subscribe(
    data =>  this.alertHandler = {msg: 'DeviceStat Version: '+data['Message'], result : data['Result'], elapsed: data['Elapsed'], type: 'success', closable: true},
    err => {
        let error = err.json();
        this.alertHandler = {msg: error['Message'], elapsed: error['Elapsed'], result : error['Result'], type: 'danger', closable: true}
      },
    () =>  { console.log("DONE")}
  );

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

  getAlertItem() {
    this.alertService.getAlertItem(null)
      .subscribe(
      data => {
        this.select_alert = [];
        this.select_alert = this.createMultiselectArray(data, 'ID', 'ID');
      },
      err => console.error(err),
      () => console.log('DONE')
      );
  }

  getProductItem() {
    this.productService.getProductItem(null)
      .subscribe(
      data => {
        this.select_product = [];
        this.select_product = this.createMultiselectArray(data, 'ID', 'ID');
      },
      err => console.error(err),
      () => console.log('DONE')
      );
  }

  createMultiselectArray(tempArray, ID?, Name?, extraData?) : any {
    let myarray = [];
    if(tempArray){
      for (let entry of tempArray) {
        myarray.push({ 'id': ID ? entry[ID] : entry, 'name': Name ? entry[Name] : entry, 'extraData': extraData ? entry[extraData] : null });
      };
    }
    return myarray;
  }

}
