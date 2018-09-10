//Menu Components  to load them dynamically
import { IfxServerComponent } from '../ifxserver/ifxserver.component';
import { IfxDBComponent } from '../ifxdb/ifxdb.component';
import { IfxMeasurementComponent } from '../ifxmeasurement/ifxmeasurement.component';
import { KapacitorComponent } from '../kapacitor/kapacitor.component';
import { KapacitorTasksComponent } from '../kapacitortasks/kapacitortasks.component';
import { RangeTimeComponent } from '../rangetime/rangetime.component';
import { ProductComponent } from '../product/product.component';
import { ProductGroupComponent } from '../productgroup/productgroup.component';
import { TemplateComponent } from '../template/template.component';
import { EndpointComponent } from '../endpoint/endpoint.component';
import { AlertComponent } from '../alert/alert.component';
import { AlertEventComponent } from '../alertevent/alertevent.component';
import { DeviceStatComponent } from '../devicestat/devicestat.component'
import { NavbarComponent } from './navbar/navbar.component'
import { SideMenuComponent } from './sidemenu/sidemenu.component'

export const HomeComponentConfig: any =
  {
    'name' : 'Influx DB Servers',
    'table-columns' : [
      { 'title': 'ID', 'name': 'ID' },
      { 'title': 'Connection URL', 'name': 'URL' },
      { 'title': 'Admin User', 'name': 'AdminUser' },
      { 'title': 'AdminPasswd', 'name': 'AdminPasswd' },
      { 'title': 'Description', 'name': 'Description' }
    ]
  };

  export const TableRole : string = 'fulledit';
  export const OverrideRoleActions : Array<Object> = [
    {'name':'export', 'type':'icon', 'icon' : 'glyphicon glyphicon-download-alt text-info', 'tooltip': 'Export item'},
    {'name':'view', 'type':'icon', 'icon' : 'glyphicon glyphicon-eye-open text-success', 'tooltip': 'View item'},
    {'name':'test-connection', 'type':'icon', 'icon' : 'glyphicon glyphicon-flash text-info', 'tooltip': 'Test connection'},
    {'name':'edit', 'type':'icon', 'icon' : 'glyphicon glyphicon-edit text-warning', 'tooltip': 'Edit item'},
    {'name':'remove', 'type':'icon', 'icon' : 'glyphicon glyphicon glyphicon-remove text-danger', 'tooltip': 'Remove item'},
    {'name':'importcatalog', 'type':'icon', 'icon' : 'glyphicon glyphicon-import text-default', 'tooltip': 'Import catalog'}
  ]


  export var MenuItems : Array<any> = [
    {'groupName' : 'Runtime', 'icon': 'glyphicon glyphicon-play', 'expanded': true, 'items':
      [
        {'title': 'Alert Events', 'selector' : 'alertevent-component', 'type': 'component', 'data': AlertEventComponent},
        {'title': 'Kapacitor Tasks', 'selector' : 'kapacitortasks-component', 'type': 'component', 'data': KapacitorTasksComponent},
        {'title': 'Agent status', 'selector' : 'runtime', 'data': null}
      ]
    },
    {'groupName' : 'Influx Catalog', 'icon': 'glyphicon glyphicon-play', 'expanded': true, 'items':
    [
      {'title': 'Influx Databases', 'selector' : 'ifdb-component', 'type': 'component', 'data': IfxDBComponent},
      {'title': 'Influx Measurements', 'selector' : 'ifmeasurement-component', 'type': 'component', 'data': IfxMeasurementComponent}
    ]
    },
    {'groupName' : 'External Server Config', 'icon': 'glyphicon glyphicon-play', 'expanded': true, 'items':
    [
      {'title': 'Influx DB Servers ', 'selector' : 'ifxserver-component', 'type': 'component', 'data': IfxServerComponent},
      {'title': 'Kapacitor Backends', 'selector' : 'kapacitor-component', 'type': 'component', 'data': KapacitorComponent},
      {'title': 'Alerting Endpoints', 'selector' : 'endpoint-component', 'type': 'component', 'data': EndpointComponent},
    ]
    }, 
    {'groupName' : 'Configuration', 'icon': 'glyphicon glyphicon-cog', 'expanded': true, 'items':
      [
        {'title': 'RangeTime', 'selector' : 'rangetime-component', 'type': 'component', 'data': RangeTimeComponent},
        {'title': 'Product', 'selector' : 'product-component', 'type': 'component', 'data': ProductComponent},
        {'title': 'Product Groups', 'selector' : 'productgroup-component', 'type': 'component', 'data': ProductGroupComponent},
        {'title': 'Template', 'selector' : 'template-component', 'type': 'component', 'data': TemplateComponent},
        {'title': 'Alert Definition', 'selector' : 'alert-component', 'type': 'component', 'data': AlertComponent},
        {'title': 'Device Stats', 'selector' : 'devicestat-component', 'type': 'component', 'data': DeviceStatComponent},
      ]
    },
    {'groupName' : 'Data Service', 'icon': 'glyphicon glyphicon-paste', 'expanded': true, 'items':
    [
      {'title': 'Export Data ', 'selector' : null, 'type': 'button', 'data': 'exportdata'},
      {'title': 'Import Data', 'selector' : null, 'type': 'button', 'data': 'importdata'},
    ]
    },
  ];

  export const DefaultItem: any = ProductComponent;