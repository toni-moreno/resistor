import { Component, ChangeDetectionStrategy, ViewChild, OnInit } from '@angular/core';
import { FormBuilder, Validators} from '@angular/forms';
import { FormArray, FormGroup, FormControl} from '@angular/forms';

import { AlertService } from './alert.service';
import { ProductService } from '../product/product.service';
import { RangeTimeService } from '../rangetime/rangetime.service';
import { OutHTTPService } from '../outhttp/outhttp.service';
import { KapacitorService } from '../kapacitor/kapacitor.service';
import { IfxDBService } from '../ifxdb/ifxdb.service';
import { IfxMeasurementService } from '../ifxmeasurement/ifxmeasurement.service';

import { ValidationService } from '../common/validation.service'
import { ExportServiceCfg } from '../common/dataservice/export.service'

import { GenericModal } from '../common/generic-modal';
import { Observable } from 'rxjs/Rx';

import { TableListComponent } from '../common/table-list.component';
import { AlertComponentConfig } from './alert.data';

import { IMultiSelectOption, IMultiSelectSettings, IMultiSelectTexts } from '../common/multiselect-dropdown';


declare var _:any;

@Component({
  selector: 'alert-component',
  providers: [AlertService, ProductService, RangeTimeService, OutHTTPService, KapacitorService, IfxDBService,IfxMeasurementService, ValidationService],
  templateUrl: './alert.component.html',
  styleUrls: ['../../css/component-styles.css']
})

export class AlertComponent implements OnInit {
  @ViewChild('viewModal') public viewModal: GenericModal;
  @ViewChild('viewModalDelete') public viewModalDelete: GenericModal;
  @ViewChild('listTableComponent') public listTableComponent: TableListComponent;


  public editmode: string; //list , create, modify
  public componentList: Array<any>;
  public filter: string;
  public sampleComponentForm: any;
  public counterItems : number = null;
  public counterErrors: any = [];
  public defaultConfig : any = AlertComponentConfig;
  public selectedArray : any = [];

  public select_product : IMultiSelectOption[] = [];
  public select_rangetime : IMultiSelectOption[] = [];
  public select_outhttp : IMultiSelectOption[] = [];
  public select_kapacitor : IMultiSelectOption[] = [];
  public select_ifxdb : IMultiSelectOption[] = [];
  public select_ifxrp : IMultiSelectOption[] = [];
  public select_ifxms : IMultiSelectOption[] = [];
  public select_ifxfs : IMultiSelectOption[] = [];
  public select_ifxts : IMultiSelectOption[] = [];


  public ifxdb_list : any = [];
  public picked_ifxdb: any = null;
  public picked_ifxms: any = null;

  private single_select: IMultiSelectSettings = {
      singleSelect: true,
  };

  public data : Array<any>;
  public isRequesting : boolean;

  private builder;
  private oldID : string;

  ngOnInit() {
    this.editmode = 'list';
    this.reloadData();
  }

  constructor(public alertService: AlertService,public productService :ProductService, public rangetimeService : RangeTimeService, public ifxDBService : IfxDBService, public ifxMeasurementService : IfxMeasurementService, public outhttpService: OutHTTPService, public kapacitorService: KapacitorService, public exportServiceCfg : ExportServiceCfg, builder: FormBuilder) {
    this.builder = builder;
  }

