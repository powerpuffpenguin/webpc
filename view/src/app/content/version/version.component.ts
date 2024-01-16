import { HttpClient } from '@angular/common/http';
import { Component, OnInit, VERSION, OnDestroy } from '@angular/core';
import { ToasterService } from 'src/app/core/toaster.service';
import { interval } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { Closed } from 'src/app/core/utils/closed';
import { RequireNet } from 'src/app/core/utils/requirenet';
import { durationString } from 'src/app/core/utils/utils';
interface VersionResponse {
  platform: string
  version: string
  date?: string
  commit?: string
}
interface StartAtResponse {
  result: number
}

@Component({
  selector: 'app-version',
  templateUrl: './version.component.html',
  styleUrls: ['./version.component.scss']
})
export class VersionComponent implements OnInit, OnDestroy {
  VERSION = VERSION
  response: VersionResponse | undefined
  private closed_ = new Closed()
  startAt: any
  started: string = ''
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
  ) { }
  ngOnInit(): void {
    ServerAPI.v1.system.child('version').get<VersionResponse>(this.httpClient).pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((response) => {
      this.response = response
    }, (e) => {
      this.toasterService.pop('error',
        undefined,
        e,
      )
    })

    RequireNet('moment').then((moment) => {
      ServerAPI.v1.system.child('start_at').get<StartAtResponse>(this.httpClient).pipe(
        takeUntil(this.closed_.observable),
      ).subscribe((response) => {
        this.startAt = moment.unix(response.result)
        const d = moment.duration(moment.unix(moment.now() / 1000).diff(this.startAt))
        this.started = durationString(d)
      }, (e) => {
        this.toasterService.pop('error',
          undefined,
          e,
        )
      })
      interval(1000).pipe(
        takeUntil(this.closed_.observable),
      ).subscribe(() => {
        if (this.startAt) {
          const d = moment.duration(moment.unix(moment.now() / 1000).diff(this.startAt))
          this.started = durationString(d)
        }
      })
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
  ngOnDestroy() {
    this.closed_.close()
  }
}
