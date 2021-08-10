import { Injectable } from '@angular/core';
import { CanDeactivate, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Observable } from 'rxjs';
import { TextComponent } from './text.component'

@Injectable({
  providedIn: 'root'
})
export class TextGuard implements CanDeactivate<TextComponent> {
  constructor() { }

  canDeactivate(
    component: TextComponent,
    currentRoute: ActivatedRouteSnapshot,
    currentState: RouterStateSnapshot,
    nextState: RouterStateSnapshot
  ): Observable<boolean> | Promise<boolean> | boolean {
    return component.canDeactivate()
  }
}
