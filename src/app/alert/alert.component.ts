import { Component, ChangeDetectionStrategy, ViewChild, OnInit } from '@angular/core';
import { FormBuilder, Validators} from '@angular/forms';
import { FormArray, FormGroup, FormControl} from '@angular/forms';

import { AlertService } from './alert.service';
import { ProductService } from '../product/product.service';
import { RangeTimeService } from '../rangetime/rangetime.service';
import { EndpointService } from '../endpoint/endpoint.service';
import { KapacitorService } from '../kapacitor/kapacitor.service';
import { OperationService } from '../operation/operation.service';
import { IfxDBService } from '../ifxdb/ifxdb.service';
import { IfxMeasurementService } from '../ifxmeasurement/ifxmeasurement.service';

import { ValidationService } from '../common/custom-validation/validation.service'
import { ExportServiceCfg } from '../common/dataservice/export.service'
import { ExportFileModal } from '../common/dataservice/export-file-modal';
import { WindowRef } from '../common/windowref';

import { GenericModal } from '../common/custom-modal/generic-modal';
import { Observable } from 'rxjs/Rx';

import { TableListComponent } from '../common/table-list.component';
import { AlertComponentConfig, TableRole, OverrideRoleActions } from './alert.data';

import { IMultiSelectOption, IMultiSelectSettings, IMultiSelectTexts } from '../common/multiselect-dropdown';


declare var _:any;

@Component({
  selector: 'alert-component',
  providers: [AlertService, ProductService, RangeTimeService, EndpointService, KapacitorService, OperationService, IfxDBService,IfxMeasurementService, ValidationService],
  templateUrl: './alert.component.html',
  styleUrls: ['../../css/component-styles.css']
})

export class AlertComponent implements OnInit {
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
  public defaultConfig : any = AlertComponentConfig;
  public tableRole : any = TableRole;
  public overrideRoleActions: any = OverrideRoleActions;
  public selectedArray : any = [];

  public select_product : IMultiSelectOption[] = [];
  public select_rangetime : IMultiSelectOption[] = [];
  public select_endpoint : IMultiSelectOption[] = [];
  public select_kapacitor : IMultiSelectOption[] = [];
  public select_operation : IMultiSelectOption[] = [];
  public select_ifxdb : IMultiSelectOption[] = [];
  public select_ifxrp : IMultiSelectOption[] = [];
  public select_ifxms : IMultiSelectOption[] = [];
  public select_ifxfs : IMultiSelectOption[] = [];
  public select_ifxts : IMultiSelectOption[] = [];
  public select_baseline : IMultiSelectOption[] = [];
  public select_alertgroup : IMultiSelectOption[] = [];
  public select_fieldresolution : IMultiSelectOption[] = [];
  public select_idtag : IMultiSelectOption[] = [];


  public ifxdb_list : any = [];
  public picked_ifxdb: any = null;
  public picked_ifxms: any = null;

  public product_list : any = [];
  public picked_product: any = null;

  private single_select: IMultiSelectSettings = {
      singleSelect: true,
  };

  public data : Array<any>;
  public isRequesting : boolean;

  nativeWindow: any
  private builder;
  private oldID : string;

  ngOnInit() {
    this.editmode = 'list';
    this.reloadData();
  }

  constructor(private winRef: WindowRef,public alertService: AlertService,public productService :ProductService, public rangetimeService : RangeTimeService, public ifxDBService : IfxDBService, public ifxMeasurementService : IfxMeasurementService, public endpointService: EndpointService, public kapacitorService: KapacitorService, public operationService: OperationService, public exportServiceCfg : ExportServiceCfg, builder: FormBuilder) {
    this.nativeWindow = winRef.nativeWindow;
    this.builder = builder;
  }

  link(url: string) {
    this.nativeWindow.open(url);
  }

