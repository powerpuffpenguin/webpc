import { HttpClient } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { interval } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { Closed } from 'src/app/core/utils/closed';
import { Loader } from 'src/app/core/utils/loader';
import { RequireState } from 'src/app/core/utils/requirenet';
import { durationString } from 'src/app/core/utils/utils';
import { VersionState, VersionResponse, StartAtState, StartAtResponse, DataState, DataResponse } from './load_state'

interface Error {
  id: string
  err: any
}
const DefaultValue: any = {}
@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  private target_ = ''
  ready = false
  loader: Loader = DefaultValue
  moment: any
  data: DataResponse = DefaultValue
  version: VersionResponse = DefaultValue
  startAt: StartAtResponse = DefaultValue
  errs: Array<Error> = []
  hasErr = false
  constructor(private readonly activatedRoute: ActivatedRoute,
    private readonly httpClient: HttpClient,
  ) {
  }
  ngOnInit(): void {
    const opts = {
      httpClient: this.httpClient,
      target: this.activatedRoute.snapshot.params["id"],
      cancel: this.closed_.observable,
    }
    this.target_ = opts.target
    this.loader = new Loader([
      new RequireState('moment', (moment) => {
        this.moment = moment
      }, (e) => {
        this.errs.push({
          id: 'moment',
          err: e,
        })
      }),
      new DataState(opts, (data) => {
        this.data = data
      }, (e) => {
        this.errs.push({
          id: 'data',
          err: e,
        })
      }),
      new VersionState(opts, (data) => {
        this.version = data
      }, (e) => {
        this.errs.push({
          id: 'VersionState',
          err: e,
        })
      }),
      new StartAtState(opts, (data) => {
        this.startAt = data
      }, (e) => {
        this.errs.push({
          id: 'StartAtState',
          err: e,
        })
      }),
    ])
    this.onClickRefresh()
    this.onClickRefresh()
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  onClickRefresh() {
    this.hasErr = false
    this.errs = []
    this.loader.load().then(() => {
      const moment = this.moment
      const startAt = this.startAt
      startAt.at = moment.unix(startAt.result)
      const d = moment.duration(moment.unix(moment.now() / 1000).diff(startAt.at))
      startAt.started = durationString(d)


      interval(1000).pipe(
        takeUntil(this.closed_.observable),
      ).subscribe(() => {
        const startAt = this.startAt
        const d = moment.duration(moment.unix(moment.now() / 1000).diff(startAt.at))
        startAt.started = durationString(d)
      })
      this.ready = true
    }).catch((_) => {
      this.hasErr = true
    })
  }
}
