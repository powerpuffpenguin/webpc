import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { AnimationsComponent } from './animations/animations.component';
import { IndexComponent } from './index/index.component';
import { RequireComponent } from './require/require.component';
import { ToasterComponent } from './toaster/toaster.component';

const routes: Routes = [
  {
    path: '',
    component: IndexComponent,
  },
  {
    path: 'toaster',
    component: ToasterComponent,
  },
  {
    path: 'require',
    component: RequireComponent,
  },
  {
    path: 'animations',
    component: AnimationsComponent,
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class DevRoutingModule { }
