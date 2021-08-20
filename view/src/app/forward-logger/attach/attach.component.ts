import { Component, OnInit, OnDestroy, AfterViewInit, ViewChild, ElementRef } from '@angular/core';
import { debounceTime, filter, map, takeUntil } from 'rxjs/operators';
import { SessionService } from 'src/app/core/session/session.service';
import { Listener } from './listener';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import { Subject } from 'rxjs';
import { Closed } from 'src/app/core/utils/closed';
import { HttpClient } from '@angular/common/http';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'logger-attach',
  templateUrl: './attach.component.html',
  styleUrls: ['./attach.component.scss']
})
export class AttachComponent implements OnInit, OnDestroy, AfterViewInit {
  private token_ = ''
  constructor(private readonly sessionService: SessionService,
    private readonly httpClient: HttpClient,
    private readonly activatedRoute: ActivatedRoute,
  ) { }
  private closed_ = new Closed()
  listener: Listener | undefined
  get isAttach(): boolean {
    return this.listener ? true : false
  }
  get isNotAttach(): boolean {
    return this.listener ? false : true
  }
  checked = true
  id = ''
  ngOnInit(): void {
    this.id = this.activatedRoute.snapshot.params['id']
    this.sessionService.observable.pipe(
      takeUntil(this.closed_.observable),
      filter((session) => {
        if (session?.root) {
          return true
        }
        return false
      }),
      map((session) => session?.access)
    ).subscribe((token) => {
      this.token_ = token ?? ''
    })
  }
  private subject_ = new Subject()
  @ViewChild("xterm")
  xterm: ElementRef | undefined
  private xterm_: Terminal | undefined
  private fitAddon_: FitAddon | undefined
  ngAfterViewInit() {
    // new xterm
    const xterm = new Terminal({
      cursorBlink: true,
      screenReaderMode: true,
      rendererType: 'canvas',
    })
    this.xterm_ = xterm
    // addon
    const fitAddon = new FitAddon()
    this.fitAddon_ = fitAddon
    xterm.loadAddon(fitAddon)
    xterm.loadAddon(new WebLinksAddon())

    xterm.open(this.xterm?.nativeElement)
    fitAddon.fit()
    this.xterm_ = xterm

    // window size change
    this.subject_.pipe(
      debounceTime(100),
      takeUntil(this.closed_.observable),
    ).subscribe((_) => {
      fitAddon.fit()
    })
  }
  onResize() {
    this.subject_.next()
  }
  ngOnDestroy() {
    this.closed_.close()
    this.onClickDetach()
  }
  onClickAttach() {
    if (this.listener) {
      return
    }
    this.listener = new Listener(
      this.httpClient,
      this.sessionService,
      this,
      this.id,
    )
  }
  onClickDetach() {
    if (!this.listener) {
      return
    }
    this.listener.close()
    this.listener = undefined
    this.xterm_?.writeln(`detach logger console`)
  }
  onClickClear() {
    this.xterm_?.clear()
  }
  writeln(text: string, log?: boolean) {
    if (log || this.checked) {
      this.xterm_?.writeln(text)
    }
  }
  write(text: string, log?: boolean) {
    if (log || this.checked) {
      this.xterm_?.write(text)
    }
  }
}
