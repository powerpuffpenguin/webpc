import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { GroupRoutingModule } from './group-routing.module';
import { ListComponent } from './list/list.component';

import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatTreeModule } from '@angular/material/tree';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatDialogModule } from '@angular/material/dialog';

import { SharedModule } from '../shared/shared.module';

import { NodeComponent } from './list/node/node.component';
import { EditComponent } from './dialog/edit/edit.component';
import { AddComponent } from './dialog/add/add.component';
import { DeleteComponent } from './dialog/delete/delete.component';
import { SelectComponent } from './dialog/select/select.component';


@NgModule({
  declarations: [
    ListComponent,
    NodeComponent,
    EditComponent,
    AddComponent,
    DeleteComponent,
    SelectComponent
  ],
  imports: [
    CommonModule, FormsModule,
    MatProgressBarModule, MatButtonModule, MatCardModule,
    MatTreeModule, MatIconModule, MatListModule,
    MatTooltipModule, MatFormFieldModule, MatInputModule,
    MatProgressSpinnerModule, MatDialogModule,
    SharedModule,
    GroupRoutingModule
  ]
})
export class GroupModule { }
