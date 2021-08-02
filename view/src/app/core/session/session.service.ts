import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Session, Manager } from './session';

@Injectable({
  providedIn: 'root'
})
export class SessionService {
  get session(): Session | undefined {
    return Manager.instance.session
  }
  get observable(): Observable<Session | undefined> {
    return Manager.instance.observable
  }
  constructor(private readonly httpClient: HttpClient,
  ) {
    Manager.instance.load()
  }
  get signining(): Observable<boolean> {
    return Manager.instance.signining
  }
  signin(name: string, password: string, remember: boolean): Promise<Session | undefined> {
    return Manager.instance.signin(this.httpClient, name, password, remember)
  }
  signout() {
    Manager.instance.signout(this.httpClient)
  }
}
