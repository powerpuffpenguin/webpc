import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';

import { DevRoutingModule } from './dev-routing.module';

import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';

import { IndexComponent } from './index/index.component';
import { ToasterComponent } from './toaster/toaster.component';
import { AnimationsComponent } from './animations/animations.component';


@NgModule({
  declarations: [IndexComponent, ToasterComponent, AnimationsComponent],
  imports: [
    CommonModule, RouterModule, FormsModule,

    MatCardModule, MatButtonModule, MatFormFieldModule,
    MatCheckboxModule, MatSelectModule, MatInputModule,
    MatIconModule, MatTooltipModule,

    DevRoutingModule,
  ]
})
export class DevModule { }
