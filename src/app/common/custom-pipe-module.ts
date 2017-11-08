import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ElapsedSecondsPipe }  from './elapsedseconds.pipe';
import {ObjectParserPipe,SplitCommaPipe } from './custom_pipe';

@NgModule({
  imports: [CommonModule],
  declarations: [ElapsedSecondsPipe,ObjectParserPipe,SplitCommaPipe],
  exports: [ElapsedSecondsPipe,ObjectParserPipe,SplitCommaPipe]
})
export class CustomPipesModule {
}
