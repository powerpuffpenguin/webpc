import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SharedModule } from "../shared/shared.module";

import { ForwardShellRoutingModule } from './forward-shell-routing.module';

import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatAutocompleteModule } from '@angular/material/autocomplete';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';

import { ListComponent } from './list/list.component';
import { ViewComponent } from './view/view.component';
import { SettingsComponent } from './dialog/settings/settings.component';


@NgModule({
  declarations: [
    ListComponent,
    ViewComponent,
    SettingsComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    SharedModule,
    MatProgressBarModule, MatButtonModule, MatCardModule,
    MatIconModule, MatListModule, MatTooltipModule,
    MatToolbarModule, MatAutocompleteModule, MatFormFieldModule,
    MatInputModule,
    ForwardShellRoutingModule
  ]
})
export class ForwardShellModule { }
