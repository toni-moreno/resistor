import { Pipe, PipeTransform } from '@angular/core';

@Pipe({name: 'splitComma'})
export class SplitCommaPipe implements PipeTransform {
  transform(value) : any {
    let valArray = [];
    valArray = value.split(',');
    console.log(valArray);
    return valArray;
}
}
