import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ForwardRoutingModule } from './forward-routing.module';
import { SharedModule } from "../shared/shared.module";


import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';

import { HomeComponent } from './home/home.component';
import { SharedComponent } from './shared/shared.component';

@NgModule({
  declarations: [
    HomeComponent,
    SharedComponent
  ],
  imports: [
    CommonModule,
    SharedModule,
    MatProgressBarModule, MatButtonModule, MatCardModule,
    MatIconModule, MatListModule,
    ForwardRoutingModule
  ]
})
export class ForwardModule { }