  createStaticForm() {
    this.sampleComponentForm = this.builder.group({
      ID: [this.sampleComponentForm ? this.sampleComponentForm.value.ID : '', Validators.required],
      Active: [this.sampleComponentForm ? this.sampleComponentForm.value.Active : '', Validators.required],
      BaselineID: [this.sampleComponentForm ? this.sampleComponentForm.value.BaselineID : '', Validators.required],
      ProductID: [this.sampleComponentForm ? this.sampleComponentForm.value.ProductID : '', Validators.required],
      AlertGroup: [this.sampleComponentForm ? this.sampleComponentForm.value.AlertGroup : '', Validators.required],
      NumAlertID: [this.sampleComponentForm ? this.sampleComponentForm.value.NumAlertID : '', Validators.compose([Validators.required, ValidationService.uintegerNotZeroValidator])],
      TriggerType: [this.sampleComponentForm ? this.sampleComponentForm.value.TriggerType : 'THRESHOLD', Validators.required],
      InfluxDB: [this.sampleComponentForm ? this.sampleComponentForm.value.InfluxDB : null, Validators.required],
      InfluxRP: [this.sampleComponentForm ? this.sampleComponentForm.value.InfluxRP : null, Validators.required],
      InfluxMeasurement: [this.sampleComponentForm ? this.sampleComponentForm.value.InfluxMeasurement : null, Validators.required],
      InfluxFilter: [this.sampleComponentForm ? this.sampleComponentForm.value.InfluxFilter : ''],
      IntervalCheck: [this.sampleComponentForm ? this.sampleComponentForm.value.IntervalCheck : '', Validators.compose([Validators.required, ValidationService.durationValidator])],
      AlertFrequency: [this.sampleComponentForm ? this.sampleComponentForm.value.AlertFrequency : '', ValidationService.durationValidator],
      AlertNotify: [this.sampleComponentForm ? this.sampleComponentForm.value.AlertNotify : '', ValidationService.uintegerValidator],
      OperationID: [this.sampleComponentForm ? this.sampleComponentForm.value.OperationID : ''],
      IsCustomExpression: [this.sampleComponentForm ? this.sampleComponentForm.value.IsCustomExpression : false, Validators.required],
      Field: [this.sampleComponentForm ? this.sampleComponentForm.value.Field : ''],
      FieldDesc: [this.sampleComponentForm ? this.sampleComponentForm.value.FieldDesc : ''],
      GrafanaServer: [this.sampleComponentForm ? this.sampleComponentForm.value.GrafanaServer : ''],
      GrafanaDashLabel: [this.sampleComponentForm ? this.sampleComponentForm.value.GrafanaDashLabel : ''],
      GrafanaDashPanelID: [this.sampleComponentForm ? this.sampleComponentForm.value.GrafanaDashPanelID : ''],
      ProductTag: [this.sampleComponentForm ? this.sampleComponentForm.value.ProductTag : '', Validators.required],
      ProductTagRO: [this.sampleComponentForm ? this.sampleComponentForm.value.ProductTagRO : ''],
      DeviceIDLabel: [this.sampleComponentForm ? this.sampleComponentForm.value.DeviceIDLabel : ''],
      ExtraTag: [this.sampleComponentForm ? this.sampleComponentForm.value.ExtraTag : ''],
      ExtraLabel: [this.sampleComponentForm ? this.sampleComponentForm.value.ExtraLabel : ''],
      AlertExtraText: [this.sampleComponentForm ? this.sampleComponentForm.value.AlertExtraText : ''],
      IDTag: [this.sampleComponentForm ? this.sampleComponentForm.value.IDTag : ''],
      KapacitorID: [this.sampleComponentForm ? this.sampleComponentForm.value.KapacitorID : '', Validators.required],
      Endpoint: [this.sampleComponentForm ? this.sampleComponentForm.value.Endpoint : ''],
      Description: [this.sampleComponentForm ? this.sampleComponentForm.value.Description : '']
    });
  }

  createDynamicForm(fieldsArray: any) : void {

    //Generates the static form:
    //Saves the actual to check later if there are shared values
    let tmpform : any;
    if (this.sampleComponentForm)  tmpform = this.sampleComponentForm.value;
    this.createStaticForm();
    //Set new values and check if we have to mantain the value!
    for (let entry of fieldsArray) {
      let value = entry.defVal;
      //Check if there are common values from the previous selected item
      if (tmpform) {
        if (tmpform[entry.ID] !== null && entry.override !== true) {
          value = tmpform[entry.ID];
        }
      }
      //Set different controls:
      this.sampleComponentForm.addControl(entry.ID, new FormControl(value, entry.Validators));
    }
}

