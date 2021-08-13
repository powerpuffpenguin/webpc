import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ListComponent } from './list/list.component';
import { RootComponent } from './root/root.component';
import { AudioComponent } from './view/audio/audio.component';
import { ImageComponent } from './view/image/image.component';
import { TextComponent } from './view/text/text.component';
import { TextGuard } from './view/text/text.guard';
import { VideoComponent } from './view/video/video.component';

const routes: Routes = [
  {
    path: '',
    component: RootComponent
  },
  {
    path: 'list',
    component: ListComponent,
  },
  {
    path: 'view/video',
    component: VideoComponent,
  },
  {
    path: 'view/audio',
    component: AudioComponent,
  },
  {
    path: 'view/image',
    component: ImageComponent,
  },
  {
    path: 'view/text',
    component: TextComponent,
    canDeactivate: [TextGuard],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ForwardFsRoutingModule { }
