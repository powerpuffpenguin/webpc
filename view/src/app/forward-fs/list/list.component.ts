import { HttpClient } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { takeUntil } from 'rxjs/operators';
import { NavigationService } from 'src/app/core/navigation/navigation.service';
import { Closed } from 'src/app/core/utils/closed';
import { State } from './state';

@Component({
  selector: 'app-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss']
})
export class ListComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  state = {} as State
  constructor(private readonly activatedRoute: ActivatedRoute,
    private readonly httpClient: HttpClient,
    private readonly navigationService: NavigationService,
  ) {
  }

  ngOnInit(): void {
    this.activatedRoute.queryParams.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((params) => {
      const id = params['id']
      this.navigationService.target = id
      const root = params['root']
      const path = params['path'] ?? '/'
      const state = new State(this.httpClient, id, root, path)
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
}
