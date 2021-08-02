import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { PageEvent } from '@angular/material/paginator';
import { ActivatedRoute, Router } from '@angular/router';
import { ToasterService } from 'angular2-toaster';
import { finalize, takeUntil } from 'rxjs/operators';
import { AuthorizationName, ServerAPI } from 'src/app/core/core/api';
import { Closed } from 'src/app/core/utils/closed';
import { AddComponent } from '../dialog/add/add.component';
import { DeleteComponent } from '../dialog/delete/delete.component';
import { EditComponent } from '../dialog/edit/edit.component';
import { PasswordComponent } from '../dialog/password/password.component';
import { Request, Response, Data, DefaultLimit } from './query'

@Component({
  selector: 'app-query',
  templateUrl: './query.component.html',
  styleUrls: ['./query.component.scss']
})
export class QueryComponent implements OnInit, OnDestroy {
  disabled = false
  private request_ = new Request()
  request = new Request()
  lastRequest: Request | undefined
  source = new Array<Data>()
  readonly displayedColumns: string[] = ['id', 'name', 'nickname', 'authorization', 'buttons']
  constructor(
    private readonly httpClient: HttpClient,
    private readonly activatedRoute: ActivatedRoute,
    private readonly router: Router,
    private readonly toasterService: ToasterService,
    private readonly matDialog: MatDialog,
  ) {
  }
  private closed_ = new Closed()
  ngOnInit(): void {
    this.request.name = undefined
    this.activatedRoute.queryParams.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((params) => {
      if (params['limit'] !== undefined) {
        if (this.lastRequest) {
          const request = this.lastRequest
          this.lastRequest = undefined
          if (request.count == 0) {
            request.cloneTo(this.request)
            request.cloneTo(this.request_)
          } else {
            this._query(request)
          }
        } else {
          const request = new Request(params)
          this._query(request)
        }
      }
    })
  }
  ngOnDestroy() {
    this.closed_.close()
    console.log(`ngOnDestroy`)
  }
  private _query(request: Request) {
    if (this.disabled) {
      return
    }
    this.disabled = true
    ServerAPI.v1.users.get<Response>(this.httpClient, {
      params: request.toArgs()
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe((response) => {
      try {
        if ((response.result === 'DATA_COUNT' || response.result === 'COUNT')
          && typeof response.count === "string") {
          const count = parseInt(response.count)
          if (count >= 0) {
            request.count = Math.floor(count)
          }
        }
      } catch (e) {
        console.warn(`parseInt count error : `, e)
      }
      request.cloneTo(this.request)
      request.cloneTo(this.request_)
      this.source = response.data
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
  get isNotQueryChange(): boolean {
    let lname = this.request.name
    let rname = this.request_.name

    return typeof lname === "string" &&
      typeof rname === "string" &&
      lname.trim().toLowerCase() == rname.trim().toLowerCase() &&
      (lname.length == 0 || this.request.nameFuzzy == this.request_.nameFuzzy) &&
      this.request.limit == this.request_.limit && this.request_.offset == 0
  }
  onClickQuery() {
    if (this.disabled || this.isNotQueryChange) {
      return
    }
    const reuqest = new Request()
    this.request.cloneTo(reuqest)
    reuqest.offset = 0
    this.router.navigate(['/user'], {
      queryParams: reuqest.toQuery(),
    })
  }
  onPage(evt: PageEvent) {
    if (this.disabled) {
      return
    }
    const reuqest = new Request()
    this.request_.cloneTo(reuqest)
    reuqest.offset = evt.pageIndex * evt.pageSize
    reuqest.limit = evt.pageSize
    this.lastRequest = reuqest
    this.router.navigate(['/user'], {
      queryParams: reuqest.toQuery(),
    })
  }
  get length(): number {
    const reuqest = this.request
    if (reuqest && reuqest.count >= 0) {
      return reuqest.count
    }
    return 0
  }
  get pageIndex(): number {
    const reuqest = this.request
    if (reuqest) {
      return reuqest.offset / reuqest.limit
    }
    return 0
  }
  get pageSize(): number {
    const reuqest = this.request
    if (reuqest && reuqest.count >= 0) {
      return reuqest.limit
    }
    return DefaultLimit
  }
  getAuthorization(authorization: Array<number>): Array<string> {
    return authorization?.map((v) => {
      return AuthorizationName(v)
    })
  }
  onClickPassword(data: Data) {
    this.matDialog.open(PasswordComponent, {
      data: data,
      disableClose: true,
    })
  }
  onClickEdit(data: Data) {
    this.matDialog.open(EditComponent, {
      data: data,
      disableClose: true,
    })
  }
  onClickDelete(data: Data) {
    this.matDialog.open(DeleteComponent, {
      data: data,
      disableClose: true,
    }).afterClosed().toPromise<boolean>().then((deleted) => {

      if (this.closed_.isClosed || !deleted || typeof deleted !== "boolean") {
        return
      }
      const index = this.source.indexOf(data)

      if (index > -1) {
        const source = new Array<Data>()
        this.source.splice(index, 1)
        source.push(...this.source)
        this.source = source
        if (this.request.count > 0) {
          this.request.count--
        }
        if (this.request_.count > 0) {
          this.request_.count--
        }
      }
    })
  }
  onClickAdd() {
    this.matDialog.open(AddComponent, {
      data: {
        onAdded: (data: Data) => {
          if (this.request.count >= 0) {
            this.request.count++
          }
          if (this.request_.count >= 0) {
            this.request_.count++
          }

          const source = [data]
          source.push(...this.source)
          this.source = source
        },
      },
      disableClose: true,
    })
  }
}
