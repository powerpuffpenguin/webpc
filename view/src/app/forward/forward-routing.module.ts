import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AccessGuard } from '../core/guard/access.guard';
import { HomeComponent } from './home/home.component';
import { SharedComponent } from './shared/shared.component';

const routes: Routes = [
  {
    path: 'view/:id',
    component: HomeComponent,
    canActivate: [AccessGuard],
  },
  {
    path: 'fs',
    loadChildren: () => import('../forward-fs/forward-fs.module').then(m => m.ForwardFsModule),
  },
  {
    path: 'shared/:id',
    component: SharedComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ForwardRoutingModule { }
