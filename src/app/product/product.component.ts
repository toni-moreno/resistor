import { Component, ChangeDetectionStrategy, ViewChild, OnInit } from '@angular/core';
import { FormBuilder, Validators} from '@angular/forms';
import { FormArray, FormGroup, FormControl} from '@angular/forms';

import { ProductService } from './product.service';
import { IfxMeasurementService } from '../ifxmeasurement/ifxmeasurement.service';
import { ValidationService } from '../common/custom-validation/validation.service'
import { ExportServiceCfg } from '../common/dataservice/export.service'
import { ExportFileModal } from '../common/dataservice/export-file-modal';

import { GenericModal } from '../common/custom-modal/generic-modal';
import { Observable } from 'rxjs/Rx';

import { TableListComponent } from '../common/table-list.component';
import { ProductComponentConfig, TableRole, OverrideRoleActions } from './product.data';

import { IMultiSelectOption, IMultiSelectSettings, IMultiSelectTexts } from '../common/multiselect-dropdown';

declare var _:any;

@Component({
  selector: 'product-component',
  providers: [ProductService, IfxMeasurementService, ValidationService],
  templateUrl: './product.component.html',
  styleUrls: ['../../css/component-styles.css']
})

/*
TODO
When a tag is selected in one of the three taglist fields, 
this tag should be removed from the other two taglist fields.
*/
export class ProductComponent implements OnInit {
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
  public defaultConfig : any = ProductComponentConfig;
  public tableRole : any = TableRole;
  public overrideRoleActions: any = OverrideRoleActions;
  public selectedArray : any = [];

  public ifxms_list : any = [];
  public picked_ifxms: any = null;

  public select_ifxms : IMultiSelectOption[] = [];
  public select_ifxpts : IMultiSelectOption[] = [];
  public select_ifxcts : IMultiSelectOption[] = [];
  public select_ifxets : IMultiSelectOption[] = [];

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

  constructor(public productService: ProductService, public exportServiceCfg : ExportServiceCfg, public ifxMeasurementService : IfxMeasurementService, builder: FormBuilder) {
    this.builder = builder;
  }

  createStaticForm() {
    this.sampleComponentForm = this.builder.group({
      ID: [this.sampleComponentForm ? this.sampleComponentForm.value.ID : '', Validators.required],
      BaseLines: [this.sampleComponentForm ? this.sampleComponentForm.value.BaseLines : '', Validators.required],
      ProductTag: [this.sampleComponentForm ? this.sampleComponentForm.value.ProductTag : '', Validators.required],
      CommonTags: [this.sampleComponentForm ? this.sampleComponentForm.value.CommonTags : '', null],
      ExtraTags: [this.sampleComponentForm ? this.sampleComponentForm.value.ExtraTags : '', null],
      Measurements: [this.sampleComponentForm ? this.sampleComponentForm.value.Measurements : '', null],
      AlertGroups: [this.sampleComponentForm ? this.sampleComponentForm.value.AlertGroups : '', null],
      Description: [this.sampleComponentForm ? this.sampleComponentForm.value.Description : '']
    });
  }

  reloadData() {
    // now it's a simple subscription to the observable
  this.productService.getProductItem(null)
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
          this.updateAllSelectedItems(this.selectedArray,action.field,action.value);
          break;
       }
       case "AppendProperty": {
         this.updateAllSelectedItems(this.selectedArray,action.field,action.value,true);
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
    for (let i in myArray) {
      console.log("Removing ",myArray[i].ID);
      this.deleteSampleItem(myArray[i].ID,true);
      obsArray.push(this.deleteSampleItem(myArray[i].ID,true));
    }
    this.genericForkJoin(obsArray);
  }

  removeItem(row) {
    let id = row.ID;
    console.log('remove', id);
    this.productService.checkOnDeleteProductItem(id)
      .subscribe(
        data => {
        this.viewModalDelete.parseObject(data)
      },
      err => console.error(err),
      () => { }
      );
  }
  newItem() {
    this.getIfxMeasurementNamesArray();
    //No hidden fields, so create fixed Form
    this.createStaticForm();
    this.editmode = "create";
  }

