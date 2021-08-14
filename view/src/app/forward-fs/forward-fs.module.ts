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
import { MatButtonToggleModule } from '@angular/material/button-toggle';
import { MatDividerModule } from '@angular/material/divider';
import { MatDialogModule } from '@angular/material/dialog';

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
import { PropertyComponent } from './dialog/property/property.component';
import { NewFileComponent } from './dialog/new-file/new-file.component';
import { NewFolderComponent } from './dialog/new-folder/new-folder.component';
import { RemoveComponent } from './dialog/remove/remove.component';
import { RenameComponent } from './dialog/rename/rename.component';


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
    ImageComponent,
    PropertyComponent,
    NewFileComponent,
    NewFolderComponent,
    RemoveComponent,
    RenameComponent
  ],
  imports: [
    CommonModule,
    RouterModule, FormsModule,
    SharedModule,
    MatProgressBarModule, MatButtonModule, MatCardModule,
    MatIconModule, MatListModule, MatTooltipModule,
    MatFormFieldModule, MatInputModule, MatRippleModule,
    MatMenuModule, MatToolbarModule, MatProgressSpinnerModule,
    MatCheckboxModule, MatButtonToggleModule, MatDividerModule,
    MatDialogModule,
    ForwardFsRoutingModule
  ]
})
export class ForwardFsModule { }