  createStaticForm() {
    this.sampleComponentForm = this.builder.group({
      ID: [this.sampleComponentForm ? this.sampleComponentForm.value.ID : '', Validators.required],
      BaselineID: [this.sampleComponentForm ? this.sampleComponentForm.value.BaselineID : '', Validators.required],
      ProductID: [this.sampleComponentForm ? this.sampleComponentForm.value.ProductID : '', Validators.required],
      GroupID: [this.sampleComponentForm ? this.sampleComponentForm.value.GroupID : '', Validators.required],
      NumAlertID: [this.sampleComponentForm ? this.sampleComponentForm.value.NumAlertID : '', Validators.required],
      TrigerType: [this.sampleComponentForm ? this.sampleComponentForm.value.TrigerType : 'THRESHOLD', Validators.required],
      InfluxDB: [this.sampleComponentForm ? this.sampleComponentForm.value.InfluxDB : null, Validators.required],
      InfluxRP: [this.sampleComponentForm ? this.sampleComponentForm.value.InfluxRP : null, Validators.required],
      InfluxMeasurement: [this.sampleComponentForm ? this.sampleComponentForm.value.InfluxMeasurement : null, Validators.required],
      TagDescription: [this.sampleComponentForm ? this.sampleComponentForm.value.TagDescription : ''],
      InfluxFilter: [this.sampleComponentForm ? this.sampleComponentForm.value.InfluxFilter : '', Validators.required],
      IntervalCheck: [this.sampleComponentForm ? this.sampleComponentForm.value.IntervalCheck : '', Validators.required],
      OperationID: [this.sampleComponentForm ? this.sampleComponentForm.value.OperationID : ''],
      Field: [this.sampleComponentForm ? this.sampleComponentForm.value.Field : '', Validators.required],
      GrafanaServer: [this.sampleComponentForm ? this.sampleComponentForm.value.GrafanaServer : ''],
      GrafanaDashLabel: [this.sampleComponentForm ? this.sampleComponentForm.value.GrafanaDashLabel : ''],
      GrafanaDashPanelID: [this.sampleComponentForm ? this.sampleComponentForm.value.GrafanaDashPanelID : ''],
      DeviceIDTag: [this.sampleComponentForm ? this.sampleComponentForm.value.DeviceIDTag : ''],
      DeviceIDLabel: [this.sampleComponentForm ? this.sampleComponentForm.value.DeviceIDLabel : ''],
      ExtraTag: [this.sampleComponentForm ? this.sampleComponentForm.value.ExtraTag : ''],
      ExtraLabel: [this.sampleComponentForm ? this.sampleComponentForm.value.ExtraLabel : ''],
      KapacitorID: [this.sampleComponentForm ? this.sampleComponentForm.value.KapacitorID : '', Validators.required],
      OutHTTP: [this.sampleComponentForm ? this.sampleComponentForm.value.OutHTTP : '', Validators.required],
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
        if (tmpform[entry.ID] && entry.override !== true) {
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
      controlArray.push({'ID': 'ThresholdType', 'defVal' : 'absolute', 'Validators' : Validators.required });
      controlArray.push({'ID': 'StatFunc', 'defVal' : 'MEAN', 'Validators' : Validators.required });
      controlArray.push({'ID': 'CritDirection', 'defVal' : 'CC', 'Validators' : Validators.required });
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

      case 'TREND':
      controlArray.push({'ID': 'StatFunc', 'defVal' : 'MEAN', 'Validators' : Validators.required });
      controlArray.push({'ID': 'CritDirection', 'defVal' : 'CC', 'Validators' : Validators.required });
      controlArray.push({'ID': 'Shift', 'defVal' : '', 'Validators' : Validators.required });
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
      break
      default: //Default mode is THRESHOLD
      controlArray.push({'ID': 'StatFunc', 'defVal' : 'MEAN', 'Validators' : Validators.required });
      controlArray.push({'ID': 'ThresholdType', 'defVal' : 'absolute', 'Validators' : Validators.required });
      controlArray.push({'ID': 'CritDirection', 'defVal' : 'CC', 'Validators' : Validators.required });
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
    this.getOutHTTPItem();
    this.getKapacitorItem();
    this.getIfxDBItem();
    if (this.sampleComponentForm) {
      this.setDynamicFields(this.sampleComponentForm.value.TigerType);
    } else {
      this.setDynamicFields(null);
    }
    this.editmode = "create";
  }

  editSampleItem(row) {
    let id = row.ID;
    this.getProductItem();
    this.getRangeTimeItem();
    this.getOutHTTPItem();
    this.getKapacitorItem();
    this.getIfxDBItem();
    this.alertService.getAlertItemById(id)
      .subscribe(data => {
        this.sampleComponentForm = {};
        this.sampleComponentForm.value = data;
        this.oldID = data.ID
        this.setDynamicFields(row.TigerType);
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

  getOutHTTPItem() {
    this.outhttpService.getOutHTTPItem(null)
      .subscribe(
      data => {
        this.select_outhttp = [];
        this.select_outhttp = this.createMultiselectArray(data, 'ID','ID');
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

  getIfxDBItem() {
    this.ifxDBService.getIfxDBItem(null)
      .subscribe(
      data => {
        this.ifxdb_list = data;
        this.select_ifxdb = [];
        this.select_ifxdb = this.createMultiselectArray(data, 'ID','Name');
      },
      err => console.error(err),
      () => console.log('DONE')
      );
  }

  pickDBItem(ifxdb_picked) {

    if (this.picked_ifxdb) {
      if (ifxdb_picked !== this.picked_ifxdb['ID']) {
        this.sampleComponentForm.controls.InfluxRP.setValue(null);
        this.sampleComponentForm.controls.InfluxMeasurement.setValue(null);
      }
    }
    //Clear Vars:
    this.select_ifxrp = null;
    this.picked_ifxdb = this.ifxdb_list.filter((x) => x['ID'] === ifxdb_picked)[0];

    if(this.picked_ifxdb) {
      this.select_ifxrp = this.createMultiselectArray(this.picked_ifxdb['Retention']);
      this.select_ifxms = this.createMultiselectArray(this.picked_ifxdb['Measurements'], 'Name', 'Name','ID');
    }
  }

  pickMeasItem(ifxms_picked) {
    this.sampleComponentForm.controls.Field.setValue(null);
    this.sampleComponentForm.controls.TagDescription.setValue(null);

    if (ifxms_picked){
      this.select_ifxfs = [];
      this.ifxMeasurementService.getIfxMeasurementItemById(_.find(this.select_ifxms, {'id' : ifxms_picked})['extraData'])
      .subscribe(
        data => {
          console.log(data);
          this.select_ifxfs = this.createMultiselectArray(data.Fields);
          this.select_ifxts = this.createMultiselectArray(data.Tags);
        },
        err => console.error(err),
        () => console.log('DONE')
      );
    }
  }

  createMultiselectArray(tempArray, ID?, Name?, extraData?) : any {
    let myarray = [];
    if(tempArray)
    for (let entry of tempArray) {
      myarray.push({ 'id': ID ? entry[ID] : entry, 'name': Name ? entry[Name] : entry, 'extraData': extraData ? entry[extraData] : null });
    };
    return myarray;
  }

}
