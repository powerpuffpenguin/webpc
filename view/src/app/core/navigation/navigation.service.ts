import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
@Injectable({
  providedIn: 'root'
})
export class NavigationService {
  private fullscreen_ = new BehaviorSubject<boolean>(false)
  get fullscreenObservable(): Observable<boolean> {
    return this.fullscreen_
  }
  get fullscreen(): boolean {
    return this.fullscreen_.value
  }
  set fullscreen(val: boolean) {
    if (val == this.fullscreen_.value) {
      return
    }
    this.fullscreen_.next(val)
  }

  private target_ = new BehaviorSubject<string>('')
  get targetObservable(): Observable<string> {
    return this.target_
  }
  get target(): string {
    return this.target_.value
  }
  set target(val: string) {
    if (typeof val !== "string") {
      return
    }
    if (val == this.target_.value) {
      return
    }
    this.target_.next(val)
  }
}
