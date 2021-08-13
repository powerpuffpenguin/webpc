import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';

import { SharedModule } from "../shared/shared.module";

import { ForwardFsRoutingModule } from './forward-fs-routing.module';

import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatListModule } from '@angular/material/list';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatRippleModule } from '@angular/material/core';
import { MatMenuModule } from '@angular/material/menu';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatCheckboxModule } from '@angular/material/checkbox';

import { RootComponent } from './root/root.component';
import { ListComponent } from './list/list.component';
import { ManagerComponent } from './manager/manager.component';
import { PathComponent } from './path/path.component';
import { FileComponent } from './file/file.component';
import { TextComponent } from './view/text/text.component';
import { NavigationComponent } from './navigation/navigation.component';
import { AudioComponent } from './view/audio/audio.component';
import { VideoComponent } from './view/video/video.component';
import { ImageComponent } from './view/image/image.component';


@NgModule({
  declarations: [
    RootComponent,
    ListComponent,
    ManagerComponent,
    PathComponent,
    FileComponent,
    TextComponent,
    NavigationComponent,
    AudioComponent,
    VideoComponent,
    ImageComponent
  ],
  imports: [
    CommonModule,
    RouterModule, FormsModule,
    SharedModule,
    MatProgressBarModule, MatButtonModule, MatCardModule,
    MatIconModule, MatListModule, MatTooltipModule,
    MatFormFieldModule, MatInputModule, MatRippleModule,
    MatMenuModule, MatToolbarModule, MatProgressSpinnerModule,
    MatCheckboxModule,
    ForwardFsRoutingModule
  ]
})
export class ForwardFsModule { }
