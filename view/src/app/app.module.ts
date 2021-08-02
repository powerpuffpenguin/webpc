import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FormsModule } from '@angular/forms';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';

import { AppRoutingModule } from './app-routing.module';
import { ToasterModule, ToasterService } from 'angular2-toaster';

import { SharedModule } from "./shared/shared.module";

import { AppComponent } from './app.component';
import { HeaderInterceptor } from './core/interceptor/header.service';
import { HomeComponent } from './app/home/home.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
  ],
  imports: [
    BrowserModule, BrowserAnimationsModule, FormsModule,
    HttpClientModule,
    SharedModule,
    AppRoutingModule, ToasterModule.forRoot()
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