  setDynamicFields (field : any, override? : boolean) : void  {
    //Saves on the array all values to push into formGroup
    let controlArray : Array<any> = [];
    switch (field) {
      case 'THRESHOLD':
      controlArray.push({'ID': 'CritDirection', 'defVal' : 'AC', 'Validators' : Validators.required });
      controlArray.push({'ID': 'FieldType', 'defVal' : 'GAUGE', 'Validators' : Validators.required });
      controlArray.push({'ID': 'StatFunc', 'defVal' : 'MEAN', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ExtraData', 'defVal' : '' });
      controlArray.push({'ID': 'Rate', 'defVal' : '' });
      controlArray.push({'ID': 'FieldResolution', 'defVal' : '' });
      controlArray.push({'ID': 'ThCritDef', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThCritEx1', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThCritEx2', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThCritRangeTimeID', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThWarnDef', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThWarnEx1', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThWarnEx2', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThWarnRangeTimeID', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThInfoDef', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThInfoEx1', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThInfoEx2', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThInfoRangeTimeID', 'defVal' : '', 'Validators' : Validators.required });
      break;
      case 'TREND':
      controlArray.push({'ID': 'CritDirection', 'defVal' : 'AC', 'Validators' : Validators.required });
      controlArray.push({'ID': 'TrendType', 'defVal' : 'absolute', 'Validators' : Validators.required });
      controlArray.push({'ID': 'TrendSign', 'defVal' : 'positive', 'Validators' : Validators.required });
      controlArray.push({'ID': 'FieldType', 'defVal' : 'GAUGE', 'Validators' : Validators.required });
      controlArray.push({'ID': 'StatFunc', 'defVal' : 'MEAN', 'Validators' : Validators.required });
      controlArray.push({'ID': 'Shift', 'defVal' : '', 'Validators' : Validators.compose([Validators.required, ValidationService.durationValidator]) });
      controlArray.push({'ID': 'ExtraData', 'defVal' : '' });
      controlArray.push({'ID': 'Rate', 'defVal' : '' });
      controlArray.push({'ID': 'FieldResolution', 'defVal' : '' });
      controlArray.push({'ID': 'ThCritDef', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThCritEx1', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThCritEx2', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThCritRangeTimeID', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThWarnDef', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThWarnEx1', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThWarnEx2', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThWarnRangeTimeID', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThInfoDef', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThInfoEx1', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThInfoEx2', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThInfoRangeTimeID', 'defVal' : '', 'Validators' : Validators.required });
      break;
      case 'DEADMAN':
      break;
      default: //Default mode is THRESHOLD
      controlArray.push({'ID': 'CritDirection', 'defVal' : 'AC', 'Validators' : Validators.required });
      controlArray.push({'ID': 'FieldType', 'defVal' : 'GAUGE', 'Validators' : Validators.required });
      controlArray.push({'ID': 'StatFunc', 'defVal' : 'MEAN', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ExtraData', 'defVal' : '' });
      controlArray.push({'ID': 'Rate', 'defVal' : '' });
      controlArray.push({'ID': 'FieldResolution', 'defVal' : '' });
      controlArray.push({'ID': 'ThCritDef', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThCritEx1', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThCritEx2', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThCritRangeTimeID', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThWarnDef', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThWarnEx1', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThWarnEx2', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThWarnRangeTimeID', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThInfoDef', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThInfoEx1', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThInfoEx2', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThInfoRangeTimeID', 'defVal' : '', 'Validators' : Validators.required });
      break;
    }
    //Reload the formGroup with new values saved on controlArray
    this.createDynamicForm(controlArray);
  }


  reloadData() {
    // now it's a simple subscription to the observable
    this.alertService.getAlertItem(null)
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
      case 'deploy':
        this.deployItem(action.event);
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
  }

  removeItem(row) {
    let id = row.ID;
    console.log('remove', id);
    this.alertService.checkOnDeleteAlertItem(id)
      .subscribe(
        data => {
        this.viewModalDelete.parseObject(data)
      },
      err => console.error(err),
      () => { }
      );
  }
  newItem() {
    //Check for subhidden fields
    this.getProductItem();
    this.getRangeTimeItem();
    this.getEndpointItem();
    this.getKapacitorItem();
    this.getOperationItem();
    if (this.sampleComponentForm) {
      this.setDynamicFields(this.sampleComponentForm.value.TriggerType);
    } else {
      this.setDynamicFields(null);
    }
    this.editmode = "create";
  }

  editSampleItem(row) {
    this.picked_product = null;
    this.picked_ifxdb = null;
    let id = row.ID;
    this.getProductItem();
    this.getRangeTimeItem();
    this.getEndpointItem();
    this.getKapacitorItem();
    this.getOperationItem();
    this.alertService.getAlertItemById(id)
      .subscribe(data => {
        this.sampleComponentForm = {};
        this.sampleComponentForm.value = data;
        this.oldID = data.ID
        this.setDynamicFields(data.TriggerType);

        this.editmode = "modify";
      },
      err => console.error(err)
      );
 	}

  deleteSampleItem(id, recursive?) {
    if (!recursive) {
    this.alertService.deleteAlertItem(id)
      .subscribe(data => { },
      err => console.error(err),
      () => { this.viewModalDelete.hide(); this.reloadData() }
      );
    } else {
      return this.alertService.deleteAlertItem(id)
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
    if (this.sampleComponentForm.valid) {
      this.alertService.addAlertItem(this.sampleComponentForm.value)
        .subscribe(data => { console.log(data) },
        err => {
          console.log(err);
        },
        () => { this.editmode = "list"; this.reloadData() }
        );
    }
  }

  deployItem(row) {
    this.alertService.editAlertItem(row, row.ID)
    .subscribe(data => { console.log(data) },
    err => {
      console.log(err);
    },
    () => { this.editmode = "list"; this.reloadData() }
    );
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
      for (let component of mySelectedArray) {
        console.log(value);
        //check if there is some new object to append
        let newEntries = _.differenceWith(value,component[field],_.isEqual);
        tmpArray = newEntries.concat(component[field])
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
          r = confirm("Changing Alert Instance ID from " + this.oldID + " to " + this.sampleComponentForm.value.ID + ". Proceed?");
        }
        if (r == true) {
          this.alertService.editAlertItem(this.sampleComponentForm.value, this.oldID)
            .subscribe(data => { console.log(data) },
            err => console.error(err),
            () => { this.editmode = "list"; this.reloadData() }
            );
        }
      }
    } else {
      return this.alertService.editAlertItem(component, component.ID)
      .do(
        (test) =>  { this.counterItems++ },
        (err) => { this.counterErrors.push({'ID': component['ID'], 'error' : err['_body']})}
      )
      .catch((err) => {
        return Observable.of({'ID': component.ID , 'error': err['_body']})
      })
    }
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

  getProductItem() {
    this.productService.getProductItem(null)
      .subscribe(
      data => {
        this.product_list = data;
        this.select_product = [];
        this.select_product = this.createMultiselectArray(data, 'ID','ID');
      },
      err => console.error(err),
      () => console.log('DONE')
      );
  }


  getRangeTimeItem() {
    this.rangetimeService.getRangeTimeItem(null)
      .subscribe(
      data => {
        this.select_rangetime = [];
        this.select_rangetime = this.createMultiselectArray(data, 'ID','ID');
      },
      err => console.error(err),
      () => console.log('DONE')
      );
  }

  getEndpointItem() {
    this.endpointService.getEndpointItem(null)
      .subscribe(
      data => {
        this.select_endpoint = [];
        this.select_endpoint = this.createMultiselectArray(data, 'ID','ID');
      },
      err => console.error(err),
      () => console.log('DONE')
      );
  }
  getKapacitorItem() {
    this.kapacitorService.getKapacitorItem(null)
      .subscribe(
      data => {
        this.select_kapacitor = [];
        this.select_kapacitor = this.createMultiselectArray(data, 'ID','ID');
      },
      err => console.error(err),
      () => console.log('DONE')
      );
  }

  getOperationItem() {
    this.operationService.getOperationItem(null)
      .subscribe(
      data => {
        this.select_operation = [];
        this.select_operation = this.createMultiselectArray(data, 'ID','ID');
      },
      err => console.error(err),
      () => console.log('DONE')
      );
  }

  getIfxDBItem() {
    this.ifxDBService.getIfxDBItem(null)
      .subscribe(
      data => {
        this.ifxdb_list = data;
        this.select_ifxdb = [];
        this.select_ifxdb = this.createMultiselectArray(data, 'ID','Name','IfxServer');
      },
      err => console.error(err),
      () => console.log('DONE')
      );
  }

  pickDBItem(ifxdb_picked) {

    if (this.picked_ifxdb) {
      if (ifxdb_picked !== this.picked_ifxdb['ID']) {
        this.sampleComponentForm.controls.InfluxRP.setValue(null);
        this.sampleComponentForm.controls.Field.setValue(null);
        this.select_ifxrp = null;
      }
    }
    //Clear Vars:
    this.picked_ifxdb = this.ifxdb_list.filter((x) => x['ID'] === ifxdb_picked)[0];

    if(this.picked_ifxdb) {
      this.select_ifxrp = this.createMultiselectArray(this.picked_ifxdb['Retention']);
      this.select_ifxfs = [];
      this.ifxMeasurementService.getIfxMeasurementItemByDbIdMeasName(ifxdb_picked, this.picked_ifxms)
      .subscribe(
        data => {
          this.select_ifxfs = this.createMultiselectArray(data.Fields);
        },
        err => console.error(err),
        () => console.log('DONE')
      );
    }
  }

  pickMeasItem(ifxms_picked) {
    //Only reset values when default values are loaded
    if (ifxms_picked !== this.sampleComponentForm.value.InfluxMeasurement){
      this.sampleComponentForm.controls.InfluxDB.setValue(null);
      this.sampleComponentForm.controls.InfluxRP.setValue(null);
      this.sampleComponentForm.controls.Field.setValue(null);
    }

    if (ifxms_picked){
      this.picked_ifxms = ifxms_picked;
      this.ifxDBService.getIfxDBCfgArrayByMeasName(ifxms_picked)
      .subscribe(
        data => {
          console.log(data);
          this.ifxdb_list = data;
          this.select_ifxdb = [];
          this.select_ifxdb = this.createMultiselectArray(data, 'ID','Name','IfxServer');
        },
        err => console.error(err),
        () => console.log('DONE')
      );
    }
  }

  pickProductItem(product_picked) {

    if (this.picked_product) {
      if (product_picked !== this.picked_product['ID']) {
        this.sampleComponentForm.controls.ProductTag.setValue(null);
        this.sampleComponentForm.controls.BaselineID.setValue(null);
        this.sampleComponentForm.controls.AlertGroup.setValue(null);
        this.sampleComponentForm.controls.FieldResolution.setValue(null);
        this.sampleComponentForm.controls.InfluxMeasurement.setValue(null);
        this.sampleComponentForm.controls.IDTag.setValue(null);
      }
    }
    //Clear Vars:
    this.picked_product = this.product_list.filter((x) => x['ID'] === product_picked)[0];
    this.select_baseline = null;
    this.select_alertgroup = null;
    this.select_fieldresolution = null;
    this.select_ifxms = null;
    this.select_idtag = null;

    if(this.picked_product) {
      this.sampleComponentForm.controls.ProductTag.setValue(this.picked_product['ProductTag']);
      this.select_baseline = this.createMultiselectArray(this.picked_product['BaseLines']);
      this.select_alertgroup = this.createMultiselectArray(this.picked_product['AlertGroups']);
      this.select_fieldresolution = this.createMultiselectArray(this.picked_product['FieldResolutions']);
      this.select_ifxms = this.createMultiselectArray(this.picked_product['Measurements']);
      this.select_idtag = this.createMultiselectArray(this.picked_product['CommonTags']);
    }
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
