import { HttpClient } from '@angular/common/http';
import { Component, ElementRef, OnDestroy, OnInit, ViewChild, AfterViewInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute } from '@angular/router';
import { BehaviorSubject, fromEvent, of, Subject } from 'rxjs';
import { takeUntil, concatAll, debounceTime } from 'rxjs/operators';
import { SessionService } from 'src/app/core/session/session.service';
import { Closed } from 'src/app/core/utils/closed';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import { Info, Shell, Target } from './state';
const DefaultFontFamily = "Lucida Console"
@Component({
  selector: 'app-view',
  templateUrl: './view.component.html',
  styleUrls: ['./view.component.scss']
})
export class ViewComponent implements OnInit, OnDestroy, AfterViewInit {
  private closed_ = new Closed()
  constructor(private readonly activatedRoute: ActivatedRoute,
    private readonly matDialog: MatDialog,
    private readonly httpClient: HttpClient,
    private readonly sessionService: SessionService,
  ) { }
  fullscreen = false
  ok = false
  duration = ''
  private info_ = new BehaviorSubject<Info | undefined>(undefined)
  get info(): Info | undefined {
    return this.info_.value
  }
  private xterm_: Terminal | undefined
  private fitAddon_: FitAddon | undefined
  private textarea_: Document | undefined
  private shell_: Shell | undefined
  fontSize = 15
  fontFamily = DefaultFontFamily
  private target_ = new BehaviorSubject<Target>({ id: '', shellid: '' })
  private resize_ = new Subject()
  ngOnInit(): void {
    this.activatedRoute.params.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((params) => {
      const id = params["id"]
      const shellid = params["shellid"]
      if (typeof id === "string" && typeof shellid === "string") {
        this.target_.next({
          id: id,
          shellid: shellid,
        })
      }
    })
    of(
      of(true),
      fromEvent(window, 'resize'),
    ).pipe(
      concatAll(),
      takeUntil(this.closed_.observable),
    ).subscribe(() => {
      const vh = window.innerHeight * 0.01;
      document.documentElement.style.setProperty('--vh', `${vh}px`);
    })
  }
  ngOnDestroy() {
    this.closed_.close()
    if (this.shell_) {
      this.shell_.close()
      this.shell_ = undefined
    }
    if (this.xterm_) {
      this.xterm_.dispose()
      this.xterm_ = undefined
    }
  }
  @ViewChild("xterm")
  xterm: ElementRef | undefined
  @ViewChild("view")
  view: ElementRef | undefined
  ngAfterViewInit() {
    if (!this.xterm) {
      return
    }
    // create xterm
    const xterm = new Terminal({
      cursorBlink: true,
      screenReaderMode: true,
      fontFamily: this.fontFamily,
      rendererType: 'canvas',
    })
    this.xterm_ = xterm
    this.fontSize = xterm.getOption("fontSize")

    const fitAddon = new FitAddon()
    this.fitAddon_ = fitAddon
    xterm.loadAddon(fitAddon)
    xterm.loadAddon(new WebLinksAddon())
    xterm.open(this.xterm.nativeElement)
    this.textarea_ = this.xterm.nativeElement.querySelector('textarea')
    fitAddon.fit()

    // fix resize
    this.resize_.pipe(
      debounceTime(100),
      takeUntil(this.closed_.observable),
    ).subscribe((_) => {
      fitAddon.fit()
    })

    xterm.onData((data) => {
      this.shell_?.send(data)
    })
    xterm.onResize((evt) => {
      this.shell_?.sendResize(evt.cols, evt.rows)
    })

    this.info_.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((info) => {
      if (!info) {
        return
      }
      this._fontSize(info.fontSize)
      if (this.fontFamily != info.fontFamily) {
        this.fontFamily = info.fontFamily
        console.log(`set font`, this.fontFamily)
        xterm.setOption("fontFamily", this.fontFamily)
        xterm.resize(1, 1)
        fitAddon.fit()
      }
    })
    // on change
    this.target_.pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((target) => {
      let shell = this.shell_
      if (shell && shell.id == target.id && shell.shellid == target.shellid) {
        return
      } else if (shell) {
        shell.close()
      }
      shell = new Shell(
        this.httpClient, this.sessionService,
        target.id, target.shellid,
        xterm, this.info_,
      )
      this.ok = true
      shell.result.catch((e) => {
        if (this.closed_.isClosed) {
          console.log(e)
        }
      }).finally(() => {
        if (this.closed_.isClosed && this.shell_ == shell) {
          this.ok = false
          this.shell_ = undefined
        }
      })
      this.shell_ = shell
    })
  }
  onResize() {
    this.resize_.next(new Date())
  }
  private _fontSize(fontSize: number) {
    if (typeof fontSize !== "number") {
      return
    }
    fontSize = Math.floor(fontSize)
    if (fontSize < 5 || fontSize == this.fontSize) {
      return
    }
    this.fontSize = fontSize

    const xterm = this.xterm_
    const fitAddon = this.fitAddon_
    const shell = this.shell_
    if (!xterm || !fitAddon || this.fontSize < 5 || !shell) {
      return
    }
    if (fontSize == xterm.getOption("fontSize")) {
      return
    }
    xterm.setOption("fontSize", fontSize)
  }
  onClickFontSize() {
    const xterm = this.xterm_
    const fitAddon = this.fitAddon_
    const shell = this.shell_
    const fontSize = this.fontSize
    if (!xterm || !fitAddon || this.fontSize < 5 || !shell) {
      return
    }
    if (fontSize == xterm.getOption("fontSize")) {
      return
    }
    xterm.setOption("fontSize", fontSize)
    fitAddon.fit()
    shell.sendFontSize(fontSize)
  }
  onClickFullscreen(ok: boolean) {

  }
  onClickConnect() {

  }
}
