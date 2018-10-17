import { Component, Input, Output, EventEmitter, forwardRef, IterableDiffers, SimpleChanges } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Observable } from 'rxjs/Rx';

import { ItemsPerPageOptions } from './global-constants';
import { TableActions } from './table-actions';
import { AvailableTableActions } from './table-available-actions';
import { OnInit, OnChanges } from '@angular/core';

import { ChangeDetectionStrategy } from '@angular/core';
import { ChangeDetectorRef } from '@angular/core';


declare var _: any;

@Component({
  selector: 'table-list',
  styles: [`
		a { outline: none !important; }
	`],
  template: `
    <div class="row">
    <div class="col-md-8 text-left">
      <!--Show/Hide multiselect-->
      <ng-container *ngIf="tableRole !== 'viewonly'">
        <button style ="margin-top: -5px;" type="button" title="{{editEnabled === false ? 'Show multiselect' : 'Hide multiselect' }}" (click)="enableEdit()" class="btn btn-primary"><i class="glyphicon glyphicon-edit"></i></button>
      </ng-container>
      <!--Filtering section-->
      <label [tooltip]="'Clear Filter'" container="body" (click)="onResetFilter()" style="margin-top: 10px"><i class="glyphicon glyphicon-trash text-primary"></i></label>
      <input *ngIf="config.filtering" placeholder="Filter all columns" required = "false" [(ngModel)]="myFilterValue" [ngTableFiltering]="config.filtering" class="form-control select-pages" (tableChanged)="onChangeTable(config)" />
      <span [ngClass]="length > 0 ? ['label label-info'] : ['label label-warning']" style="font-size : 100%">{{length}} Results</span>
      <!--Table Actions-->
      <ng-container *ngIf="typeComponent === 'alerteventhist-component' || typeComponent === 'alertevent-component' || typeComponent === 'kapacitortasks-component'">
        <button style ="margin-top: -5px;" type="button" title="Refresh" (click)="customClick('reloaddata')" class="btn btn-primary"><i class="glyphicon glyphicon-refresh"></i></button>
        <span [ngClass]="['label label-info']" style="font-size : 100%">Last Refresh: {{this.LastUpdate | date : 'HH:mm:ss - Z'}}</span>
      </ng-container>
      <ng-container *ngIf="tableRole === 'fulledit'">
        <button style ="margin-top: -5px;" type="button" (click)="customClick('new')" class="btn btn-primary"><i class="glyphicon glyphicon-plus"></i> New</button>
      </ng-container>
    </div>
    <!--Items per page selection-->
    <div class="col-md-4 text-right">
        <span style="margin-left: 20px"> Items per page: </span>
        <select class="select-pages" style="width:auto" [ngModel]="itemsPerPage || 'All'" (ngModelChange)="changeItemsPerPage($event)">
            <option *ngFor="let option of itemsPerPageOptions" style="padding-left:2px" [value]="option.value">{{option.title}}</option>
        </select>
      </div>
    </div>
    <br>
    <!--Table Actions-->
    <ng-container *ngIf="typeComponent === 'alerteventhist-component' || typeComponent === 'alertevent-component'">
      <div class="row well" style="margin-top: 10px; padding: 10px 0px 10px 15px;">
        <span> Status: </span>
        <label style="font-size:100%" [ngClass]="['label label-success']" (click)="toggleFilter('OK')" container="body" tooltip="Filter OK">{{counterOKs}} OKs <i [ngClass]="OKFilter === true ? ['glyphicon glyphicon-ok'] : ['glyphicon glyphicon-unchecked']"></i></label>
        <label style="font-size:100%;margin-left:5px" [ngClass]="['label label-default']" (click)="toggleFilter('I')" container="body" tooltip="Filter Open">{{counterNOKs}} Open <i [ngClass]="OpenFilter === true ? ['glyphicon glyphicon-ok'] : ['glyphicon glyphicon-unchecked']"></i></label>
        <label style="font-size:100%;margin-left:5px" [ngClass]="['label label-danger']" (click)="toggleFilter('CRITICAL')" container="body" tooltip="Filter CRITICAL">{{counterCrits}} Criticals <i [ngClass]="CritFilter === true ? ['glyphicon glyphicon-ok'] : ['glyphicon glyphicon-unchecked']"></i></label>
        <label style="font-size:100%;margin-left:5px" [ngClass]="['label label-warning']" (click)="toggleFilter('WARNING')" container="body" tooltip="Filter WARNING">{{counterWarns}} Warnings <i [ngClass]="WarnFilter === true ? ['glyphicon glyphicon-ok'] : ['glyphicon glyphicon-unchecked']"></i></label>
        <label style="font-size:100%;margin-left:5px" [ngClass]="['label label-info']" (click)="toggleFilter('INFO')" container="body" tooltip="Filter INFO">{{counterInfos}} Infos <i [ngClass]="InfoFilter === true ? ['glyphicon glyphicon-ok'] : ['glyphicon glyphicon-unchecked']"></i></label>
      </div>
    </ng-container>
    <!--Table available actions-->
    <table-actions [editEnabled]="editEnabled" [counterErrors]="counterErrors" [counterItems]="counterItems || 0" [itemsSelected]="selectedArray.length" [tableAvailableActions]="tableAvailableActions" (actionApply)="customClick('tableaction',$event,selectedArray)"
    [counterTotal]="counterTotal" [counterOKs]="counterOKs" [counterNOKs]="counterNOKs" 
    [counterCrits]="counterCrits" [counterWarns]="counterWarns" [counterInfos]="counterInfos"></table-actions>
    <my-spinner [isRunning]="isRequesting"></my-spinner>
    <!--Table with data -->
    <ng-table *ngIf="isRequesting === false && data"
      [rows]="rows"
      [columns]="columns"
      [sanitizeCell]="sanitizeCell"
      [config]="config"
      [(checkedItems)]="selectedArray"
      [editMode]="editEnabled"
      (tableChanged)="onChangeTable(config)"
      (customClicked)="customClick($event.action, $event.row)"
      [tableRole]="tableRole"
      [roleActions]="roleActions">
    </ng-table>

    <!-- Pagination -->
    <pagination *ngIf="config.paging" class="pagination-sm" [ngModel]="page" [totalItems]="length" [itemsPerPage]="itemsPerPage" [maxSize]="maxSize" [boundaryLinks]="true" [rotate]="false" (pageChanged)="onChangeTable(config, $event)" (numPages)="numPages = $event">
    </pagination>
    <pre *ngIf="config.paging" class="card card-block card-header">Page: {{page}} / {{numPages}}</pre>
    `,
  styleUrls: ['../../css/component-styles.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class TableListComponent implements OnInit, OnChanges {

  //Inputs
  @Input() typeComponent: string;
  @Input() columns: Array<any>;
  @Input() data: Array<any>;
  @Input() counterItems: any = 0;
  @Input() counterErrors: any = [];
  @Input() counterTotal: any = 0;
  @Input() counterOKs: any = 0;
  @Input() counterNOKs: any = 0;
  @Input() counterCrits: any = 0;
  @Input() counterWarns: any = 0;
  @Input() counterInfos: any = 0;
  @Input() selectedArray: any = [];
  @Input() isRequesting: boolean = false;

  @Input() tableRole : any = 'fulledit';
  @Input() roleActions : any = [
    {'name':'export', 'type':'icon', 'icon' : 'glyphicon glyphicon-download-alt text-info', 'tooltip': 'Export item'},
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'edit', 'type':'icon', 'icon' : 'glyphicon glyphicon-edit text-warning', 'tooltip': 'Edit item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'}
  ]

  @Input() sanitizeCell: Function;
  @Output() public customClicked: EventEmitter<any> = new EventEmitter();

  //Vars
  public editEnabled: boolean = false;
  private rows: Array<any> = [];
  public page: number = 1;
  public itemsPerPage: number = 20;
  public itemsPerPageOptions: any = ItemsPerPageOptions;
  public maxSize: number = 5;
  public numPages: number = 1;
  public numPagesShown: number = 1;
  public firstPageShown: number = 1;
  public lastPageShown: number = 1;
  public length: number = 0;
  public tableAvailableActions: any;
  public myFilterValue: any;
  public sortColumn = '';
  public sortDir = '';
  public LastUpdate = new Date();

  public OKFilter: boolean = false;
  public OpenFilter: boolean = false;
  public CritFilter: boolean = false;
  public WarnFilter: boolean = false;
  public InfoFilter: boolean = false;

  //Set config
  public config: any = {
    paging: true,
    sorting: { columns: this.columns },
    filtering: { filterString: '' },
    className: ['table-striped', 'table-bordered']
  };

  ngOnChanges(changes: SimpleChanges) {
    if (!this.data) this.data = [];
    this.onChangeTable(this.config);
    this.cd.markForCheck();
  }

  ngOnInit() {
    this.onResetFilterColumns();
    this.config.sorting = { columns: this.columns };
    this.onChangeTable(this.config)
  }

  constructor(public cd: ChangeDetectorRef) { }

  //Enable edit tables
  enableEdit() {
    this.editEnabled = !this.editEnabled;
    let obsArray = [];
    this.tableAvailableActions = new AvailableTableActions(this.typeComponent).availableOptions;
  }

  public changePage(page: any, data: Array<any> = this.data): Array<any> {
    //Check if we have to change the actual page

    let maxPage = Math.ceil(data.length / this.itemsPerPage);
    if (page.page > maxPage && page.page != this.page) this.page = page.page = maxPage;
    if (page.page != this.page) this.page = page.page;
    let start = (page.page - 1) * page.itemsPerPage;
    let end = page.itemsPerPage > -1 ? (start + page.itemsPerPage) : data.length;
    return data.slice(start, end);
  }

  public changeSort(data: any, config: any): any {
    if (!config.sorting) {
      return data;
    }

    let columns = this.config.sorting.columns || [];
    let columnName: string = void 0;
    let sort: string = void 0;

    for (let i = 0; i < columns.length; i++) {
      if (typeof columns[i].sort !== 'undefined' && columns[i].sort !== '' && columns[i].sort !== false) {
        columnName = columns[i].name;
        sort = columns[i].sort;
        this.sortColumn = columnName;
        this.sortDir = sort;
      }
    }

    if (!columnName) {
      return data;
    }

    // simple sorting
    return data.sort((previous: any, current: any) => {
      if (previous[columnName] > current[columnName]) {
        return sort === 'desc' ? -1 : 1;
      } else if (previous[columnName] < current[columnName]) {
        return sort === 'asc' ? -1 : 1;
      }
      return 0;
    });
  }

  public changeFilter(data: any, config: any): any {
    let filteredData: Array<any> = data;
    if (!config.filtering) {
      return filteredData;
    }

    if (config.filtering.columnName && config.filtering.columnName.length > 0) {
      filteredData = filteredData.filter((item: any) => 
        (item[config.filtering.columnName] == null ? '' : item[config.filtering.columnName]).toString().match(this.config.filtering.filterString)
      );
      return filteredData;
    }

    this.columns.forEach((column: any) => {
      if (column.filtering) {
        filteredData = filteredData.filter((item: any) => {
          return (item[column.name] == null ? '' : item[column.name]).toString().match(column.filtering.filterString);
        });
      }
    });

    return this.filterData(filteredData, this.config.filtering.filterString);
  }

  filterData(srcArray: Array<any>, filterString: string): Array<any> {
    let tempArray: Array<any> = [];
    srcArray.forEach((item: any) => {
      let flag = false;
      this.columns.forEach((column: any) => {
        if ((item[column.name] == null ? '' : item[column.name]).toString().match(filterString)) {
          flag = true;
        }
      });
      if (flag) {
        tempArray.push(item);
      }
    });
    return tempArray;
  }

  changeItemsPerPage(items) {
    this.itemsPerPage = parseInt(items);
    let maxPage = Math.ceil(this.length / this.itemsPerPage);
    if (this.page > maxPage) this.page = maxPage;
    this.onChangeTable(this.config);
  }

  public onChangeTable(config: any, page: any = { page: this.page, itemsPerPage: this.itemsPerPage }): any {
    if (config) {
      if (config.filtering) {
        Object.assign(this.config.filtering, config.filtering);
      }
      if (config.sorting) {
        Object.assign(this.config.sorting, config.sorting);
      }
    }
    let filteredData = this.changeFilter(this.data, this.config);
    let sortedData = this.changeSort(filteredData, this.config);
    this.rows = page && config.paging ? this.changePage(page, sortedData) : sortedData;
    this.length = sortedData.length;
    let maxPage = Math.ceil(this.length / this.itemsPerPage);
    this.numPagesShown = maxPage < this.maxSize ? maxPage : this.maxSize;
    this.firstPageShown = ((Math.ceil(this.page / this.maxSize) - 1) * this.numPagesShown) + 1;
    this.lastPageShown = (this.firstPageShown + this.numPagesShown - 1) < maxPage ? (this.firstPageShown + this.numPagesShown - 1) : maxPage;
    if (this.firstPageShown + this.numPagesShown > maxPage) this.numPagesShown = maxPage - this.firstPageShown + 1;
  }

  onResetFilter(): void {
    this.page = 1;
    this.myFilterValue = "";
    this.config.filtering = { filtering: { filterString: '' } };
    this.onResetFilterColumns();
    this.onChangeTable(this.config);
  }

  onResetFilterColumns(): void {
    this.columns.forEach((column: any) => {
      if (column.filtering) {
        column.filtering.filterString = '';
      }
    });
  }

  onFilterColumn(columnName: string, filterString: string): void {
    this.columns.forEach((column: any) => {
      if (column.filtering && column.name == columnName) {
        column.filtering.filterString = filterString;
      }
    });
    this.onChangeTable(this.config);
  }

  toggleFilter(filterString: string) {
    if (this.OKFilter === false && filterString === 'OK') {
      this.OKFilter = true;
      this.OpenFilter = false;
      this.CritFilter = false;
      this.WarnFilter = false;
      this.InfoFilter = false;
    } else if (this.OpenFilter === false && filterString === 'I') {
      this.OKFilter = false;
      this.OpenFilter = true;
      this.CritFilter = false;
      this.WarnFilter = false;
      this.InfoFilter = false;
    } else if (this.CritFilter === false && filterString === 'CRITICAL') {
      this.OKFilter = false;
      this.OpenFilter = false;
      this.CritFilter = true;
      this.WarnFilter = false;
      this.InfoFilter = false;
    } else if (this.WarnFilter === false && filterString === 'WARNING') {
      this.OKFilter = false;
      this.OpenFilter = false;
      this.CritFilter = false;
      this.WarnFilter = true;
      this.InfoFilter = false;
    } else if (this.InfoFilter === false && filterString === 'INFO') {
      this.OKFilter = false;
      this.OpenFilter = false;
      this.CritFilter = false;
      this.WarnFilter = false;
      this.InfoFilter = true;
    } else {
      this.OKFilter = false;
      this.OpenFilter = false;
      this.CritFilter = false;
      this.WarnFilter = false;
      this.InfoFilter = false;
      filterString = '';
    }
    this.onFilterColumn('Level', filterString);
  }

  customClick(clicked: string, event: any = "", data: any = ""): void {
    if (clicked == "reloaddata") {
      //console.log("customClick with reloaddata");
      this.LastUpdate = new Date();
    }
    this.customClicked.emit({ 'option': clicked, 'event': event, 'data': data, 'sortColumn': this.sortColumn, 'sortDir': this.sortDir });
    //pending change for future, to get only the list of results to show, not all the list
    //this.customClicked.emit({ 'option': clicked, 'event': event, 'data': data, 'sortColumn': this.sortColumn, 'sortDir': this.sortDir, 'page': this.page, 'itemsPerPage': this.itemsPerPage, 'maxSize': this.maxSize });
  }

}
