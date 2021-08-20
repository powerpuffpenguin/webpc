import { HttpClient, HttpParams } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { filter, finalize, map, takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { SessionService } from 'src/app/core/session/session.service';
import { Closed } from 'src/app/core/utils/closed';
interface Response {
  names?: Array<string>
}
@Component({
  selector: 'app-download',
  templateUrl: './download.component.html',
  styleUrls: ['./download.component.scss']
})
export class DownloadComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  constructor(private readonly activatedRoute: ActivatedRoute,
    private readonly httpClient: HttpClient,
    private readonly sessionService: SessionService,
  ) { }
  err: any
  ready = false
  source = new Array<string>()
  private set_ = new Set<string>()
  private token_ = ''
  private id_ = ''
  ngOnInit(): void {
    this.id_ = this.activatedRoute.snapshot.params['id']
    this.load()
    this.sessionService.observable.pipe(
      takeUntil(this.closed_.observable),
      filter((session) => {
        if (session?.access) {
          return true
        }
        return false
      }),
      map((session) => session?.access)
    ).subscribe((token) => {
      this.token_ = token ?? ''
    })
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  load() {
    this.err = null
    this.ready = false
    ServerAPI.forward.v1.logger.get<Response>(this.httpClient, {
      params: {
        slave_id: this.id_,
      }
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.ready = true
      })
    ).subscribe((response) => {
      if (response && response.names && response.names.length > 0) {
        for (let i = 0; i < response.names.length; i++) {
          const element = response.names[i]
          if (typeof element === "string" && element.length > 0) {
            if (this.set_.has(element)) {
              continue
            }
            this.source.push(element)
            this.set_.add(element)
          }
        }
        this.source.sort()
      }
    }, (e) => {
      this.err = e
    })
  }
  getURL(name: string): string {
    const parms = new HttpParams({
      fromObject: {
        slave_id: this.id_,
        access_token: `${this.token_}`,
      }
    })
    return ServerAPI.forward.v1.logger.httpURL('download', name) + '?' + parms.toString()
  }
}
