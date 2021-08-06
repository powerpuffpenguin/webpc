import { HttpClient } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ToasterService } from 'angular2-toaster';
import { takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { Closed } from 'src/app/core/utils/closed';
import { Loader } from 'src/app/core/utils/loader';
import { RequireNet } from 'src/app/core/utils/requirenet';
import { durationString } from 'src/app/core/utils/utils';
interface VersionResponse {
  platform: string
  version: string
}
interface StartAtResponse {
  result: number
}
@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  err: any
  response: VersionResponse | undefined
  startAt: any
  started: string = ''
  loader: Loader | undefined
  constructor(private readonly activatedRoute: ActivatedRoute,
    private readonly httpClient: HttpClient,
    private readonly toasterService: ToasterService,
  ) {
    this.loader = new Loader([

    ])
  }

  ngOnInit(): void {
    this.onClickRefresh()
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  onClickRefresh() {
    const target = this.activatedRoute.snapshot.params["id"]
    ServerAPI.forward.v1.system.child('version').get<VersionResponse>(this.httpClient, {
      params: {
        slave_id: target,
      }
    }).pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((resp) => {
      this.response = resp
    }, (e) => {
      this.err = e
    })
  }

}
