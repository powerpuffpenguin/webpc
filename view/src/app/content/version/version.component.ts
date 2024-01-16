import { HttpClient } from '@angular/common/http';
import { Component, OnInit, VERSION, OnDestroy } from '@angular/core';
import { ToasterService } from 'src/app/core/toaster.service';
import { timer } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { Closed } from 'src/app/core/utils/closed';
import { datetimeString, usedString } from 'src/app/core/utils/datetime';
interface VersionResponse {
  platform: string
  version: string
  date?: string
  commit?: string
}
interface StartAtResponse {
  result: string | number
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
    ServerAPI.v1.system.child('start_at').get<StartAtResponse>(this.httpClient).pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((response) => {
      const startAt = typeof response.result === "number" ? response.result : parseInt(response.result)
      if (Number.isSafeInteger(startAt)) {
        this.startAt = datetimeString(new Date(startAt * 1000))
        timer(0, 1000).pipe(takeUntil(this.closed_.observable)).subscribe({
          next: () => {
            let used = Math.floor(Date.now() / 1000)
            if (used > startAt) {
              used -= startAt
            } else {
              used = 0
            }
            this.started = usedString(used)
          },
        })
      }
    }, (e) => {
      this.toasterService.pop('error',
        undefined,
        e,
      )
    })
  }
  ngOnDestroy() {
    this.closed_.close()
  }
}
