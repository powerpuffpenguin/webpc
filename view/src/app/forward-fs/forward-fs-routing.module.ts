import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ListComponent } from './list/list.component';
import { RootComponent } from './root/root.component';

const routes: Routes = [
  {
    path: ':id',
    component: RootComponent
  },
  {
    path: ':id/list',
    component: ListComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ForwardFsRoutingModule { }
