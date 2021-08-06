import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Completer } from '../utils/completer';
import { NetElement } from './tree';
import { ServerAPI } from '../core/api';
import { Observable, Subject } from 'rxjs';
interface ListResponse {
  items: Array<NetElement>
}

@Injectable({
  providedIn: 'root'
})
export class GroupService {
  private completer_: Completer<Array<NetElement>> | undefined
  constructor(private readonly httpClient: HttpClient) { }
  private resetSubject_ = new Subject<boolean>()
  get resetObservable(): Observable<boolean> {
    return this.resetSubject_
  }
  get promise(): Promise<Array<NetElement>> {

    if (this.completer_) {
      return this.completer_.promise
    }
    const completer = new Completer<Array<NetElement>>()
    this.completer_ = completer
    ServerAPI.v1.groups.get<ListResponse>(this.httpClient).toPromise().then((resp) => {
      completer.resolve(resp.items)
    }).catch((e) => {
      completer.reject(e)
      if (completer == this.completer_) {
        this.completer_ = undefined
      }
    })
    return completer.promise
  }
  reset() {
    this.resetSubject_.next(true)
    if (this.completer_) {
      this.completer_ = undefined
    }
  }
}
