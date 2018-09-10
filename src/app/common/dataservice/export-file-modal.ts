import { Component, Input, Output, Pipe, PipeTransform, ViewChild, EventEmitter } from '@angular/core';
import { ModalDirective } from 'ngx-bootstrap';
import { Validators, FormGroup, FormControl, FormArray, FormBuilder } from '@angular/forms';
import { ExportServiceCfg } from './export.service'
import { TreeView} from './treeview';
import { IMultiSelectOption, IMultiSelectSettings, IMultiSelectTexts } from '../multiselect-dropdown';

//Services
import { IfxServerService } from '../../ifxserver/ifxserver.service';
import { AlertService } from '../../alert/alert.service';
import { DeviceStatService } from '../../devicestat/devicestat.service';
import { KapacitorService } from '../../kapacitor/kapacitor.service';
import { EndpointService } from '../../endpoint/endpoint.service';
import { ProductService } from '../../product/product.service';
import { ProductGroupService } from '../../productgroup/productgroup.service';
import { RangeTimeService } from '../../rangetime/rangetime.service';
import { TemplateService } from '../../template/template.service';

import { Subscription } from 'rxjs';

@Component({
  selector: 'export-file-modal',
  template: `
      <div bsModal #childModal="bs-modal" class="modal fade" tabindex="-1" role="dialog" aria-labelledby="myLargeModalLabel" aria-hidden="true">
          <div class="modal-dialog" style="width:90%">
            <div class="modal-content" >
              <div class="modal-header">
                <button type="button" class="close" (click)="childModal.hide()" aria-label="Close">
                  <span aria-hidden="true">&times;</span>
                </button>
                <h4 class="modal-title" *ngIf="exportObject != null">{{titleName}} <b>{{ exportObject.ID }}</b> - <label [ngClass]="['label label-'+colorsObject[exportType]]">{{exportType}}</label></h4>
                <h4 class="modal-title" *ngIf="exportObject == null">Export</h4>
              </div>
              <div class="modal-body">

              <div *ngIf="prepareExport === false">
              <div class="row">
              <div class="col-md-2">
              <div class="panel-heading">
              1.Select type:
              </div>
                <div class="panel panel-default" *ngFor="let items of objectTypes; let i = index" style="margin-bottom: 0px" >
                  <div class="panel-heading" (click)="loadSelection(i, items.Type)" role="button">
                  <i [ngClass]="selectedType ? (selectedType.Type === items.Type ? ['glyphicon glyphicon-eye-open'] : ['glyphicon glyphicon-eye-close'] ) : ['glyphicon glyphicon-eye-close']"  style="padding-right: 10px"></i>
                  <label [ngClass]="['label label-'+items.Class]"> {{items.Type}}  </label>
                  </div>
                </div>
                </div>
                <div class="col-md-5">
                <div *ngIf="selectedType">
                  <div class="panel-heading">
                    <div>
                      2. Select Items of type <label [ngClass]="['label label-'+selectedType.Class]"> {{selectedType.Type}}</label>  <span class="badge" style="margin-left: 10px">{{resultArray.length}} Results</span>
                    </div>
                    <div dropdown class="text-left" style="margin-top: 10px">
                    <span class="dropdown-toggle-split">Filter by</span>
                    <ss-multiselect-dropdown style="border: none" [options]="listFilterProp" [texts]="myTexts" [settings]="propSettings" [(ngModel)]="selectedFilterProp" (ngModelChange)="onChange(filter)"></ss-multiselect-dropdown>
                      <input type=text [(ngModel)]="filter" placeholder="Filter items..." (ngModelChange)="onChange($event)">
                      <label [tooltip]="'Clear Filter'" container="body" (click)="filter=''; onChange(filter)"><i class="glyphicon glyphicon-trash text-primary"></i></label>
                    </div>
                    <div class="text-right">
                    <label class="label label-success" (click)=selectAllItems(true)>Select All</label>
                    <label class="label label-danger" (click)=selectAllItems(false)>Deselect All</label>
                    </div>
                  </div>
                  <div style="max-height: 400px; overflow-y:auto">
                    <div *ngFor="let res of resultArray;  let i = index" style="margin-bottom: 0px" >
                      <treeview [keyProperty]="selectedFilterProp" [showType]="false" [visible]="false" [title]="res.ID" [object]="res" [alreadySelected]="checkItems(res.ID, selectedType.Type)" [type]="selectedType.Type" [visibleToogleEnable]="true" [addClickEnable]="true" (addClicked)="selectItem($event,index)"  style="margin-bottom:0px !important">{{res}}</treeview>
                    </div>
                  </div>
                  </div>
                  </div>
                  <div class="col-md-5">
                  <div *ngIf="finalArray.length !== 0">
                    <div class="panel-heading"> 3. Items ready to export: <span class="badge">{{finalArray.length}}</span>
                    <div class="text-right">
                      <label class="label label-danger" (click)="finalArray = []">Clear All</label>
                    </div>
                    </div>
                    <div style="max-height: 400px; overflow-y:auto">
                      <div *ngFor="let res of finalArray;  let i = index" class="col-md-12">
                      <i class="text-danger glyphicon glyphicon-remove-sign col-md-1" role="button" style="margin-top: 15px;" (click)="removeItem(i)"> </i>
                        <treeview [visible]="false" [title]="res.ObjectID" [object]="res" [type]="res.ObjectTypeID" [recursiveToogle]="true" [index] = "i" (recursiveClicked)="toogleRecursive($event)" class="col-md-11">{{res}}</treeview>
                      </div>
                    </div>
                    </div>
                    </div>
                </div>
                </div>
              <div *ngIf="prepareExport === true">
              <div *ngIf="exportResult === true" style="overflow-y: scroll; max-height: 350px">
                <h4 class="text-success"> <i class="glyphicon glyphicon-ok-circle" style="padding-right:10px"></i>Succesfully exported {{exportedItem.Objects.length}} items </h4>
                <div *ngFor="let object of exportedItem.Objects; let i=index">
                  <treeview [visible]="false" [visibleToogleEnable]="true" [title]="object.ObjectID" [type]="object.ObjectTypeID" [object]="object.ObjectCfg"> </treeview>
              </div>
              </div>
              <div  *ngIf="exportForm && exportResult === false">
              <form class="form-horizontal" *ngIf="bulkExport === false">
              <div class="form-group">
                <label class="col-sm-2 col-offset-sm-2" for="Recursive">Recursive</label>
                <i placement="top" style="float: left" class="info control-label glyphicon glyphicon-info-sign" tooltipAnimation="true" tooltip="Select if the export request will include all related componentes"></i>
                <div class="col-sm-9">
                <select name="recursiveObject" class="form-control" id="recursiveObject" [(ngModel)]="recursiveObject">
                  <option value="true">True</option>
                  <option value="false">False</option>
                </select>
                </div>
              </div>
              </form>
              <form [formGroup]="exportForm" class="form-horizontal"  >
                    <div class="form-group">
                      <label for="FileName" class="col-sm-2 col-offset-sm-2 control-FileName">FileName</label>
                      <i placement="top" style="float: left" class="info control-label glyphicon glyphicon-info-sign" tooltipAnimation="true" tooltip="Desired file name"></i>
                      <div class="col-sm-9">
                      <input type="text" class="form-control" placeholder="file.json" formControlName="FileName" id="FileName">
                      </div>
                    </div>
                    <div class="form-group">
                      <label for="Author" class="col-sm-2 control-Author">Author</label>
                      <i placement="top" style="float: left" class="info control-label glyphicon glyphicon-info-sign" tooltipAnimation="true" tooltip="Author of the export"></i>
                      <div class="col-sm-9">
                      <input type="text" class="form-control" placeholder="pseriescollector" formControlName="Author" id="Author">
                      </div>
                    </div>
                    <div class="form-group">
                      <label for="Tags" class="col-sm-2 control-Tags">Tags</label>
                      <i placement="top" style="float: left" class="info control-label glyphicon glyphicon-info-sign" tooltipAnimation="true" tooltip="Related tags to identify exported data"></i>
                      <div class="col-sm-9">
                      <input type="text" class="form-control" placeholder="cisco,catalyst,..." formControlName="Tags" id="Tags">
                      </div>
                    </div>

                    <div class="form-group">
                      <label for="FileName" class="col-sm-2 control-FileName">Description</label>
                      <i placement="top" style="float: left" class="info control-label glyphicon glyphicon-info-sign" tooltipAnimation="true" tooltip="Description of the exported file"></i>
                      <div class="col-sm-9">
                      <textarea class="form-control" style="width: 100%" rows="2" formControlName="Description" id="Description"> </textarea>
                      </div>
                    </div>
              </form>
              </div>
              </div>
              </div>
              <div class="modal-footer" *ngIf="showValidation === true">
               <button type="button" class="btn btn-default" (click)="childModal.hide()">Close</button>
               <button *ngIf="exportResult === false && prepareExport === true" type="button" class="btn btn-primary" (click)="exportBulkItem()">{{textValidation ? textValidation : Save}}</button>
               <button *ngIf="prepareExport === false" type="button" class="btn btn-primary" [disabled]="finalArray.length === 0" (click)="showExportForm()">Continue</button>
             </div>
            </div>
          </div>
        </div>`,
        styleUrls: ['./import-modal-styles.css'],
        providers: [IfxServerService,AlertService,DeviceStatService,KapacitorService ,EndpointService ,ProductService,ProductGroupService,RangeTimeService ,TemplateService, TreeView]
})

