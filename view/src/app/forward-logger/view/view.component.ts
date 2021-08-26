import { HttpClient } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ToasterService } from 'angular2-toaster';
import { finalize, takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { NavigationService } from 'src/app/core/navigation/navigation.service';
import { SessionService } from 'src/app/core/session/session.service';
import { Closed } from 'src/app/core/utils/closed';
interface Response {
  file?: string
  console?: string
}
@Component({
  selector: 'app-view',
  templateUrl: './view.component.html',
  styleUrls: ['./view.component.scss']
})
export class ViewComponent implements OnInit, OnDestroy {
  constructor(private readonly activatedRoute: ActivatedRoute,
    private readonly httpClient: HttpClient,
    private readonly toasterService: ToasterService,
    private readonly navigationService: NavigationService,
  ) { }
  private closed_ = new Closed()
  ready = false
  err: any
  disabled = false
  private backup_: Response = {}
  data: Response = {}
  readonly levels = ['debug', 'info', 'warn', 'error', 'dpanic', 'panic', 'fatal']
  id = ''
  ngOnInit(): void {
    this.id = this.activatedRoute.snapshot.params['id']
    this.navigationService.target = this.id
    this.load()
  }
  ngOnDestroy() {
    this.closed_.close()
    this.navigationService.target = ''
  }
  load() {
    this.err = null
    this.ready = false
    ServerAPI.forward.v1.logger.child('level').get<Response>(this.httpClient, {
      params: {
        slave_id: this.id,
      }
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.ready = true
      })
    ).subscribe((response) => {
      this.backup_.file = this.data.file = this._formatLevel(response.file)
      this.backup_.console = this.data.console = this._formatLevel(response.console)
    }, (e) => {
      this.err = e
    })
  }
  private _formatLevel(v?: string): string | undefined {
    if (typeof v === "string") {
      v = v.trim().toLowerCase()
      const levels = this.levels
      for (let i = 0; i < levels.length; i++) {
        const element = levels[i]
        if (element == v) {
          return v
        }
      }
    }
    return undefined
  }
  get isFileNotChanged(): boolean {
    return this.data.file == this.backup_.file
  }
  get isConsoleNotChanged(): boolean {
    return this.data.console == this.backup_.console
  }
  onClickResetFile() {
    if (this.disabled) {
      return
    }
    this.data.file = this.backup_.file
  }
  onClickResetConsole() {
    if (this.disabled) {
      return
    }
    this.data.console = this.backup_.console
  }
  onClickSetFile() {
    if (this.disabled || this.isFileNotChanged) {
      return
    }
    this.disabled = true
    ServerAPI.forward.v1.logger.child('level').post(this.httpClient, {
      tag: 'file',
      level: this.data.file,
    }, {
      params: {
        slave_id: this.id,
      },
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      }),
    ).subscribe(() => {
      this.backup_.file = this.data.file
      this.toasterService.pop('success', undefined, `set file level to ${this.data.file} successed`)
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
  onClickSetConsole() {
    if (this.disabled || this.isConsoleNotChanged) {
      return
    }
    this.disabled = true
    ServerAPI.forward.v1.logger.child('level').post(this.httpClient, {
      tag: 'console',
      level: this.data.console,
    }, {
      params: {
        slave_id: this.id,
      }
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      }),
    ).subscribe(() => {
      if (this.closed_.isClosed) {
        return
      }
      this.backup_.console = this.data.console
      this.toasterService.pop('success', undefined, `set console level to ${this.data.console} successed`)
    }, (e) => {
      if (this.closed_.isClosed) {
        return
      }
      this.toasterService.pop('error', undefined, e)
    })
  }
}
