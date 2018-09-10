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
        case 'alertevent-component':
        return this.getAlertEventAvailableActions();
        case 'devicestats-component':
        return this.getDeviceStatsAvailableActions();
        case 'outhttp-component':
        return this.getOutHTTPAvailableActions();
        case 'product-component':
        return this.getProductAvailableActions();
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
      }
    ];
    return tableAvailableActions;
  }
  getOutHTTPAvailableActions (data ? : any) : any {
    let tableAvailableActions = [
    //Remove Action
      {'title': 'Remove', 'content' :
        {'type' : 'button','action' : 'RemoveAllSelected'}
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
