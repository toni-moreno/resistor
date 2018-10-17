import { FormBuilder, Validators, FormArray, FormGroup, FormControl} from '@angular/forms';
import { ValidationService } from './custom-validation/validation.service'

export class AvailableTableActions {

  //AvailableOptions result depeding on component type
  public availableOptions : Array<any>;

  // type can be : device,...
  // data is the passed extraData when declaring AvailableTableActions on each component
  checkComponentType(type, data?) : any {
    switch (type) {
      case 'kapacitor-component':
        return this.getKapacitorAvailableActions();
        case 'alert-component':
        return this.getAlertAvailableActions();
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
    //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
      },
      // Change Property Action
      {'title': 'Change property', 'content' :
        {'type' : 'selector', 'action' : 'ChangeProperty', 'options' : [
          {'title' : 'Active', 'type':'boolean', 'options' : [
            'true','false']
          }
        ]},
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
          }
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
      }
    ];
    return tableAvailableActions;
  }
  getTemplateAvailableActions (data ? : any) : any {
    let tableAvailableActions = [
    //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
      }
    ];
    return tableAvailableActions;
  }

}
