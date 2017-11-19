import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ElapsedSecondsPipe } from './elapsedseconds.pipe';
import { ObjectParserPipe } from './objectparser.pipe';
import { SplitCommaPipe } from './commasplit.pipe';

@NgModule({
  imports: [CommonModule],
  declarations: [ElapsedSecondsPipe,ObjectParserPipe,SplitCommaPipe],
  exports: [ElapsedSecondsPipe,ObjectParserPipe,SplitCommaPipe]
})
export class CustomPipesModule {
}
