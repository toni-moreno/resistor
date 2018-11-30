import { FormBuilder, Validators, FormArray, FormGroup, FormControl} from '@angular/forms';
import { ValidationService } from './custom-validation/validation.service'
import { IMultiSelectSettings } from './multiselect-dropdown'

export class AvailableTableActions {

  //AvailableOptions result depeding on component type
  public availableOptions : Array<any>;
  private single_select: IMultiSelectSettings = {
    pullRight: false,
    enableSearch: true,
    checkedStyle: 'glyphicon',
    buttonClasses: 'btn btn-default',
    selectionLimit: 0,
    closeOnSelect: false,
    showCheckAll: true,
    showUncheckAll: true,
    dynamicTitleMaxItems: 3,
    maxHeight: '400px',
    singleSelect: true,
    allowCustomItem: false
  };
  private multi_select: IMultiSelectSettings = {
    pullRight: false,
    enableSearch: true,
    checkedStyle: 'glyphicon',
    buttonClasses: 'btn btn-default',
    selectionLimit: 0,
    closeOnSelect: false,
    showCheckAll: true,
    showUncheckAll: true,
    dynamicTitleMaxItems: 3,
    maxHeight: '400px',
    singleSelect: false,
    allowCustomItem: false
  };

  // type can be : device,...
  // data is the passed extraData when declaring AvailableTableActions on each component
  checkComponentType(type, data?) : any {
    switch (type) {
      case 'kapacitor-component':
        return this.getKapacitorAvailableActions();
        case 'alert-component':
        return this.getAlertAvailableActions(data);
        case 'alerteventhist-component':
        return this.getAlertEventHistAvailableActions();
        case 'alertevent-component':
        return this.getAlertEventAvailableActions();
        case 'devicestats-component':
        return this.getDeviceStatsAvailableActions();
        case 'endpoint-component':
        return this.getEndpointAvailableActions();
        case 'product-component':
        return this.getProductAvailableActions();
        case 'operation-component':
        return this.getOperationAvailableActions();
        case 'rangetime-component':
        return this.getRangeTimeAvailableActions();
        case 'template-component':
        return this.getTemplateAvailableActions();
      default:
        return null;
      }
  }

  constructor (componentType : string, extraData? : any) {
    this.availableOptions = this.checkComponentType(componentType, extraData);
  }

  getKapacitorAvailableActions (data ? : any) : any {
    let tableAvailableActions = [
    //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
      }
    ];
    return tableAvailableActions;
  }

  getAlertAvailableActions (data ? : any) : any {
    let tableAvailableActions = [
      //Deploy Action
      {'title': 'Deploy', 'content' :
        {'type' : 'button','action' : 'DeployAllSelected'}
      },
      //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
      },
      // Change Property Action
      {'title': 'Change property', 'content' :
        {'type' : 'selector', 'action' : 'ChangeProperty', 'options' : [
          {'title' : 'Active', 'type':'boolean', 'options' : [
            'true','false']
          },
          {'title': 'InfluxFilter','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('')
            })
          },
          {'title': 'AlertNotify','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('', ValidationService.uintegerValidator)
            })
          },
          {'title': 'GrafanaServer','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('')
            })
          },
          {'title': 'GrafanaDashLabel','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('')
            })
          },
          {'title': 'GrafanaDashPanelID','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('')
            })
          },
          {'title': 'ExtraLabel','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('')
            })
          },
          {'title': 'ExtraTag','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('')
            })
          },
          {'title' : 'Endpoint', 'type':'multiselector', 'options' :
            data[0], 'settings' : this.multi_select
          },
          {'title' : 'KapacitorID', 'type':'multiselector', 'options' :
            data[1], 'settings' : this.single_select
          },
          {'title' : 'OperationID', 'type':'multiselector', 'options' :
            data[2], 'settings' : this.single_select
          }
        ]},
      },
      //AppendProperty
      {'title': 'Append property', 'content' :
        {'type' : 'selector', 'action' : 'AppendProperty', 'options' : [
          {'title' : 'Endpoint', 'type':'multiselector', 'options' :
            data[0], 'settings' : this.multi_select
          }
          ]
        }
      }
    ];
    return tableAvailableActions;
  }

  getAlertEventHistAvailableActions (data ? : any) : any {
    let tableAvailableActions = [
    //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
      }
    ];
    return tableAvailableActions;
  }

  getAlertEventAvailableActions (data ? : any) : any {
    let tableAvailableActions = [
    //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
      }
    ];
    return tableAvailableActions;
  }

  getDeviceStatsAvailableActions (data ? : any) : any {
    let tableAvailableActions = [
    //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
      },
      // Change Property Action
      {'title': 'Change property', 'content' :
        {'type' : 'selector', 'action' : 'ChangeProperty', 'options' : [
          {'title' : 'Active', 'type':'boolean', 'options' : [
            'true','false']
          },
          {'title' : 'ExceptionID', 'type':'selector', 'options' : [
            '-1','0','1','2']
          },
        ]},
      }
    ];
    return tableAvailableActions;
  }
  getEndpointAvailableActions (data ? : any) : any {
    let tableAvailableActions = [
    //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
      },
      // Change Property Action
      {'title': 'Change property', 'content' :
        {'type' : 'selector', 'action' : 'ChangeProperty', 'options' : [
          {'title' : 'Enabled', 'type':'boolean', 'options' : [
            'true','false']
          }
        ]},
      }
    ];
    return tableAvailableActions;
  }
  getProductAvailableActions (data ? : any) : any {
    let tableAvailableActions = [
    //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
      },
      // Change Property Action
      {'title': 'Change property', 'content' :
        {'type' : 'selector', 'action' : 'ChangeProperty', 'options' : [
          {'title': 'AlertGroups','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('')
            })
          },
          {'title': 'FieldResolutions','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('')
            })
          },
        ]},
      },
      //AppendProperty
      {'title': 'Append property', 'content' :
        {'type' : 'selector', 'action' : 'AppendProperty', 'options' : [
          {'title': 'AlertGroups','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('')
            })
          },
          {'title': 'FieldResolutions','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('', ValidationService.durationValidator)
            })
          },
          ]
        }
      }
    ];
    return tableAvailableActions;
  }
  getOperationAvailableActions (data ? : any) : any {
    let tableAvailableActions = [
    //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
      }
    ];
    return tableAvailableActions;
  }
  getRangeTimeAvailableActions (data ? : any) : any {
    let tableAvailableActions = [
    //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
      },
      // Change Property Action
      {'title': 'Change property', 'content' :
        {'type' : 'selector', 'action' : 'ChangeProperty', 'options' : [
          {'title': 'MinHour','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('', Validators.compose([Validators.required, ValidationService.hourValidator]))
            })
          },
          {'title': 'MaxHour','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('', Validators.compose([Validators.required, ValidationService.hourValidator]))
            })
          },
          {'title': 'WeekDays','type':'input', 'options':
            new FormGroup({
              formControl : new FormControl('', Validators.compose([Validators.required, ValidationService.weekdaysValidator]))
            })
          }
        ]},
      }
    ];
    return tableAvailableActions;
  }
  getTemplateAvailableActions (data ? : any) : any {
    let tableAvailableActions = [
      //Deploy Action
      {'title': 'Deploy', 'content' :
        {'type' : 'button','action' : 'DeployAllSelected'}
      },
      //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
      }
    ];
    return tableAvailableActions;
  }

}
