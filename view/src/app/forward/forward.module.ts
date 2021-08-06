import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ForwardRoutingModule } from './forward-routing.module';
import { HomeComponent } from './home/home.component';


@NgModule({
  declarations: [
    HomeComponent
  ],
  imports: [
    CommonModule,
    ForwardRoutingModule
  ]
})
export class ForwardModule { }
