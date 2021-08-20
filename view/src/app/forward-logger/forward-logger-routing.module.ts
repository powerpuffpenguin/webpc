import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { DownloadComponent } from './download/download.component';
import { ViewComponent } from './view/view.component';
const routes: Routes = [
  {
    path: 'view/:id',
    component: ViewComponent,
  },
  {
    path: 'download/:id',
    component: DownloadComponent,
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ForwardLoggerRoutingModule { }
