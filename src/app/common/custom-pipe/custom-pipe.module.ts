import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ElapsedSecondsPipe } from './elapsedseconds.pipe';
import { ObjectParserPipe } from './objectparser.pipe';
import { SplitCommaPipe } from './commasplit.pipe';
import { LeadingZeroesPipe } from './leadingzeroes.pipe';

@NgModule({
  imports: [CommonModule],
  declarations: [ElapsedSecondsPipe,ObjectParserPipe,SplitCommaPipe,LeadingZeroesPipe],
  exports: [ElapsedSecondsPipe,ObjectParserPipe,SplitCommaPipe,LeadingZeroesPipe]
})
export class CustomPipesModule {
}
