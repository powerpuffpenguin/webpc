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
import { State } from './state';

const DefaultValue: any = {}
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
  ) {
  }
  ngOnInit(): void {
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
}
