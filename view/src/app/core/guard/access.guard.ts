import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot, UrlTree } from '@angular/router';
import { ToasterService } from 'src/app/core/toaster.service';
import { Observable } from 'rxjs';
import { first } from 'rxjs/operators';
import { SessionService } from '../session/session.service';

@Injectable({
  providedIn: 'root'
})
export class AccessGuard implements CanActivate {
  constructor(private sessionService: SessionService,
    private toasterService: ToasterService,
  ) {
  }
  canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot,
  ): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    return new Promise<boolean>((resolve) => {
      this.sessionService.observable.pipe(
        first()
      ).subscribe((session) => {
        if (session) {
          resolve(true)
        } else {
          resolve(false)
          this.toasterService.pop('error', undefined, 'Permission denied')
        }
      })
    })
  }
}
