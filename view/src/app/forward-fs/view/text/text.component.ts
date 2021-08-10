import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { SessionService } from 'src/app/core/session/session.service';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { HttpClient } from '@angular/common/http';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-text',
  templateUrl: './text.component.html',
  styleUrls: ['./text.component.scss']
})
export class TextComponent implements OnInit {


  constructor(private router: Router,
    private route: ActivatedRoute,
    private sessionService: SessionService,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private httpClient: HttpClient,
  ) { }

  ngOnInit(): void {
    //  this.sessionService.observable.subscribe((session) => {
    //       if (this._closed) {
    //         return
    //       }
    //       this._session = session
    //     })
    //     this.sessionService.ready.then(() => {
    //       if (this._closed) {
    //         return
    //       }
    //       const param = this.route.snapshot.queryParamMap
    //       const root = param.get(`root`)
    //       const path = param.get(`path`)
    //       this.root = root
    //       this.filepath = path
    //       this.name = path
    //       const index = path.lastIndexOf('/')
    //       if (index != -1) {
    //         this.dir = path.substring(0, index)
    //         this.name = path.substring(index + 1)
    //       }
    //       this.ready = true
    //       this.onClickLoad()
    //     })
  }
  ngOnDestroy() {
    // this._closed = true
    // if (this._subscription) {
    //   this._subscription.unsubscribe()
    // }
  }
  onPathChange(path: string) {
    // if (typeof path !== "string") {
    //   path = '/'
    // }
    // if (!path.startsWith('/')) {
    //   path = '/' + path
    // }

    // this.router.navigate(['fs', 'list'], {
    //   queryParams: {
    //     root: this.root,
    //     path: path,
    //   }
    // })
  }
  canDeactivate(): boolean {
    // if (this.saving) {
    //   this.toasterService.pop('warning',
    //     undefined,
    //     this.i18nService.get('Wait for data to be saved'),
    //   )
    //   return false
    // }
    return true
  }
  onClickLoad() {
    // this.loading = true
    // ServerAPI.v1.fs.getOne(this.httpClient,
    //   [this.root, this.filepath],
    //   {
    //     responseType: 'text',
    //   },
    // ).then((data) => {
    //   this.val = data
    //   this._val = data
    // }, (e) => {
    //   this.toasterService.pop('error', undefined, e)
    // }).finally(() => {
    //   this.loading = false
    // })
  }
  // get isNotChanged(): boolean {
  //   return this.val == this._val
  // }
  // get isNotCanWrite(): boolean {
  //   if (this._session) {
  //     if (this._session.write || this._session.root) {
  //       return false
  //     }
  //   }
  //   return true
  // }
  onCLickSave() {
    // this.saving = true
    // ServerAPI.v1.fs.putOne(this.httpClient,
    //   [
    //     this.root,
    //     this.filepath,
    //   ],
    //   this.val,
    // ).then(() => {
    //   this.toasterService.pop('success', undefined, this.i18nService.get('Data saved'))
    //   this._val = this.val
    // }, (e) => {
    //   this.toasterService.pop('error', undefined, e)
    // }).finally(() => {
    //   this.saving = false
    // })
  }
}
