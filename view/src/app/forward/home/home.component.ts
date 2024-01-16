import { HttpClient, HttpParams } from '@angular/common/http';
import { Component, ElementRef, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ToasterService } from 'src/app/core/toaster.service';
import { takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { SessionService } from 'src/app/core/session/session.service';
import { Closed } from 'src/app/core/utils/closed';
import { State } from './state';
import * as ClipboardJS from 'clipboard';

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
    private readonly toasterService: ToasterService,
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
      this.url = ServerAPI.forward.v1.forward.websocketURL(id)
    })
  }
  @ViewChild("clipboard")
  private readonly clipboard_: ElementRef | undefined
  private clipboardjs_?: ClipboardJS
  ngAfterViewInit() {
    this.clipboardjs_ = new ClipboardJS(this.clipboard_?.nativeElement).on('success', () => {
      if (this.closed_.isNotClosed) {
        this.toasterService.pop('info', '', "copied")
      }
    }).on('error', (evt: any) => {
      if (this.closed_.isNotClosed) {
        this.toasterService.pop('error', undefined, "copied error")
        console.error('Action:', evt.action)
        console.error('Trigger:', evt.trigger)
      }
    })
  }
  ngOnDestroy() {
    this.state.closed.close()
    this.closed_.close()
    this.clipboardjs_?.destroy()
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
  url = ""
  onCliCkCopyClipboard() {
    const clipboard = this.clipboard_
    if (!clipboard) {
      return
    }
    const element = clipboard.nativeElement
    if (!element) {
      return
    }
    element.setAttribute(
      'data-clipboard-text',
      this.url,
    )
    element.click()
  }
}
