import { Component, EventEmitter, Input, Output, Pipe, PipeTransform,  } from '@angular/core';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { ElapsedSecondsPipe } from '../../../custom-pipe/elapsedseconds.pipe';

@Component({
  selector: 'ng-table',
  templateUrl: './ng-table.html',
  styles: [`
    :host >>> .displayinline{
      display: inline !important;
      padding-left: 5px;
    }
`]
})
export class NgTableComponent {
  // Table values
  @Input() public rows: Array<any> = [];
  @Input() public showRoleActions: boolean = true;

  @Input() public roleActions: any;
  @Input() public tableRole : string = 'fulledit';

  @Input() public showStatus: boolean = false;
  @Input() public editMode: boolean = false;
  @Input() public exportType: string;
  @Input() public extraActions: Array<any>;
  @Input() checkedItems: Array<any>;
  @Input() checkRows: Array<any>;
  @Input() sanitizeCell: Function;

  @Input()
  public set config(conf: any) {
    if (!conf.className) {
      conf.className = 'table-striped table-bordered';
    }
    if (conf.className instanceof Array) {
      conf.className = conf.className.join(' ');
    }
    this._config = conf;
  }

  // Outputs (Events)
  @Output() public tableChanged: EventEmitter<any> = new EventEmitter();
  @Output() public cellClicked: EventEmitter<any> = new EventEmitter();
  @Output() public customClicked: EventEmitter<any> = new EventEmitter();
  @Output() public testedConnection: EventEmitter<any> = new EventEmitter();
  @Output() public extraActionClicked: EventEmitter<any> = new EventEmitter();

  public showFilterRow: Boolean = false;

  @Input()
  public set columns(values: Array<any>) {
    this._columns = [];
    values.forEach((value: any) => {
      if (value.filtering) {
        this.showFilterRow = true;
      }
      if (value.className && value.className instanceof Array) {
        value.className = value.className.join(' ');
      }
      let column = this._columns.find((col: any) => col.name === value.name);
      if (column) {
        Object.assign(column, value);
      }
      if (!column) {
        this._columns.push(value);
      }
    });
  }

  private _columns: Array<any> = [];
  private _config: any = {};

  public constructor(private sanitizer: DomSanitizer) {
  }

  public sanitize(html: any, transform?: any ): SafeHtml {

    let output: string
    if (typeof this.sanitizeCell === "function" ) {
      output = this.sanitizeCell(html,transform)
      if ( output.length > 0 ) {
        return output
      }
    }
    if  (transform === "ns2s") {
      html = html / 1.e9;
      transform = "elapsedseconds";
    }
    if  (transform === "elapsedseconds") {
      let test = new ElapsedSecondsPipe().transform(html,'3');
      html = test.toString();
    }
    if  (transform === "imgwtooltip") {
      if (html) return '<i class="glyphicon glyphicon-remove text-danger" title="'+html+'"></i>';
      else return html;
    }
    if  (transform === "color") {
      let color: string = "green";
      if (html === "INFO") {
        color = "blue";
      } else if (html === "WARNING") {
        color = "orange";
      } else if (html === "CRITICAL") {
        color = "red";
      } 
      html = "<span style='color:"+color+";font-weight:bold'>"+html+"</span>";
    }
    if (typeof html === 'object') {
      var test: any = '<ul class="list-unstyled">';
      for (var item of html) {
        if (typeof item === 'object') {
          test += "<li>";
          for (var item2 in Object(item)) {
            if (typeof item[item2] === 'boolean') {
              if (item[item2]) test += ' <i class="glyphicon glyphicon-arrow-right"></i>'
              else test += ' <i class="glyphicon glyphicon-alert"></i>'
            } else if (item2 === 'TagID') {
              test += '<h4 class="text-success displayinline">'+item[item2] +' - </h4>';
            } else {
              test +='<span>'+item[item2]+'</span>';
            }
          }
          test += "</li>";
        } else {
          test += "<li>" + item + "</li>";
        }
      }
      test += "</ul>"
      return test;
    } else if (typeof html === 'boolean') {
      if (html) return '<i class="glyphicon glyphicon-ok text-success"></i>'
      else return '<i class="glyphicon glyphicon-remove text-danger"></i>'
    }
    else {
      return this.sanitizer.bypassSecurityTrustHtml(html);
    }
  }

  public get columns(): Array<any> {
    return this._columns;
  }

  public get config(): any {
    return this._config;
  }

  public get configColumns(): any {
    let sortColumns: Array<any> = [];

    this.columns.forEach((column: any) => {
      if (column.sort) {
        sortColumns.push(column);
      }
    });
    return { columns: sortColumns };
  }

  public onChangeTable(column: any): void {
    this._columns.forEach((col: any) => {
      if (col.name !== column.name && col.sort !== false) {
        col.sort = '';
      }
    });
    this.tableChanged.emit({ sorting: this.configColumns });
  }

  selectAllItems(selectAll) {
    //Creates the form array
    if (selectAll === true) {
      for (let row of this.rows) {
        if (this.checkItems(row.ID)) {
          this.checkedItems.push(row);
        }
      }
    } else {
      for (let row of this.rows) {
        let index = this.findIndexItem(row.ID);
          if (index) this.deleteItem(index);
      }
    }
  }

  checkItems(item: any) : boolean {
    //Extract the ID from finalArray and loaded Items:
    let exist = true;
    for (let a of this.checkedItems) {
      if (item === a.ID) {
        exist = false;
      }
    }
    return exist;
  }

  selectItem(row : any) : void {
    if (this.checkItems(row.ID)) {
      this.checkedItems.push(row);
    }
    else {
      let index = this.findIndexItem(row.ID);
      this.deleteItem(index);
    }
  }
  //Remove item from Array
  deleteItem(index) {
    this.checkedItems.splice(index,1);
  }

  findIndexItem(ID) : any {
    for (let a in this.checkedItems) {
      if (ID === this.checkedItems[a].ID) {
        return a;
      }
    }
  }

  public cellClick(row: any, column: any): void {
    this.cellClicked.emit({ row, column });
  }

  public customClick(action: string, row: any) : void {
    this.customClicked.emit({'action' : action, 'row' : row});
  }

  public testConnection(row: any) : void {
    this.testedConnection.emit(row);
  }
  public extraActionClick(row: any, action: any, property? : any) : void {
    this.extraActionClicked.emit({row , action, property});
  }
}
