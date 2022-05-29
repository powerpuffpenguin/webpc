import { HttpClient } from '@angular/common/http';
import { Component, OnDestroy, OnInit, AfterViewInit, ElementRef, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ToasterService } from 'angular2-toaster';
import { takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { Closed } from 'src/app/core/utils/closed';
import { Path, Manager, Current, Source, FileUrl, } from './manager';
const isHTML5Video = (typeof (document.createElement('video').canPlayType) != 'undefined')
const emptySource = new Array()
interface DownloadAccessResponse {
  access: string
}
interface Params {
  path: Path
  name: string
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
    private readonly router: Router,
  ) { }
  path: Path | undefined
  ngOnInit(): void {
    this.activatedRoute.queryParams.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((params) => {
      const path = new Path(params['id'],
        params['root'],
        params['path'] ?? '/',
      )
      this.path = path
      if (!isHTML5Video) {
        return
      }
      const p: Params = {
        path: path,
        name: params['name'] ?? '',
      }
      if (this.access_) {
        this._play(p)
      } else {
        this.params_ = p
      }
    })
  }
  ngOnDestroy(): void {
    this.closed_.close()
  }
  @ViewChild("player")
  readonly playerRef: ElementRef | undefined
  private video_: HTMLVideoElement | undefined
  ngAfterViewInit(): void {
    const video: HTMLVideoElement | undefined = this.playerRef?.nativeElement
    if (!video) {
      return
    }
    this.video_ = video
    video.onplay = () => {
      if (video.textTracks.length != 0) {
        video.textTracks[0].mode = "showing"
      }
      this.manager_?.saveName()
    }
    video.onended = () => {
      const manager = this.manager_
      if (!manager) {
        return
      }
      const name = manager.next()
      if (name === undefined) {
        return
      }
      const path = manager.path
      this.router.navigate(['forward', 'fs', 'movie'], {
        queryParams: {
          id: path.id,
          root: path.root,
          path: path.path,
          name: name,
        },
      })
    }
    video.ontimeupdate = (evt) => {
      const manager = this.manager_
      if (manager) {
        manager.timeupdate()
      }
    }
    ServerAPI.v1.sessions.child('download_access').post<DownloadAccessResponse>(this.httpClient, undefined)
      .pipe(takeUntil(this.closed_.observable)).subscribe((resp) => {
        this.access_ = resp.access
        this._videojs()
      }, (e) => {
        this.toasterService.pop('error', undefined, e)
      })
  }
  private access_ = ''
  private params_: Params | undefined
  private manager_: Manager | undefined
  private _videojs() {
    const params = this.params_
    if (!params) {
      return
    }
    this.params_ = undefined
    this._play(params)
  }
  private _play(parmas: Params) {
    const video = this.video_
    if (!video) {
      throw new Error("video null")

    }
    const access = this.access_
    const path = parmas.path
    let manager = this.manager_
    if (!manager) {
      manager = new Manager(video,
        access,
        this.httpClient,
        path,
      )
      this._run(manager)
    } else if (!manager.path.equal(path)) {
      manager.close()
      manager = new Manager(video,
        access,
        this.httpClient,
        path,
      )
      this._run(manager)
    }
    manager.push(parmas.name)
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
  get current(): Current | undefined {
    return this.manager_?.current
  }
  get url(): string {
    const manager = this.manager_
    if (manager) {
      const current = manager.current
      if (current) {
        return current.url
      }
    }
    return ''
  }
  get textTracks(): Array<FileUrl> {
    const manager = this.manager_
    if (manager) {
      const current = manager.current
      if (current) {
        return current.source.textTracks
      }
    }
    return emptySource
  }
  get source(): Array<Source> {
    const items = this.manager_?.items
    if (items) {
      return items
    }
    return emptySource
  }
  isPlay(name: string): boolean {
    const manager = this.manager_
    if (!manager) {
      return false
    }
    return manager.current?.name == name
  }
  trackById(_: number, source: Source): string { return source.source.name; }
}
