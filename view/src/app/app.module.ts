import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { ToasterModule, ToasterService } from 'angular2-toaster';

import { SharedModule } from "./shared/shared.module";

import { AppComponent } from './app.component';
import { HeaderInterceptor } from './core/interceptor/header.service';
import { HomeComponent } from './app/home/home.component';
import { QueryComponent } from './app/query/query.component';

import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatInputModule } from '@angular/material/input';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatTableModule } from '@angular/material/table';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatIconModule } from '@angular/material/icon';
import { MatDialogModule } from '@angular/material/dialog';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { AddComponent } from './app/dialog/add/add.component';
import { CodeComponent } from './app/dialog/code/code.component';
import { EditComponent } from './app/dialog/edit/edit.component';
import { DeleteComponent } from './app/dialog/delete/delete.component';
import { GroupComponent } from './app/dialog/group/group.component';
import { ServiceWorkerModule } from '@angular/service-worker';
import { environment } from '../environments/environment';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    QueryComponent,
    AddComponent,
    CodeComponent,
    EditComponent,
    DeleteComponent,
    GroupComponent,
  ],
  imports: [
    BrowserModule, BrowserAnimationsModule, RouterModule,
    FormsModule, HttpClientModule,
    SharedModule,
    MatButtonModule, MatFormFieldModule, MatCheckboxModule,
    MatInputModule, MatPaginatorModule, MatTableModule,
    MatTooltipModule, MatIconModule, MatDialogModule,
    MatProgressSpinnerModule,
    AppRoutingModule, ToasterModule.forRoot(), ServiceWorkerModule.register('ngsw-worker.js', {
  enabled: environment.production,
  // Register the ServiceWorker as soon as the app is stable
  // or after 30 seconds (whichever comes first).
  registrationStrategy: 'registerWhenStable:30000'
})
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: HeaderInterceptor,
      multi: true,
    },
    ToasterService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
