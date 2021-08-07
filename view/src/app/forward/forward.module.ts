import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ForwardRoutingModule } from './forward-routing.module';
import { SharedModule } from "../shared/shared.module";


import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';

import { HomeComponent } from './home/home.component';

@NgModule({
  declarations: [
    HomeComponent
  ],
  imports: [
    CommonModule,
    SharedModule,
    MatProgressBarModule, MatButtonModule, MatCardModule,

    ForwardRoutingModule
  ]
})
export class ForwardModule { }