export class ExportFileModal {
  @ViewChild('childModal') public childModal: ModalDirective;
  @Input() titleName: any;
  @Input() customMessage: string;
  @Input() showValidation: boolean;
  @Input() textValidation: string;
  @Input() prepareExport : boolean = true;
  @Input() bulkExport: boolean = false;
  @Input() exportType: any = null;

  @Output() public validationClicked: EventEmitter<any> = new EventEmitter();

  public validationClick(myId: string): void {
    this.validationClicked.emit(myId);
  }

  public builder: any;
  public exportForm: any;
  public mySubscriber: Subscription;

  constructor(builder: FormBuilder, public exportServiceCfg : ExportServiceCfg,
    public alertService : AlertService, public ifxServerService : IfxServerService, 
    public deviceStatService : DeviceStatService, public kapacitorService : KapacitorService,
    public endpointService : EndpointService,public productService : ProductService, 
    public productGroupService : ProductGroupService, public rangeTimeService : RangeTimeService,
    public templateService : TemplateService) {

    this.builder = builder;
  }

//COMMON
  createStaticForm() {
    this.exportForm = this.builder.group({
      FileName: [this.prepareExport ? this.exportObject.ID+'_'+this.exportType+'_'+this.nowDate+'.json' : 'bulkexport_'+this.nowDate+'.json' , Validators.required],
      Author: ['resistor', Validators.required],
      Tags: [''],
      Recursive: [true, Validators.required],
      Description: ['Autogenerated', Validators.required]
    });
  }

