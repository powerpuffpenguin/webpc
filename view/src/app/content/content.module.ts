import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ContentRoutingModule } from './content-routing.module';

import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';

import { AboutComponent } from './about/about.component';
import { VersionComponent } from './version/version.component';


@NgModule({
  declarations: [AboutComponent, VersionComponent],
  imports: [
    CommonModule,
    MatButtonModule, MatIconModule, MatCardModule,
    ContentRoutingModule
  ]
})
export class ContentModule { }
