import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
@Injectable({
  providedIn: 'root'
})
export class FullscreenService {

  constructor() { }
  private _subject = new BehaviorSubject<boolean>(false)
  get observable(): Observable<boolean> {
    return this._subject
  }
  set fullscreen(val: boolean) {
    if (val == this._subject.value) {
      return
    }
    this._subject.next(val)
  }
}
