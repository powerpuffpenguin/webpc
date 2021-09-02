import { HttpClient, HttpParams } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { SessionService } from 'src/app/core/session/session.service';
import { Closed } from 'src/app/core/utils/closed';
import { State } from './state';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  state = {} as State
  constructor(private readonly activatedRoute: ActivatedRoute,
    private readonly httpClient: HttpClient,
    private readonly sessionService: SessionService,
  ) {
  }
  private accessToken_ = ''
  ngOnInit(): void {
    this.sessionService.observable.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((session) => {
      if (session && session.access) {
        this.accessToken_ = session.access
      }
    })
    this.activatedRoute.params.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((params) => {
      const id = params['id']
      const state = new State(this.httpClient, id)
      this.state = state
      state.refresh()
    })
  }
  ngOnDestroy() {
    this.state.closed.close()
    this.closed_.close()
  }
  onClickRefresh() {
    this.state.refresh()
  }
  get ready() {
    return this.state.ready
  }
  get hasErr() {
    return this.state.hasErr
  }
  get errs() {
    return this.state.errs
  }
  get data() {
    return this.state.data
  }
  get version() {
    return this.state.version
  }
  get startAt() {
    return this.state.startAt
  }
  get vncURL(): string {
    if (!this.data.id) {
      return ''
    }
    let params = new HttpParams({
      fromObject: {
        access_token: this.accessToken_,
      },
    })
    const path = ServerAPI.forward.v1.vnc.httpURL(this.data.id).substring(1)
    params = new HttpParams({
      fromObject: {
        path: `${path}?${params.toString()}`,
      },
    })
    return `/static/noVNC/vnc.html?${params.toString()}`
  }
  get upgraded(): string {
    return this.state?.upgraded?.version ?? ''
  }
}
