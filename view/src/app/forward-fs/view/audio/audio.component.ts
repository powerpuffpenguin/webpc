import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { takeUntil } from 'rxjs/operators';
import { SessionService } from 'src/app/core/session/session.service';
import { Closed } from 'src/app/core/utils/closed';
import { Session } from 'src/app/core/session/session';
import { split } from '../../fs'
import { ServerAPI } from 'src/app/core/core/api';
import { HttpParams } from '@angular/common/http';

@Component({
  selector: 'app-audio',
  templateUrl: './audio.component.html',
  styleUrls: ['./audio.component.scss']
})
export class AudioComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  id = ''
  root = ''
  filepath = ''
  dir = ''
  private url_ = ''
  get url(): string {
    const baseURL = this.url_
    if (baseURL != ``) {
      let params: HttpParams
      const session = this.session_
      if (session) {
        params = new HttpParams({
          fromObject: {
            access_token: session.access,
            slave_id: this.id,
            root: this.root,
            path: this.filepath,
          }
        })
      } else {
        params = new HttpParams({
          fromObject: {
            slave_id: this.id,
            root: this.root,
            path: this.filepath,
          }
        })
      }
      return `${baseURL}?${params.toString()}`
    }
    return baseURL
  }
  private session_: Session | undefined
  constructor(private router: Router,
    private activatedRoute: ActivatedRoute,
    private sessionService: SessionService,) { }

  ngOnInit(): void {
    this.sessionService.observable.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((session) => {
      if (session) {
        this.session_ = session
      }
    })
    this.activatedRoute.queryParams.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((params) => {

      this.id = params['id']
      this.root = params['root']
      this.filepath = params['path']
      const tmp = split(this.filepath)
      this.dir = tmp.dir

      this.url_ = ServerAPI.forward.v1.fs.httpURL('download')
    })
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  onPathChange(path: string) {
    if (typeof path !== "string") {
      path = '/'
    }
    if (!path.startsWith('/')) {
      path = '/' + path
    }

    this.router.navigate(['/forward/fs/list'], {
      queryParams: {
        id: this.id,
        root: this.root,
        path: path,
      }
    })
  }
}