  //Single Object Export:
  exportObject: any = null;

   //Single Object
  public colorsObject : Object = {
    "alertcfg": 'danger',
    "devicestatcfg": 'info',
    "ifxservercfg": 'success',
    "kapacitorcfg": 'primary',
    "endpointcfg": 'default',
    "productcfg": 'warning',
    "productgroupcfg": 'info',
    "rangetimecfg": 'danger',
    "templatecfg": 'primary'
   };

  //Control to load exported result
  exportResult : boolean = false;
  exportedItem : any;

  //Others
  nowDate : any;
  recursiveObject : boolean = true;

  //Bulk Export - Result Array from Loading data:

  resultArray : any = [];
  dataArray : any = [];
  selectedFilterProp : any = "ID";
  listFilterProp: IMultiSelectOption[] = [];
  private propSettings: IMultiSelectSettings = {
      singleSelect: true,
  };

  //Bulk Export - SelectedType
  selectedType : any = null;
  finalArray : any = [];
  filter : any = "";
  //Bulk Objects
  public objectTypes : any = [
    {'Type':"alertcfg", 'Class': 'danger', 'Visible':false},
    {'Type':"devicestatcfg", 'Class': 'info', 'Visible':false},
    {'Type':"ifxservercfg", 'Class' : 'success', 'Visible':false},
    {'Type':"kapacitorcfg", 'Class' : 'primary', 'Visible':false},
    {'Type':"endpointcfg", 'Class' : 'default', 'Visible':false},
    {'Type':"productcfg", 'Class' : 'warning', 'Visible':false},
    {'Type':"productgroupcfg", 'Class' : 'info', 'Visible':false},
    {'Type':"rangetimecfg", 'Class' : 'danger', 'Visible':false},
    {'Type':"templatecfg", 'Class' : 'primary', 'Visible':false}
   ]

   //Reset Vars on Init
  clearVars() {
    this.finalArray = [];
    this.resultArray = [];
    this.selectedType = null;
    this.exportResult = false;
    this.exportedItem = [];
    this.exportType = this.exportType || null;
    this.exportObject = null;
  }

