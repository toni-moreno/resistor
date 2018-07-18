import { Pipe, PipeTransform } from '@angular/core';

@Pipe({name: 'leadingZeroes'})
export class LeadingZeroesPipe implements PipeTransform {
  transform(value) : any {
    if (value) {
      return ("000"+value.toString()).slice(-3);
    }
    return "";
}

}
