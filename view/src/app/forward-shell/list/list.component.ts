import { HttpClient } from '@angular/common/http';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { ActivatedRoute } from '@angular/router';
import { takeUntil } from 'rxjs/operators';
import { NavigationService } from 'src/app/core/navigation/navigation.service';
import { Closed } from 'src/app/core/utils/closed';
import { DeleteComponent } from '../dialog/delete/delete.component';
import { EditComponent } from '../dialog/edit/edit.component';
import { ListResult } from './load_state';
import { State } from './state';
@Component({
  selector: 'app-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss']
})
export class ListComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  disabled = false
  state = {} as State
  constructor(private readonly activatedRoute: ActivatedRoute,
    private readonly httpClient: HttpClient,
    private readonly matDialog: MatDialog,
    private readonly navigationService: NavigationService,
  ) { }

  ngOnInit(): void {
    this.activatedRoute.params.pipe(
      takeUntil(this.closed_.observable)
    ).subscribe((params) => {
      const id = params['id']
      this.navigationService.target = id
      const state = new State(this.httpClient, id)
      this.state = state
      state.refresh()
    })
  }
  ngOnDestroy() {
    this.closed_.close()
    this.state.closed.close()
    this.navigationService.target = ''
  }
  onClickRefresh() {
    this.state.refresh()
  }
  get ready() {
    return this.state.ready
  }
  get hasErr() {
    return this.state.hasErr
  }
  get errs() {
    return this.state.errs
  }
  get items() {
    return this.state.list.result
  }
  get id() {
    return this.state.target
  }
  onClickEdit(node: ListResult) {
    this.matDialog.open(EditComponent, {
      data: {
        id: this.id,
        result: node,
      },
      disableClose: true,
    })
  }
  onClickDelete(node: ListResult) {
    this.matDialog.open(DeleteComponent, {
      data: {
        id: this.id,
        result: node,
      },
      disableClose: true,
    }).afterClosed().toPromise<boolean>().then((deleted) => {
      if (this.closed_.isClosed || !deleted || typeof deleted !== "boolean") {
        return
      }
      const items = this.state.list.result
      const index = items.indexOf(node)
      if (index > -1) {
        const source = new Array<ListResult>()
        items.splice(index, 1)
        source.push(...items)
        this.state.list.result = source
      }
    })
  }
}
