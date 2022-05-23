import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy, AfterViewInit, ElementRef, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ToasterService } from 'angular2-toaster';
import { takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { Closed } from 'src/app/core/utils/closed';
import videojs, { VideoJsPlayer } from "video.js"

import { Manager, Path, Source } from './manager'
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
  private _videojs(access: string) {
    const ctx = this
    videojs(this.playerRef?.nativeElement,
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
        this.on('ended', () => {
          ctx.manager_?.next()
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
    return manager.current == name
  }
}
