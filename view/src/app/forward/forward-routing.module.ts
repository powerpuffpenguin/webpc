import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { FsGuard } from '../core/guard/fs.guard';
import { RootGuard } from '../core/guard/root.guard';
import { ShellGuard } from '../core/guard/shell.guard';
import { HomeComponent } from './home/home.component';
import { SharedComponent } from './shared/shared.component';

const routes: Routes = [
  {
    path: 'view/:id',
    component: HomeComponent,
  },
  {
    path: 'fs',
    loadChildren: () => import('../forward-fs/forward-fs.module').then(m => m.ForwardFsModule),
    canActivate: [FsGuard],
  },
  {
    path: 'logger',
    loadChildren: () => import('../forward-logger/forward-logger.module').then(m => m.ForwardLoggerModule),
    canActivate: [RootGuard],
  },
  {
    path: 'shared/:id',
    component: SharedComponent,
  },
  {
    path: 'shell',
    loadChildren: () => import('../forward-shell/forward-shell.module').then(m => m.ForwardShellModule),
    canActivate: [ShellGuard],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ForwardRoutingModule { }
