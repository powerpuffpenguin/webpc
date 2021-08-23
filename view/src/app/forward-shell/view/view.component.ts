import { HttpClient } from '@angular/common/http';
import { Component, ElementRef, OnDestroy, OnInit, ViewChild, AfterViewInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute, Router } from '@angular/router';
import { BehaviorSubject, fromEvent, of, Subject } from 'rxjs';
import { takeUntil, concatAll, debounceTime } from 'rxjs/operators';
import { FullscreenService } from 'src/app/core/fullscreen/fullscreen.service';
import { SessionService } from 'src/app/core/session/session.service';
import { Closed } from 'src/app/core/utils/closed';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import { SettingsComponent } from '../dialog/settings/settings.component';
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
    private readonly router: Router,
    private readonly matDialog: MatDialog,
    private readonly httpClient: HttpClient,
    private readonly sessionService: SessionService,
    private readonly fullscreenService: FullscreenService,
  ) { }
  fullscreen = false
  ok = false
  duration = ''
  private info_ = new BehaviorSubject<Info | undefined>(undefined)
  get info(): Info | undefined {
    return this.info_.value
  }
  get id(): string {
    return this.target_.value.id
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
    this.fullscreenService.observable.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((ok) => {
      this.fullscreen = ok
    })
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

      const target = this.target_.value
      if (info.id != target.shellid) {
        this.router.navigate(['/forward/shell', target.id, info.id])
      }
    })
    // on change
    setTimeout(() => {
      this.target_.pipe(
        takeUntil(this.closed_.observable),
      ).subscribe((target) => {
        this._connect(xterm, target.id, target.shellid)
      })
    }, 0)
  }
  private _connect(xterm: Terminal, id: string, shellid: string) {
    let shell = this.shell_
    if (shell && shell.id == id && shell.shellid == shellid) {
      return
    } else if (shell) {
      shell.close()
    }
    shell = new Shell(
      this.httpClient, this.sessionService,
      id, shellid,
      xterm, this.info_,
    )
    this.ok = true
    shell.result.catch((e) => {
      if (this.closed_.isNotClosed) {
        console.log(e)
      }
    }).finally(() => {
      xterm.setOption("cursorBlink", false)
      if (this.closed_.isNotClosed && this.shell_ == shell) {
        this.ok = false
        this.shell_ = undefined
      }
    })
    this.shell_ = shell
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
    if (!xterm || !fitAddon || fontSize < 5 || !shell) {
      return
    }
    if (fontSize == xterm.getOption("fontSize")) {
      return
    }
    xterm.setOption("fontSize", fontSize)
    fitAddon.fit()
    shell.sendFontSize(fontSize)
  }
  onClickFontFamily() {
    const xterm = this.xterm_
    const fitAddon = this.fitAddon_
    const shell = this.shell_
    const fontFamily = this.fontFamily
    if (!xterm || !fitAddon || !shell) {
      return
    }
    if (xterm.getOption("fontFamily") == fontFamily) {
      return
    }
    xterm.setOption("fontFamily", fontFamily)
    xterm.resize(1, 1)
    xterm.clear()
    fitAddon.fit()
    shell.sendFontFamily(fontFamily)
  }
  onClickConnect() {
    const xterm = this.xterm_
    if (!xterm) {
      return
    }
    const target = this.target_.value
    xterm.clear()
    this._connect(xterm, target.id, target.shellid)
  }
  onClickSettings() {
    this.matDialog.open(SettingsComponent, {
      data: {
        fontFamily: this.fontFamily,
        onFontFamily: (str: string) => {
          this.fontFamily = str
          this.onClickFontFamily()
        },
        fontSize: this.fontSize,
        onFontSize: (size: number) => {
          this.fontSize = size
          this.onClickFontSize()
        },
      },
    })
  }
  onClickTab(evt: MouseEvent) {
    this._keyboardKeyDown(9, 'Tab', evt)
  }
  onClickCDHome(evt: MouseEvent) {
    this._keyboardKeyDown(192, '~', evt)
  }
  onClickESC(evt: MouseEvent) {
    this._keyboardKeyDown(27, 'Escape', evt)
  }
  onClickArrowUp(evt: MouseEvent) {
    this._keyboardKeyDown(38, 'ArrowUp', evt)
  }
  onClickArrowDown(evt: MouseEvent) {
    this._keyboardKeyDown(40, 'ArrowDown', evt)
  }
  onClickArrowLeft(evt: MouseEvent) {
    this._keyboardKeyDown(37, 'ArrowLeft', evt)
  }
  onClickArrowRight(evt: MouseEvent) {
    this._keyboardKeyDown(39, 'ArrowRight', evt)
  }
  onClickFullscreen(ok: boolean) {
    this.fullscreen = ok
    this.fullscreenService.fullscreen = ok
    this.onResize()
  }
  private _keyboardKeyDown(keyCode: number, key: string, evt: any) {
    const textarea = this.textarea_
    const xterm = this.xterm_
    if (!textarea) {
      return
    }
    textarea.dispatchEvent(new KeyboardEvent('keydown', {
      keyCode: keyCode,
      key: key,
      code: key,
    } as any))
    if (xterm) {
      setTimeout(() => {
        xterm.focus()
      }, 0)
    }
  }
}
