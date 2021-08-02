import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RootGuard } from '../core/guard/root.guard';
import { AboutComponent } from './about/about.component';
import { VersionComponent } from './version/version.component';

const routes: Routes = [
  {
    path: 'about',
    component: AboutComponent,
  },
  {
    path: 'version',
    component: VersionComponent,
    canActivate: [RootGuard],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ContentRoutingModule { }
