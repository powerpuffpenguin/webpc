import { Injectable } from '@angular/core';
import { HttpEvent, HttpInterceptor, HttpHandler, HttpRequest, HttpErrorResponse, HttpHeaders, HttpClient, HttpResponse } from '@angular/common/http'
import { from, Observable } from 'rxjs';
import { catchError, concatAll, tap, map } from 'rxjs/operators';
import { Manager, Session } from '../session/session';
import { Codes, NetError, resolveError } from '../core/restful';
import { Upgraded } from './upgraded';

@Injectable()
export class HeaderInterceptor implements HttpInterceptor {
  constructor(private readonly httpClient: HttpClient) { }
  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    let headers = req.headers

    if (req.method == "GET" || req.method == "HEAD") {
      headers = headers.set('ngsw-bypass', '')
    }
    if (headers.has(`Interceptor`)) {
      const interceptor = headers.get(`Interceptor`)
      headers = headers.delete(`Interceptor`)
      if (interceptor === 'none') {
        return next.handle(req.clone({
          headers: headers,
        }))
      }
    }
    const session = Manager.instance.session
    if (session) {
      if (!headers.has('Authorization')) {
        headers = headers.set('Authorization', `Bearer ${session.access}`)
      }
    }
    return next.handle(req.clone({
      headers: headers,
    })).pipe(
      tap((evt: any) => {
        if (evt.type === 0) {
          return
        }
        if (evt instanceof HttpResponse) {
          const version = evt.headers.get('App-Upgraded')
          if (typeof version === `string`) {
            Upgraded.instance.nextVersion(version)
          } else {
            Upgraded.instance.nextVersion('')
          }
        }
      }),
      catchError((err, caught) => {
          if (session && err instanceof HttpErrorResponse) {
            const e = resolveError(err)
            if (e.grpc == Codes.Unauthenticated) {
              return this._refreshRetry(req, next, session, e)
            } else if (e.grpc == Codes.PermissionDenied && e.message == 'token not exists') {
              Manager.instance.clear(session)
            }
          }
        throw err
      }),
    )
  }
  private _refreshRetry(req: HttpRequest<any>, next: HttpHandler, session: Session, err: NetError): Observable<HttpEvent<any>> {
    return from(new Promise<Session>((resolve, reject) => {
      Manager.instance.refresh(this.httpClient, session, err).then((session) => {
        if (session) {
          resolve(session)
        } else {
          reject()
        }
      }).catch((e) => {
        reject(e)
      })
    })).pipe(
      map((session) => {
        let headers = req.headers
        headers = headers.set('Authorization', `Bearer ${session.access}`)
        return next.handle(req.clone({
          headers: headers,
        }))
      }),
      concatAll(),
    )
  }
}
