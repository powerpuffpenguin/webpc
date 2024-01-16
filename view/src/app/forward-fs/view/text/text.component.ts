import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { SessionService } from 'src/app/core/session/session.service';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'src/app/core/toaster.service';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { HttpClient } from '@angular/common/http';
import { Closed } from 'src/app/core/utils/closed';
import { finalize, takeUntil } from 'rxjs/operators';
import { split } from '../../fs'
import { Session } from 'src/app/core/session/session';
@Component({
  selector: 'app-text',
  templateUrl: './text.component.html',
  styleUrls: ['./text.component.scss']
})
export class TextComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  get edit(): boolean {
    return !this.isNotCanWrite && this.rw == 'write'
  }
  rw: 'read' | 'write' = 'read'
  id = ''
  root = ''
  filepath = ''
  dir = ''
  name = ''
  get disabled(): boolean {
    return this.loading_ || this.saving_
  }
  val = ''
  private val_ = ''
  private session_: Session | undefined
  private loading_ = false
  private saving_ = false
  constructor(private router: Router,
    private activatedRoute: ActivatedRoute,
    private sessionService: SessionService,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private httpClient: HttpClient,
  ) { }

  ngOnInit(): void {
    this.sessionService.observable.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((session) => {
      this.session_ = session
    })
    this.activatedRoute.queryParams.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((params) => {
      if (this.disabled) {
        return
      }
      this.id = params['id']
      this.root = params['root']
      this.filepath = params['path']
      const tmp = split(this.filepath)
      this.name = tmp.name
      this.dir = tmp.dir

      this.onClickLoad()
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
  canDeactivate(): boolean {
    if (this.saving_) {
      this.toasterService.pop('warning',
        undefined,
        this.i18nService.get('Wait for data to be saved'),
      )
      return false
    }
    return true
  }
  onClickLoad() {
    this.loading_ = true
    ServerAPI.forward.v1.fs.child('download').get(this.httpClient, {
      responseType: 'text',
      params: {
        slave_id: this.id,
        root: this.root,
        path: this.filepath,
      }
    }).pipe(
      finalize(() => {
        this.loading_ = false
      }),
      takeUntil(this.closed_.observable)
    ).subscribe((data) => {
      this.val = data
      this.val_ = data
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
  get isNotChanged(): boolean {
    return this.val == this.val_
  }
  get isNotCanWrite(): boolean {
    if (this.session_) {
      if (this.session_.write || this.session_.root) {
        return false
      }
    }
    return true
  }
  onCLickSave() {
    this.saving_ = true
    ServerAPI.forward.v1.fs.child('put', this.id, this.root, this.filepath).put(this.httpClient,
      this.val
    ).pipe(
      finalize(() => {
        this.saving_ = false
      }),
      takeUntil(this.closed_.observable)
    ).subscribe(() => {
      this.toasterService.pop('success', undefined, this.i18nService.get('Data saved'))
      this.val_ = this.val
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
