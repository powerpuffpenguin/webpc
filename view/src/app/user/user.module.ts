import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { UserRoutingModule } from './user-routing.module';

import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatInputModule } from '@angular/material/input';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatIconModule } from '@angular/material/icon';
import { MatDialogModule } from '@angular/material/dialog';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

import { QueryComponent } from './query/query.component';
import { PasswordComponent } from './dialog/password/password.component';
import { DeleteComponent } from './dialog/delete/delete.component';
import { EditComponent } from './dialog/edit/edit.component';
import { AddComponent } from './dialog/add/add.component';


@NgModule({
  declarations: [QueryComponent, PasswordComponent, DeleteComponent, EditComponent, AddComponent],
  imports: [
    CommonModule, FormsModule,
    MatButtonModule, MatFormFieldModule, MatCheckboxModule,
    MatInputModule, MatPaginatorModule, MatTableModule,
    MatTooltipModule, MatIconModule, MatDialogModule,
    MatProgressSpinnerModule,
    UserRoutingModule
  ]
})
export class UserModule { }
