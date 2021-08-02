import { Injectable } from '@angular/core';
import { HttpEvent, HttpInterceptor, HttpHandler, HttpRequest, HttpErrorResponse, HttpHeaders, HttpClient } from '@angular/common/http'
import { from, Observable } from 'rxjs';
import { catchError, concatAll, map, mapTo } from 'rxjs/operators';
import { Manager, Session } from '../session/session';
import { NetError, resolveError, resolveHttpError } from '../core/restful';
import { Completer } from '../utils/completer';

@Injectable()
export class HeaderInterceptor implements HttpInterceptor {
  constructor(private readonly httpClient: HttpClient) { }
  private completer_: Completer<Session | undefined> | undefined
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
    let first = true
    return next.handle(req.clone({
      headers: headers,
    })).pipe(
      catchError((err, caught) => {
        if (first) {
          // only refresh once
          first = false
          if (session && err instanceof HttpErrorResponse) {
            if (err.status === 401 || err.status === 403) {
              return this._refreshRetry(req, next, session, err.status)
            } else if (err.status === 403) {
              return this._refreshRetry(req, next, session, err.status, resolveError(err).message)
            }
          }
        }
        throw err
      }),
    )
  }
  private _refreshRetry(req: HttpRequest<any>, next: HttpHandler, session: Session, code?: number, msg?: string): Observable<HttpEvent<any>> {
    return from(new Promise<Session>((resolve, reject) => {
      Manager.instance.refresh(this.httpClient, session, code, msg).then((session) => {
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
