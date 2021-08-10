import { Component, OnDestroy, OnInit, Input, Output, EventEmitter, ElementRef, ViewChild } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatMenuTrigger } from '@angular/material/menu';
import { Router } from '@angular/router';
import { ToasterService } from 'angular2-toaster';
import { fromEvent, Subscription } from 'rxjs';
import { first, takeUntil } from 'rxjs/operators';
import { Session } from 'src/app/core/session/session';
import { SessionService } from 'src/app/core/session/session.service';
import { Closed } from 'src/app/core/utils/closed';
import { CheckEvent } from '../file/file.component';
import { FileInfo, Dir } from '../fs';
import { Box, Point } from './box';
function isObject(object: any): boolean {
  return object !== null && typeof object === "object"
}
const DefaultValue: any = {}
@Component({
  selector: 'fs-manager',
  templateUrl: './manager.component.html',
  styleUrls: ['./manager.component.scss']
})
export class ManagerComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  private session_: Session | undefined
  private subscription_: Subscription | undefined
  constructor(private router: Router,
    private matDialog: MatDialog,
    private sessionService: SessionService,
    private toasterService: ToasterService,
  ) { }

  @Input()
  folder = {} as Dir
  private source_: Array<FileInfo> = []
  private hide_: Array<FileInfo> = []
  @Output()
  sourceChange = new EventEmitter<Array<FileInfo>>()
  @Input('source')
  set source(arrs: Array<FileInfo>) {
    this.source_ = arrs
    this.hide_ = []
    if (arrs && arrs.length > 0) {
      const items = new Array<FileInfo>()
      for (let i = 0; i < arrs.length; i++) {
        if (arrs[i].name.startsWith('.')) {
          continue
        }
        items.push(arrs[i])
      }
      this.hide_ = items
    }
  }
  get source(): Array<FileInfo> {
    return this.all ? this.source_ : this.hide_
  }
  ngOnInit(): void {
    this.sessionService.observable.pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((session) => {
      this.session_ = session
    })
  }
  get access(): string {
    if (this.session_) {
      return this.session_.access
    }
    return ''
  }
  ngOnDestroy() {
    this.closed_.close()
    if (this.subscription_) {
      this.subscription_.unsubscribe()
      this.subscription_ = undefined
    }
  }
  @ViewChild('fs')
  fs = DefaultValue as ElementRef
  @ViewChild('box')
  box = DefaultValue as ElementRef
  private trigger_ = DefaultValue as MatMenuTrigger
  @ViewChild(MatMenuTrigger)
  set trigger(trigger: MatMenuTrigger) {
    if (this.trigger_ && this.trigger_ != DefaultValue) {
      return
    }
    this.trigger_ = trigger
  }
  get trigger(): MatMenuTrigger {
    return this.trigger_
  }

  ctrl = false
  shift = false
  all = false
  onPathChange(path: string) {
    const folder = this.folder
    if (!folder) {
      return
    }

    if (typeof path !== "string") {
      path = '/'
    }
    if (!path.startsWith('/')) {
      path = '/' + path
    }

    this.router.navigate(['forward', 'fs', 'list'], {
      queryParams: {
        id: this.folder.id,
        root: folder.root,
        path: path,
      }
    })
  }
  menuLeft = 0
  menuTop = 0
  onContextmenu(evt: MouseEvent) {
    if (!this.ctrl && !this.shift && !evt.ctrlKey && !evt.shiftKey) {
      this._clearChecked()
    }
    const trigger = this.trigger
    if (trigger) {
      this._openMenu(trigger, evt.clientX, evt.clientY)
    }
    return false
  }
  onContextmenuNode(evt: CheckEvent) {
    if (!evt.target.checked) {
      if (!this.ctrl && !this.shift && !evt.event.ctrlKey && !evt.event.shiftKey) {
        this._clearChecked()
      }
      evt.target.checked = true
    }
    const trigger = this.trigger
    if (trigger) {
      this._openMenu(trigger, (evt.event as any).clientX, (evt.event as any).clientY)
    }
    return false
  }
  private box_ = new Box()
  onStart(evt: MouseEvent) {
    if (evt.button == 2 || evt.ctrlKey || evt.shiftKey || this.ctrl || this.shift) {
      return
    }
    if (this.subscription_) {
      this.subscription_.unsubscribe()
    }
    this.displayBox_ = false
    const doc = this.box.nativeElement
    let start: undefined | Date = new Date()
    this.subscription_ = fromEvent(document, 'mousemove').pipe(
      takeUntil(fromEvent(document, 'mouseup').pipe(first()))
    ).subscribe({
      next: (evt: any) => {
        if (start) {
          const now = new Date()
          const diff = now.getTime() - start.getTime()
          if (diff < 100) {
            return
          }
          this.displayBox_ = true
          start = undefined
          this.box_.setRange(doc)
          this.box_.start = new Point(evt.clientX, evt.clientY)
          this.box_.stop = this.box_.start
          return;
        }
        this.box_.setRange(doc)
        this.box_.stop = new Point(evt.clientX, evt.clientY)
        this.box_.calculate()
      },
      complete: () => {
        this._select()
      },
    })
  }
  private displayBox_ = false
  onClick(evt: any) {
    evt.stopPropagation()
    if (this.displayBox_ || evt.ctrlKey || evt.shiftKey || this.ctrl || this.shift) {
      return
    }
    // 清空選項
    this._clearChecked()
  }
  private _clearChecked() {
    const source = this.source_
    for (let i = 0; i < source.length; i++) {
      if (source[i].checked) {
        source[i].checked = false
      }
    }
  }
  private _select() {
    const arrs = this.box_.checked(this.fs.nativeElement)
    this._clearChecked()
    const source = this.source
    for (let i = 0; i < arrs.length; i++) {
      const index = arrs[i]
      if (index < source.length) {
        source[index].checked = true
      }
    }

    this.box_.reset()
  }
  get x(): number {
    return this.box_.x
  }
  get y(): number {
    return this.box_.y
  }
  get w(): number {
    return this.box_.w
  }
  get h(): number {
    return this.box_.h
  }
  onCheckChange(evt: CheckEvent) {
    if (evt.event.ctrlKey || this.ctrl) {
      evt.target.checked = !evt.target.checked
      return
    }
    let start = -1
    let stop = -1
    let index = -1
    // 清空選項
    const source = this.source
    if (source) {
      for (let i = 0; i < source.length; i++) {
        if (source[i] == evt.target) {
          index = i
        }
        if (source[i].checked) {
          if (start == -1) {
            start = i
          }
          stop = i
        }
        if (source[i].checked) {
          source[i].checked = false
        }
      }
    }
    if (index == -1) {
      return
    }
    // 設置選項
    if ((evt.event.shiftKey || this.shift) && start != -1) {
      if (index <= start) {
        for (let i = index; i <= stop; i++) {
          source[i].checked = true
        }
      } else if (index >= stop) {
        for (let i = start; i <= index; i++) {
          source[i].checked = true
        }
      } else {
        for (let i = start; i <= stop; i++) {
          source[i].checked = true
        }
      }
      return
    }
    source[index].checked = true
  }
  toggleDisplay() {
    this.all = !this.all
    this._clearChecked()
  }
  // 爲 彈出菜單 緩存 選中目標
  target = new Array<FileInfo>()
  private _openMenu(trigger: MatMenuTrigger, x: number, y: number) {
    this.menuLeft = x
    this.menuTop = y
    trigger.openMenu()
    const target = new Array<FileInfo>()
    const source = this.source
    for (let i = 0; i < source.length; i++) {
      if (source[i].checked) {
        target.push(this.source[i])
      }
    }
    this.target = target
  }
  get isNotCanWrite(): boolean {
    if (this.session_) {
      if (this.session_.root) {
        return false
      }
      if (this.session_.write && this.folder.write) {
        return false
      }
    }
    return true
  }
  get isSessionNotCanWrite(): boolean {
    if (this.session_) {
      if (this.session_.root) {
        return false
      }
      if (this.session_.write) {
        return false
      }
    }
    return true
  }
  onClickRename() {
    //   if (this.target && this.target.length == 1) {
    //     const node = this.target[0]
    //     const name = node.name
    //     this.matDialog.open(RenameComponent, {
    //       data: node,
    //       disableClose: true,
    //     }).afterClosed().toPromise().then(() => {
    //       const current = node.name;
    //       if (name == current) {
    //         return
    //       }
    //       if (name.startsWith(`.`)) {
    //         if (!current.startsWith(`.`)) {
    //           if (!this._hide) {
    //             this._hide = new Array<FileInfo>()
    //           }
    //           this._hide.push(node)
    //           this._hide.sort(FileInfo.compare)
    //         }
    //       } else {
    //         if (current.startsWith(`.`)) {
    //           if (this._hide) {
    //             const index = this._hide.indexOf(node)
    //             if (index != -1) {
    //               this._hide.splice(index, 1)
    //             }
    //           }
    //         }
    //       }
    //     })
    //   }
  }
  onClickNewFile() {
    //   if (!this.folder || this._closed) {
    //     return
    //   }
    //   this.matDialog.open(NewFileComponent, {
    //     data: this.folder,
    //     disableClose: true,
    //   }).afterClosed().toPromise().then((fileinfo: FileInfo) => {
    //     if (fileinfo && fileinfo instanceof FileInfo) {
    //       this._pushNode(fileinfo)
    //     }
    //   })
  }
  onClickNewFolder() {
    //   if (!this.folder || this._closed) {
    //     return
    //   }
    //   this.matDialog.open(NewFolderComponent, {
    //     data: this.folder,
    //     disableClose: true,
    //   }).afterClosed().toPromise().then((fileinfo: FileInfo) => {
    //     if (fileinfo && fileinfo instanceof FileInfo) {
    //       this._pushNode(fileinfo)
    //     }
    //   })
  }
  // private _pushNode(fileinfo: FileInfo) {
  //   if (!this._source) {
  //     this._source = new Array<FileInfo>()
  //     this.sourceChange.emit(this._source)
  //   }
  //   this._source.push(fileinfo)
  //   this._source.sort(FileInfo.compare)
  //   if (typeof fileinfo.name === "string" && fileinfo.name.startsWith('.')) {
  //     return
  //   }
  //   if (!this._hide) {
  //     this._hide = new Array<FileInfo>()
  //   }
  //   this._hide.push(fileinfo)
  //   this._hide.sort(FileInfo.compare)
  // }
  onClickProperty() {
    //   if (!this.target || this.target.length == 0) {
    //     return
    //   }
    //   this.matDialog.open(PropertyComponent, {
    //     data: this.target,
    //   })
  }
  onClickRemove() {
    //   const target = this.target
    //   if (!target || target.length == 0) {
    //     return
    //   }
    //   const dir = this.folder
    //   this.matDialog.open(RemoveComponent, {
    //     data: {
    //       dir: dir,
    //       source: target,
    //     },
    //     disableClose: true,
    //   }).afterClosed().toPromise().then((ok) => {
    //     if (ok) {
    //       for (let i = 0; i < target.length; i++) {
    //         const element = target[i]
    //         if (this._source) {
    //           const index = this._source.indexOf(element)
    //           if (index != -1) {
    //             this._source.splice(index, 1)
    //           }
    //         }
    //         if (this._hide) {
    //           const index = this._hide.indexOf(element)
    //           if (index != -1) {
    //             this._hide.splice(index, 1)
    //           }
    //         }
    //       }
    //     }
    //   })
  }
  onClickCompress() {
    //   const target = this.target
    //   if (!target || target.length == 0) {
    //     return
    //   }
    //   const dir = this.folder
    //   this.matDialog.open(CompressComponent, {
    //     data: {
    //       dir: dir,
    //       source: target,
    //     },
    //     disableClose: true,
    //   }).afterClosed().toPromise().then((fileinfo: FileInfo) => {
    //     if (fileinfo instanceof FileInfo) {
    //       this._pushOrUpdate(fileinfo)
    //     }
    //   })
  }
  // private _pushOrUpdate(fileinfo: FileInfo) {
  //   if (!this._source) {
  //     this._source = new Array<FileInfo>()
  //     this.sourceChange.emit(this._source)
  //   }
  //   if (this._source.length == 0) {
  //     this._source.push(fileinfo)
  //   } else {
  //     let ok = false
  //     for (let i = 0; i < this._source.length; i++) {
  //       if (this._source[i].name == fileinfo.name) {
  //         ok = true
  //         this._source[i] = fileinfo
  //         break
  //       }
  //     }
  //     if (!ok) {
  //       this._source.push(fileinfo)
  //       this._source.sort(FileInfo.compare)
  //     }
  //   }

  //   if (!this._hide) {
  //     this._hide = new Array<FileInfo>()
  //   }
  //   if (this._hide.length == 0) {
  //     this._hide.push(fileinfo)
  //   } else {
  //     let ok = false
  //     for (let i = 0; i < this._hide.length; i++) {
  //       if (this._hide[i].name == fileinfo.name) {
  //         ok = true
  //         this._hide[i] = fileinfo
  //         break
  //       }
  //     }
  //     if (!ok) {
  //       this._hide.push(fileinfo)
  //       this._hide.sort(FileInfo.compare)
  //     }
  //   }
  // }
  onClickRefresh() {
    //   const folder = this.folder
    //   if (!folder) {
    //     return
    //   }
    //   this.router.navigate(['fs', 'list'], {
    //     queryParams: {
    //       root: folder.root,
    //       path: folder.dir,
    //       tick: new Date().getTime(),
    //     }
    //   })
  }
  onClickUncompress() {
    //   const target = this.target
    //   if (!target || target.length == 0) {
    //     return
    //   }
    //   const dir = this.folder
    //   this.matDialog.open(UncompressComponent, {
    //     data: {
    //       dir: dir,
    //       source: target[0],
    //     },
    //     disableClose: true,
    //   }).afterClosed().toPromise().then((ok) => {
    //     if (ok) {
    //       this.onClickRefresh()
    //     }
    //   })
  }

  // private _copy(iscopy: boolean): boolean {
  //   const target = this.target
  //   if (!target || target.length == 0) {
  //     return false
  //   }
  //   const folder = this.folder
  //   if (!folder) {
  //     return false
  //   }
  //   const names = new Array<string>()
  //   for (let i = 0; i < target.length; i++) {
  //     names.push(target[i].name)
  //   }
  //   this.fileService.files = {
  //     copy: iscopy,
  //     root: folder.root,
  //     dir: folder.dir,
  //     names: names,
  //   }
  //   return true
  // }
  onClickCopy() {
    //   if (this._copy(true)) {
    //     this.toasterService.pop(`success`, undefined, this.i18nService.get(`File has been copied`))
    //   }
  }
  onClickCut() {
    //   if (this._copy(false)) {
    //     this.toasterService.pop(`success`, undefined, this.i18nService.get(`File has been cut`))
    //   }
  }
  onClickPaste() {
    //   const files = this.fileService.files
    //   if (!files || !isObject(files.names) || !Array.isArray(files.names)) {
    //     return
    //   }
    //   if (files.copy) {
    //     this.matDialog.open(CopyComponent, {
    //       data: {
    //         names: files.names,
    //         src: {
    //           root: files.root,
    //           dir: files.dir,
    //         },
    //         dst: {
    //           root: this.folder.root,
    //           dir: this.folder.dir,
    //         },
    //       },
    //       disableClose: true,
    //     }).afterClosed().toPromise().then(() => {
    //       if (!this._closed) {
    //         this.onClickRefresh()
    //       }
    //     })
    //   } else {
    //     this.matDialog.open(CutComponent, {
    //       data: {
    //         names: files.names,
    //         src: {
    //           root: files.root,
    //           dir: files.dir,
    //         },
    //         dst: {
    //           root: this.folder.root,
    //           dir: this.folder.dir,
    //         },
    //       },
    //       disableClose: true,
    //     }).afterClosed().toPromise().then((ok) => {
    //       if (ok) {
    //         this.fileService.clear(files)
    //       }
    //       if (!this._closed) {
    //         this.onClickRefresh()
    //       }
    //     })
    //   }
  }
  onClickUpload() {
    //   this.matDialog.open(UploadComponent, {
    //     data: {
    //       root: this.folder.root,
    //       dir: this.folder.dir,
    //     },
    //     disableClose: true,
    //   }).afterClosed().toPromise().then(() => {
    //     if (this._closed) {
    //       return
    //     }
    //     this.onClickRefresh()
    //   })
  }
}
