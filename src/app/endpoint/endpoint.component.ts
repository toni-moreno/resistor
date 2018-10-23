import { Component, ChangeDetectionStrategy, ViewChild, OnInit } from '@angular/core';
import { FormBuilder, Validators} from '@angular/forms';
import { FormArray, FormGroup, FormControl} from '@angular/forms';

import { EndpointService } from './endpoint.service';
import { ValidationService } from '../common/custom-validation/validation.service'
import { ExportServiceCfg } from '../common/dataservice/export.service'
import { ExportFileModal } from '../common/dataservice/export-file-modal';

import { GenericModal } from '../common/custom-modal/generic-modal';
import { Observable } from 'rxjs/Rx';

import { TableListComponent } from '../common/table-list.component';
import { EndpointComponentConfig, TableRole, OverrideRoleActions } from './endpoint.data';

declare var _:any;

@Component({
  selector: 'endpoint-component',
  providers: [EndpointService, ValidationService],
  templateUrl: './endpoint.component.html',
  styleUrls: ['../../css/component-styles.css']
})

export class EndpointComponent implements OnInit {
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
  public defaultConfig : any = EndpointComponentConfig;
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

  constructor(public endpointService: EndpointService, public exportServiceCfg : ExportServiceCfg, builder: FormBuilder) {
    this.builder = builder;
  }

  createStaticForm() {
    this.sampleComponentForm = this.builder.group({
      ID: [this.sampleComponentForm ? this.sampleComponentForm.value.ID : '', Validators.required],
      Type: [this.sampleComponentForm ? this.sampleComponentForm.value.Type : '', Validators.required],
      Enabled: [this.sampleComponentForm ? this.sampleComponentForm.value.Enabled : '', Validators.required],
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
      case 'httppost':
      controlArray.push({'ID': 'URL', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'Headers', 'defVal' : '' });
      controlArray.push({'ID': 'BasicAuthUsername', 'defVal' : '' });
      controlArray.push({'ID': 'BasicAuthPassword', 'defVal' : '' });
      break;
      case 'logging':
      controlArray.push({'ID': 'LogFile', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'LogLevel', 'defVal' : '', 'Validators' : Validators.required });
      break;
      case 'slack':
      controlArray.push({'ID':'URL','defVal' :'', 'Validators' : Validators.required });
      controlArray.push({'ID':'Channel','defVal' :'','Validators' : Validators.required });
      controlArray.push({'ID':'SlackUsername','defVal' :'','Validators' : Validators.required });
      controlArray.push({'ID':'IconEmoji','defVal' :'' });
      controlArray.push({'ID':'SslCa','defVal' :'' });
      controlArray.push({'ID':'SslCert','defVal' :'' });
      controlArray.push({'ID':'SslKey','defVal' :'' });
      controlArray.push({'ID':'InsecureSkipVerify','defVal' :'' });
      break;
      case 'email':
      controlArray.push({'ID':'Host','defVal' :'', 'Validators' : Validators.required });
      controlArray.push({'ID':'Port','defVal' :'','Validators' : Validators.required });
      controlArray.push({'ID':'Username','defVal' :'' });
      controlArray.push({'ID':'Password','defVal' :'' });
      controlArray.push({'ID':'From','defVal' :'','Validators' : Validators.required });
      controlArray.push({'ID':'To','defVal' :'','Validators' : Validators.required });
      controlArray.push({'ID':'InsecureSkipVerify','defVal' :'' });
      break;
      default: //Default mode is httppost
      controlArray.push({'ID': 'URL', 'defVal' : '', 'Validators' : Validators.required });
      controlArray.push({'ID': 'Headers', 'defVal' : '' });
      controlArray.push({'ID': 'BasicAuthUsername', 'defVal' : '' });
      controlArray.push({'ID': 'BasicAuthPassword', 'defVal' : '' });
      break;
    }
    //Reload the formGroup with new values saved on controlArray
    this.createDynamicForm(controlArray);
  }

  reloadData() {
    // now it's a simple subscription to the observable
    this.endpointService.getEndpointItem(null)
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
        this.newItem();
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
    this.endpointService.checkOnDeleteEndpointItem(id)
      .subscribe(
        data => {
        this.viewModalDelete.parseObject(data)
      },
      err => console.error(err),
      () => { }
      );
  }
  newItem() {
    if (this.sampleComponentForm) {
      this.setDynamicFields(this.sampleComponentForm.value.Type);
    } else {
      this.setDynamicFields(null);
    }
    this.editmode = "create";
  }

  editSampleItem(row) {
    let id = row.ID;
    this.endpointService.getEndpointItemById(id)
      .subscribe(data => {
        this.sampleComponentForm = {};
        this.sampleComponentForm.value = data;
        this.oldID = data.ID;
        this.setDynamicFields(data.Type);
        this.editmode = "modify";
      },
      err => console.error(err)
      );
 	}

  deleteSampleItem(id, recursive?) {
    if (!recursive) {
    this.endpointService.deleteEndpointItem(id)
      .subscribe(data => { },
      err => console.error(err),
      () => { this.viewModalDelete.hide(); this.reloadData() }
      );
    } else {
      return this.endpointService.deleteEndpointItem(id)
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
      this.endpointService.addEndpointItem(this.sampleComponentForm.value)
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
          r = confirm("Changing Endpoint Instance ID from " + this.oldID + " to " + this.sampleComponentForm.value.ID + ". Proceed?");
        }
        if (r == true) {
          this.endpointService.editEndpointItem(this.sampleComponentForm.value, this.oldID)
            .subscribe(data => { console.log(data) },
            err => console.error(err),
            () => { this.editmode = "list"; this.reloadData() }
            );
        }
      }
    } else {
      return this.endpointService.editEndpointItem(component, component.ID)
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

  createMultiselectArray(tempArray) : any {
    let myarray = [];
    for (let entry of tempArray) {
      myarray.push({ 'id': entry.ID, 'name': entry.ID, 'extraData': entry.Description });
    }
    return myarray;
  }

}
