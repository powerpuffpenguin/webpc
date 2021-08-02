import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { LoggerRoutingModule } from './logger-routing.module';

import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatCardModule } from '@angular/material/card';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatListModule } from '@angular/material/list';

import { ViewComponent } from './view/view.component';
import { DownloadComponent } from './download/download.component';
import { AttachComponent } from './attach/attach.component';


@NgModule({
  declarations: [ViewComponent, DownloadComponent, AttachComponent],
  imports: [
    CommonModule, FormsModule,
    MatProgressBarModule, MatCardModule, MatSlideToggleModule,
    MatButtonModule, MatIconModule, MatTooltipModule,
    MatFormFieldModule, MatSelectModule, MatListModule,
    LoggerRoutingModule
  ]
})
export class LoggerModule { }