  editSampleItem(row) {
    this.getIfxMeasurementNamesArray();
    let id = row.ID;
    this.productService.getProductItemById(id)
      .subscribe(data => {
        this.sampleComponentForm = {};
        this.sampleComponentForm.value = data;
        this.oldID = data.ID;
        this.createStaticForm();
        this.editmode = "modify";
      },
      err => console.error(err)
      );
 	}

  deleteSampleItem(id, recursive?) {
    if (!recursive) {
    this.productService.deleteProductItem(id)
      .subscribe(data => { },
      err => console.error(err),
      () => { this.viewModalDelete.hide(); this.reloadData() }
      );
    } else {
      return this.productService.deleteProductItem(id)
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
      this.productService.addProductItem(this.sampleComponentForm.value)
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
        tmpArray = newEntries.concat(component[field]);
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
          r = confirm("Changing Product Instance ID from " + this.oldID + " to " + this.sampleComponentForm.value.ID + ". Proceed?");
        }
        if (r == true) {
          this.productService.editProductItem(this.sampleComponentForm.value, this.oldID)
            .subscribe(data => { console.log(data) },
            err => console.error(err),
            () => { this.editmode = "list"; this.reloadData() }
            );
        }
      }
    } else {
      return this.productService.editProductItem(component, component.ID)
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

  getIfxMeasurementNamesArray() {
    this.ifxMeasurementService.getIfxMeasurementNamesArray(null)
      .subscribe(
      data => {
        this.ifxms_list = data;
        this.select_ifxms = [];
        this.select_ifxms = this.createMultiselectArray(data, 'Name', 'Name', 'ID');
      },
      err => console.error(err),
      () => console.log('DONE')
      );
  }

  pickMeasItem(ifxms_picked) {
    //Only reset values when default values are loaded
    if (ifxms_picked !== this.sampleComponentForm.value.Measurements){
      this.sampleComponentForm.controls.ProductTag.setValue(null);
      this.sampleComponentForm.controls.CommonTags.setValue(null);
      this.sampleComponentForm.controls.ExtraTags.setValue(null);
    }

    if (ifxms_picked) {
      if (ifxms_picked.length > 0) {
        this.ifxMeasurementService.getIfxMeasurementTagsArray(ifxms_picked)
        .subscribe(
          data => {
            this.select_ifxpts = [];
            this.select_ifxcts = [];
            this.select_ifxets = [];
            this.select_ifxpts = this.createMultiselectArray(data);
            this.select_ifxcts = this.createMultiselectArray(data);
            this.select_ifxets = this.createMultiselectArray(data);
            this.sampleComponentForm.controls.ProductTag.setValue(this.cleanSingleselectValue(this.sampleComponentForm.value.ProductTag, data));
            this.sampleComponentForm.controls.CommonTags.setValue(this.cleanMultiselectValue(this.sampleComponentForm.value.CommonTags, data));
            this.sampleComponentForm.controls.ExtraTags.setValue(this.cleanMultiselectValue(this.sampleComponentForm.value.ExtraTags, data));
          },
          err => console.error(err),
          () => console.log('DONE')
        );
      } else {
        //Empty Tag Controls
        this.sampleComponentForm.controls.ProductTag.setValue(null);
        this.sampleComponentForm.controls.CommonTags.setValue(null);
        this.sampleComponentForm.controls.ExtraTags.setValue(null);
        this.select_ifxpts = [];
        this.select_ifxcts = [];
        this.select_ifxets = [];
      }
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

  cleanMultiselectValue(oldValuesArray, currentArray) : string[] {
    let myarray : string[] = [];
    if(oldValuesArray){
      for (let oldValue of oldValuesArray) {
        if (currentArray.indexOf(oldValue) >= 0) {
          myarray.push(oldValue);
        }
      }
    }
    return myarray;
  }

  cleanSingleselectValue(oldValue, currentArray) : string {
    let myValue : string = "";
    if(oldValue){
      if (currentArray.indexOf(oldValue) >= 0) {
        myValue = oldValue;
      }
    }
    return myValue;
  }

}
