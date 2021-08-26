import { HttpClient } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { NavigationService } from 'src/app/core/navigation/navigation.service';
import { Closed } from 'src/app/core/utils/closed';
import { State } from './state';
@Component({
  selector: 'app-shared',
  templateUrl: './shared.component.html',
  styleUrls: ['./shared.component.scss']
})
export class SharedComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  state = {} as State
  constructor(private readonly activatedRoute: ActivatedRoute,
    private readonly httpClient: HttpClient,
    private readonly navigationService: NavigationService,
  ) { }

  ngOnInit(): void {
    this.activatedRoute.params.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((params) => {
      const id = params['id']
      this.navigationService.target = id
      const state = new State(this.httpClient, id)
      this.state = state
      state.refresh()
    })
  }
  ngOnDestroy() {
    this.state.closed.close()
    this.closed_.close()
    this.navigationService.target = ''
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
  get names() {
    return this.state.mount.name
  }
  get id() {
    return this.state.target
  }
  url(root: string): string {
    return ServerAPI.forward.v1.static.httpURL(this.state.target, root)
  }
}
