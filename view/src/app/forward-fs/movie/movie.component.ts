import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy, AfterViewInit, ElementRef, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ToasterService } from 'angular2-toaster';
import { takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { Closed } from 'src/app/core/utils/closed';
import videojs, { VideoJsPlayer } from "video.js"
import { DB } from "./db"

import { Manager, Path, Source, Current } from './manager'
var emptySource = new Array<Source>()


interface DownloadAccessResponse {
  access: string
}
@Component({
  selector: 'app-movie',
  templateUrl: './movie.component.html',
  styleUrls: ['./movie.component.scss']
})
export class MovieComponent implements OnInit, OnDestroy, AfterViewInit {
  private closed_ = new Closed()
  constructor(private readonly activatedRoute: ActivatedRoute,
    private readonly httpClient: HttpClient,
    private toasterService: ToasterService,
  ) { }

  ngOnInit(): void {
  }
  ngOnDestroy(): void {
    this.closed_.close()
    this.manager_?.close()
    this.player_?.dispose()
  }
  @ViewChild("player")
  readonly playerRef: ElementRef | undefined
  private manager_: Manager | undefined
  ngAfterViewInit() {
    ServerAPI.v1.sessions.child('download_access').post<DownloadAccessResponse>(this.httpClient, undefined)
      .pipe(takeUntil(this.closed_.observable)).subscribe((resp) => {
        this._videojs(resp.access)
      }, (e) => {
        this.toasterService.pop('error', undefined, e)
      })
  }
  private player_: VideoJsPlayer | undefined
  private _videojs(access: string) {
    const ctx = this
    this.player_ = videojs(this.playerRef?.nativeElement,
      {
        controls: true,
        preload: 'auto',
        autoplay: true,
      },
      function () {
        // 更新字幕設定
        const settings = (this as any).textTrackSettings
        settings.setValues({
          "backgroundColor": "#000",
          "backgroundOpacity": "0",
          "edgeStyle": "uniform",
        })
        settings.updateDisplay()

        ctx._subscribe(this, access)
        if (DB.isSupported) {
          this.on('timeupdate', () => {
            ctx.manager_?.save(this.src(), this.currentTime())
          })
        }
        this.on('ended', () => {
          ctx.manager_?.next()
        })
        this.on('seeking', () => {
          const current = ctx.manager_?.current
          if (current && current.skipTo != 0) {
            current.skipTo = 0
          }
        })
      },
    )
  }
  path: Path | undefined
  private _subscribe(player: VideoJsPlayer, access: string) {
    this.activatedRoute.queryParams.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((params) => {
      const path = new Path(params['id'],
        params['root'],
        params['path'] ?? '/',
      )
      this.path = path
      let manager = this.manager_
      if (!manager) {
        manager = new Manager(player,
          access,
          this.httpClient,
          path,
        )
        this._run(manager)
      } else if (!manager.path.equal(path)) {
        manager.close()
        manager = new Manager(player,
          access,
          this.httpClient,
          path,
        )
        this._run(manager)
      }
      manager.push(params['name'] ?? '')
    })
  }
  private _run(manager: Manager) {
    this.manager_ = manager
    manager.run().catch((e) => {
      if (this.closed_.isClosed) {
        return
      }
      console.log(e)
      if (this.manager_ == manager) {
        this.manager_ = undefined
        this.toasterService.pop('warning',
          undefined,
          `${e}`,
        )
      }
    })
  }
  get source(): Array<Source> {
    const items = this.manager_?.items
    if (items) {
      return items
    }
    return emptySource
  }

  trackById(_: number, source: Source): string { return source.source.name; }
  onClick(name: string) {
    this.manager_?.push(name)
  }
  isPlay(name: string): boolean {
    const manager = this.manager_
    if (!manager) {
      return false
    }
    return manager.currentName == name
  }
  get skipTo(): number {
    return this.manager_?.current?.skipTo ?? 0
  }
  get skipToString(): string {
    const strs = new Array<string>()
    let v = Math.floor(this.skipTo)
    const h = Math.floor(v / 3600)
    if (h > 0) {
      strs.push(h.toString())
      v -= h * 3600
    }
    const m = Math.floor(v / 60)
    if (m > 0) {
      let s = m.toString()
      if (s.length == 1) {
        s = '0' + s
      }
      strs.push(s)
      v -= m * 60
    } else {
      strs.push('00')
    }
    let s = v.toString()
    if (s.length == 1) {
      s = '0' + s
    }
    strs.push(s)
    return strs.join(':')
  }
  onClickSkipTo() {
    const manager = this.manager_
    if (!manager) {
      return
    }
    const current = manager.current
    if (current) {
      manager.skipTo(current)
    }
  }
}