  //Init Modal, depending from where is called
  initExportModal(exportObject: any, prepareExport? : boolean) {
    this.clearVars();
    if (prepareExport === false) {
      this.prepareExport = false;
    } else {
      this.prepareExport = true;
    };
    //Single export
    if (this.prepareExport === true) {
      console.log(exportObject);
      this.exportObject = exportObject.row || exportObject;
      this.exportType = exportObject.exportType || this.exportType;
      //Sets the FinalArray to export the items, in this case only be 1
      this.finalArray = [{
        //ensure string on object ID
        'ObjectID' : this.exportObject.ID.toString(),
        'ObjectTypeID' :  this.exportType,
        'Options' : {
          Recursive: this.recursiveObject
        }
      }]
    //Bulk export
  } else {
    this.exportObject = exportObject;
    this.exportType = null;
  }
   
    this.nowDate = this.getCustomTimeString()
    this.createStaticForm();
    this.childModal.show();
  }

  getCustomTimeString(){
    let date  = new Date();
    let day = ('0' + date.getDate()).slice(-2);
    let year  = date.getFullYear().toString();
    let month  = ('0' + (date.getMonth()+1)).slice(-2);
    let ymd =year+month+day;
    let hm  =  date.getHours().toString()+date.getMinutes().toString();
    return ymd + '_' + hm;
  }

  onChange(event){
    let tmpArray = this.dataArray.filter((item: any) => {
      if (item[this.selectedFilterProp]) return item[this.selectedFilterProp].toString().match(event);
      else if (event === "" && !item[this.selectedFilterProp]) return item;
    });
    this.resultArray = tmpArray;
  }
  changeFilterProp(prop){
    this.selectedFilterProp = prop;
  }

  //Load items from selection type
   loadSelection(i, type) {
     for (let a of this.objectTypes) {
       if(type !== this.objectTypes[i].Type) {
         this.objectTypes[i].Visible = false;
       }
     }
     this.objectTypes[i].Visible = true;
     this.selectedType = this.objectTypes[i];
     this.filter = "";
     this.loadItems(type,null);
   }

   checkItems(checkItem: any,type) : boolean {
     //Extract the ID from finalArray and loaded Items:
     let exist = true;
     for (let a of this.finalArray) {
       if (checkItem == a.ObjectID) {
         exist = false;
       }
     }
     return exist;
   }

   //Common function to find given object property inside an array
   findIndexItem(checkArray, checkItem: any) : any {
     for (let a in checkArray) {
       if (checkItem == checkArray[a].ObjectID) {
         return a;
       }
     }
   }

   selectAllItems(selectAll) {
     //Creates the form array
     if (selectAll === true) {
       for (let a of this.resultArray) {
         if (this.checkItems(a.ID, this.selectedType)) {
           //ensure string on object ID
           this.finalArray.push({ "ObjectID" : a.ID.toString(), ObjectTypeID: this.selectedType.Type, "Options" : {'Recursive': false }});
         }
       }
     } else {
       for (let a of this.resultArray) {
         let index = this.findIndexItem(this.finalArray, a.ID);
           if (index) this.removeItem(index);
       }
     }
   }

   //Select item to add it to the FinalArray or delete it if its alreay selected
  selectItem(event) {
    if (this.checkItems(event.ObjectID, event.ObjectTypeID)) {
      //ensure string on object ID
      event.ObjectID =  event.ObjectID.toString()
      this.finalArray.push(event);
    }
    else {
      let index = this.findIndexItem(this.finalArray, event.ObjectID);
      this.removeItem(index);
    }
  }
  //Remove item from Array
  removeItem(index) {
    this.finalArray.splice(index,1);
  }

  //Change Recursive option on the FinalArray objects
  toogleRecursive(event) {
    this.finalArray[event.Index].Options.Recursive = event.Recursive;
  }

  showExportForm() {
    this.prepareExport = true;
  }

