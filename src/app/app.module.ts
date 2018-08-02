//Auth examples from: https://github.com/auth0-blog/angular2-authentication-sample
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { RouterModule } from '@angular/router';
import { HttpModule } from '@angular/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { BsDropdownModule } from 'ngx-bootstrap';

// external libs

import { Ng2TableModule } from './common/ng-table/ng2-table';
/*import { TestNg2TableModule } from './common/ng-table-test/ng2-table';*/

import { HomeComponent } from './home/home.component';
import { NavbarComponent } from './home/navbar/navbar.component';
import { SideMenuComponent } from './home/sidemenu/sidemenu.component';

import { LoginComponent } from './login/login.component';
import { App } from './app';

import { AppRoutes } from './app.routes';
//common
import { ControlMessagesComponent } from './common/custom-validation/control-messages.component';
import { MultiselectDropdownModule } from './common/multiselect-dropdown'
import { PasswordToggleDirective } from './common/custom-directives'
import { TableActions } from './common/table-actions';
//Resistor Components

import { BlockUIService } from './common/blockui/blockui-service';

import { AccordionModule , PaginationModule ,TabsModule } from 'ngx-bootstrap';
import { TooltipModule } from 'ngx-bootstrap';
import { ModalModule } from 'ngx-bootstrap';
import { ModalDirective } from 'ngx-bootstrap';
import { ProgressbarModule } from 'ngx-bootstrap';
import { TimepickerModule } from 'ngx-bootstrap';

import { GenericModal } from './common/custom-modal/generic-modal';
import { AboutModal } from './home/about-modal';
import { TreeView } from './common/dataservice/treeview';
import { ImportFileModal } from './common/dataservice/import-file-modal'
import { ExportFileModal } from './common/dataservice/export-file-modal'
import { CoreModule } from './core/core.module';

//others
import { WindowRef } from './common/windowref';
import { ValidationService } from './common/custom-validation/validation.service';
import { ExportServiceCfg } from './common/dataservice/export.service'

import { CustomPipesModule } from './common/custom-pipe/custom-pipe.module';

import { BlockUIComponent } from './common/blockui/blockui-component';
import { SpinnerComponent } from './common/spinner';

import { TableListComponent } from './common/table-list.component';

import { KapacitorComponent } from './kapacitor/kapacitor.component';
import { KapacitorTasksComponent } from './kapacitortasks/kapacitortasks.component';
import { RangeTimeComponent } from './rangetime/rangetime.component';
import { IfxServerComponent } from './ifxserver/ifxserver.component';
import { IfxDBComponent } from './ifxdb/ifxdb.component';
import { IfxMeasurementComponent } from './ifxmeasurement/ifxmeasurement.component';
import { ProductComponent } from './product/product.component';
import { ProductGroupComponent } from './productgroup/productgroup.component';
import { TemplateComponent } from './template/template.component';
import { OutHTTPComponent } from './outhttp/outhttp.component';
import { AlertComponent } from './alert/alert.component';
import { AlertEventComponent } from './alertevent/alertevent.component';
import { DeviceStatComponent } from './devicestat/devicestat.component';


@NgModule({
  bootstrap: [App],
  declarations: [
    PasswordToggleDirective,
    TableActions,
    ControlMessagesComponent,
    GenericModal,
    AboutModal,
    ImportFileModal,
    ExportFileModal,
    BlockUIComponent,
    TreeView,
    SpinnerComponent,
    HomeComponent,
    LoginComponent,
    KapacitorComponent,
    KapacitorTasksComponent,
    RangeTimeComponent,
    IfxServerComponent,
    IfxDBComponent,
    IfxMeasurementComponent,
    ProductComponent,
    ProductGroupComponent,
    TemplateComponent,
    OutHTTPComponent,
    AlertComponent,
    AlertEventComponent,
    DeviceStatComponent,
    NavbarComponent,
    SideMenuComponent,
    TableListComponent,
    App,
  ],
  imports: [

    CoreModule,
    CustomPipesModule,
    HttpModule,
    BrowserModule,
    FormsModule,
    ReactiveFormsModule,
    MultiselectDropdownModule,
    ProgressbarModule.forRoot(),
    AccordionModule.forRoot(),
    TooltipModule.forRoot(),
    ModalModule.forRoot(),
    PaginationModule.forRoot(),
    TabsModule.forRoot(),
    BsDropdownModule.forRoot(),
    TimepickerModule.forRoot(),
    Ng2TableModule,
    RouterModule.forRoot(AppRoutes)
  ],
  providers: [
    WindowRef,
    ExportServiceCfg,
    ValidationService,
    BlockUIService
  ],
  entryComponents: [
      BlockUIComponent,
      KapacitorComponent,
      KapacitorTasksComponent,
      RangeTimeComponent,
      IfxServerComponent,
      IfxDBComponent,
      IfxMeasurementComponent,
      ProductComponent,
      ProductGroupComponent,
      TemplateComponent,
      OutHTTPComponent,
      AlertComponent,
      AlertEventComponent,
      DeviceStatComponent
    ]
})
export class AppModule {}
