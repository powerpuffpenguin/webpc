import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { SharedModule } from "../shared/shared.module";

import { ForwardFsRoutingModule } from './forward-fs-routing.module';

import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';

import { RootComponent } from './root/root.component';
import { ListComponent } from './list/list.component';


@NgModule({
  declarations: [
    RootComponent,
    ListComponent
  ],
  imports: [
    CommonModule,
    SharedModule,
    MatProgressBarModule, MatButtonModule, MatCardModule,
    MatIconModule, MatListModule,
    ForwardFsRoutingModule
  ]
})
export class ForwardFsModule { }