  exportBulkItem(){
    if (this.bulkExport === false) this.finalArray[0].Options.Recursive = this.recursiveObject;

    let finalValues = {"Info": this.exportForm.value, "Objects" : this.finalArray}
    this.exportServiceCfg.bulkExport(finalValues)
    .subscribe(
      data => {
        this.exportedItem = data[1];
        saveAs(data[0],data[1].Info.FileName);
        this.exportResult = true;
      },
      err => console.error(err),
      () => console.log("DONE"),
    );
  }

//Load items functions from services depending on items selected Type
  loadItems(type, filter?) {
    this.resultArray = [];
    this.selectedFilterProp = "ID";
    this.listFilterProp = [];
    if (this.mySubscriber) this.mySubscriber.unsubscribe();
    switch (type) {
      case 'ifxservercfg':
       this.mySubscriber = this.ifxServerService.getIfxServerItem(filter)
       .subscribe(
       data => {
         //Load items on selection
         this.dataArray = data;
         this.resultArray = this.dataArray;
         for (let i in this.dataArray[0]) {
           this.listFilterProp.push({ 'id': i, 'name': i });
         }
       },
       err => {console.log(err)},
       () => {console.log("DONE")}
       );
      break;
      case 'alertcfg':
      this.mySubscriber = this.alertService.getAlertItem(filter)
       .subscribe(
       data => {
         this.dataArray=data;
         this.resultArray = this.dataArray;
         for (let i in this.dataArray[0]) {
           this.listFilterProp.push({ 'id': i, 'name': i });
         }
       },
       err => {console.log(err)},
       () => {console.log("DONE")}
       );
      break;
      case 'devicestatcfg':
      this.mySubscriber = this.deviceStatService.getDeviceStatItem(filter)
       .subscribe(
       data => {
         this.dataArray=data;
         this.resultArray = this.dataArray;
         for (let i in this.dataArray[0]) {
           this.listFilterProp.push({ 'id': i, 'name': i });
         }
       },
       err => {console.log(err)},
       () => {console.log("DONE")}
       );
      break;
      case 'kapacitorcfg':
      this.mySubscriber = this.kapacitorService.getKapacitorItem(filter)
       .subscribe(
       data => {
         this.dataArray=data;
         this.resultArray = this.dataArray;
         for (let i in this.dataArray[0]) {
           this.listFilterProp.push({ 'id': i, 'name': i });
         }
       },
       err => {console.log(err)},
       () => {console.log("DONE")}
       );
      break;
      case 'endpointcfg':
      this.mySubscriber = this.endpointService.getEndpointItem(filter)
       .subscribe(
       data => {
         this.dataArray=data;
         this.resultArray = this.dataArray;
         for (let i in this.dataArray[0]) {
           this.listFilterProp.push({ 'id': i, 'name': i });
         }
       },
       err => {console.log(err)},
       () => {console.log("DONE")}
       );
      break;
      case 'productcfg':
      this.mySubscriber = this.productService.getProductItem(filter)
       .subscribe(
       data => {
         this.dataArray=data;
         this.resultArray = this.dataArray;
         for (let i in this.dataArray[0]) {
           this.listFilterProp.push({ 'id': i, 'name': i });
         }
       },
       err => {console.log(err)},
       () => {console.log("DONE")}
       );
      break;
      case 'productgroupcfg':
      this.mySubscriber = this.productGroupService.getProductGroupItem(filter)
       .subscribe(
       data => {
         this.dataArray=data;
         this.resultArray = this.dataArray;
         for (let i in this.dataArray[0]) {
           this.listFilterProp.push({ 'id': i, 'name': i });
         }
       },
       err => {console.log(err)},
       () => {console.log("DONE")}
       );
      break;
      case 'rangetimecfg':
      this.mySubscriber = this.rangeTimeService.getRangeTimeItem(filter)
       .subscribe(
       data => {
         this.dataArray=data;
         this.resultArray = this.dataArray;
         for (let i in this.dataArray[0]) {
           this.listFilterProp.push({ 'id': i, 'name': i });
         }
       },
       err => {console.log(err)},
       () => {console.log("DONE")}
       );
      break;
      case 'templatecfg':
      this.mySubscriber = this.templateService.getTemplateItem(filter)
       .subscribe(
       data => {
         this.dataArray=data;
         this.resultArray = this.dataArray;
         for (let i in this.dataArray[0]) {
           this.listFilterProp.push({ 'id': i, 'name': i });
         }
       },
       err => {console.log(err)},
       () => {console.log("DONE")}
       );
      break;

      default:
      break;
    }
  }
  //Common Functions
  isArray(myObject) {
    return myObject instanceof Array;
  }

  isObject(myObject) {
    return typeof myObject === 'object'
  }

  hide() {
    if (this.mySubscriber) this.mySubscriber.unsubscribe();
    this.childModal.hide();
    
  }

}
