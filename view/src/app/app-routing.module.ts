import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './app/home/home.component';
import { RootGuard } from './core/guard/root.guard';

const routes: Routes = [
  {
    path: '',
    component: HomeComponent,
  },
  {
    path: 'dev',
    loadChildren: () => import('./dev/dev.module').then(m => m.DevModule),
  },
  {
    path: 'content',
    loadChildren: () => import('./content/content.module').then(m => m.ContentModule),
  },
  {
    path: 'logger',
    loadChildren: () => import('./logger/logger.module').then(m => m.LoggerModule),
    canActivate: [RootGuard],
  },
  {
    path: 'user',
    loadChildren: () => import('./user/user.module').then(m => m.UserModule),
    canActivate: [RootGuard],
  },
  {
    path: 'group',
    loadChildren: () => import('./group/group.module').then(m => m.GroupModule),
    canActivate: [RootGuard],
  },
  {
    path: 'forward',
    loadChildren: () => import('./forward/forward.module').then(m => m.ForwardModule),
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { relativeLinkResolution: 'legacy' })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
