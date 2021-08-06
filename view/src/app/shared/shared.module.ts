import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';

import { MatToolbarModule } from '@angular/material/toolbar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatIconModule } from '@angular/material/icon';
import { MatMenuModule } from '@angular/material/menu';
import { MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider';
import { MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatInputModule } from '@angular/material/input';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTreeModule } from '@angular/material/tree';
import { MatListModule } from '@angular/material/list';
import { MatRadioModule } from '@angular/material/radio';

import { NavigationBarComponent } from './navigation-bar/navigation-bar.component';
import { SignInComponent } from './sign-in/sign-in.component';
import { PasswordComponent } from './password/password.component';
import { TreeComponent } from './tree/tree.component';
import { TreeSelectComponent } from './tree-select/tree-select.component';


@NgModule({
  declarations: [NavigationBarComponent, SignInComponent, PasswordComponent, TreeComponent, TreeSelectComponent],
  imports: [
    CommonModule, RouterModule, FormsModule,

    MatToolbarModule, MatTooltipModule, MatIconModule,
    MatMenuModule, MatButtonModule, MatDividerModule,
    MatDialogModule, MatFormFieldModule, MatCheckboxModule,
    MatInputModule, MatProgressSpinnerModule, MatTreeModule,
    MatRadioModule,
    MatListModule,
  ],
  exports: [NavigationBarComponent, TreeComponent, TreeSelectComponent]
})
export class SharedModule { }
