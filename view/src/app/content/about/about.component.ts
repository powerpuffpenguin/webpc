import { Component, OnInit, OnDestroy } from '@angular/core';
import { ServerAPI } from 'src/app/core/core/api';
import { HttpClient } from '@angular/common/http';
import { ToasterService } from 'src/app/core/toaster.service';
import { takeUntil } from 'rxjs/operators';
import { Closed } from 'src/app/core/utils/closed';

@Component({
  selector: 'app-about',
  templateUrl: './about.component.html',
  styleUrls: ['./about.component.scss']
})
export class AboutComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  licenses: string | undefined
  license: string | undefined
  constructor(
    private readonly httpClient: HttpClient,
    private readonly toasterService: ToasterService,
  ) { }
  ngOnInit(): void {
    ServerAPI.static.license.get(this.httpClient, {
      responseType: 'text',
    }).pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((text) => {
      this.license = text
    }, (e) => {
      this.toasterService.pop('error',
        undefined,
        e,
      )
    })
    ServerAPI.static.licenses.get(this.httpClient, {
      responseType: 'text',
    }).pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((text) => {
      this.licenses = text
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
