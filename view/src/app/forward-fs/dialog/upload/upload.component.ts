import { HttpClient } from '@angular/common/http';
import { Component, ElementRef, Inject, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { fromEvent } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { Closed } from 'src/app/core/utils/closed';
import { EventCode } from '../event';
import { Uploader } from './uploader';
import { UploadFile, Dir, States, Workers } from './state';
class Source {
  private keys_ = new Set<string>()
  private items_ = new Array<UploadFile>()

  push(uploadFile: UploadFile) {
    if (!uploadFile.file.size) {
      console.warn(`not support size 0`, uploadFile.file)
      return
    }
    const key = uploadFile.key
    if (this.keys_.has(key)) {
      return
    }
    this.keys_.add(key)
    this.items_.push(uploadFile)
  }
  get source(): Array<UploadFile> {
    return this.items_
  }
  clear() {
    const arrs = this.items_
    for (let i = arrs.length - 1; i >= 0; i--) {
      const node = arrs[i]
      if (node.isWorking()) {
        continue
      }
      arrs.splice(i, 1)
      this.keys_.delete(node.key)
    }
  }
  delete(uploadFile: UploadFile) {
    if (!uploadFile || uploadFile.isWorking()) {
      return
    }
    const index = this.items_.indexOf(uploadFile)
    if (index == -1) {
      return
    }
    this.items_.splice(index, 1)
    this.keys_.delete(uploadFile.key)
  }
  get(): UploadFile | undefined {
    const arrs = this.items_
    let find: UploadFile | undefined
    for (let i = 0; i < arrs.length; i++) {
      const element = arrs[i]
      if (element.state == States.Nil) {
        find = element
        break
      }
    }
    return find
  }
  prepare() {
    for (let i = 0; i < this.items_.length; i++) {
      const element = this.items_[i]
      if (element.state == States.Error) {
        element.state = States.Nil
      }
    }
  }
}
@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.scss']
})
export class UploadComponent implements OnInit, OnDestroy {
  private result_ = false
  private source_ = new Source()
  get source(): Array<UploadFile> {
    return this.source_.source
  }
  disabled = false
  dragover = false
  private closed_ = new Closed()
  private uploader_: Uploader | undefined
  constructor(
    private readonly matDialog: MatDialog,
    private readonly matDialogRef: MatDialogRef<UploadComponent>,
    @Inject(MAT_DIALOG_DATA) public readonly data: Dir,
    private readonly httpClient: HttpClient,) { }

  ngOnInit(): void {
    fromEvent(document, 'drop').pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((evt) => {
      evt.preventDefault()
    })
    fromEvent(document, 'dragover').pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((evt) => {
      evt.preventDefault()
    })
  }
  ngOnDestroy() {
    this.closed_.close()
    if (this.uploader_) {
      this.uploader_.close()
      this.uploader_ = undefined
    }
  }
  @ViewChild("drop")
  private drop_: ElementRef | undefined
  ngAfterViewInit() {
    const drop = this.drop_
    if (!drop) {
      return
    }
    fromEvent(drop.nativeElement, 'dragover').pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((evt: any) => {
      this.dragover = true
      evt.stopPropagation()
      evt.preventDefault()
    })
    fromEvent(drop.nativeElement, 'dragenter').pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((evt: any) => {
      this.dragover = true
      evt.stopPropagation()
      evt.preventDefault()
    })
    fromEvent(drop.nativeElement, 'dragexit').pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((evt: any) => {
      this.dragover = false
      evt.stopPropagation()
      evt.preventDefault()
    })
    fromEvent(drop.nativeElement, 'drop').pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((evt: any) => {
      this.dragover = false
      evt.stopPropagation()
      evt.preventDefault()
      this._drop(evt)
    })
  }
  private _drop(evt: any) {
    let dataTransfer = evt.dataTransfer
    if (!dataTransfer) {
      if (evt.originalEvent) {
        dataTransfer = evt.originalEvent.dataTransfer
        console.warn(`use evt.originalEvent`)
      }
    }
    if (!dataTransfer) {
      console.warn(`dataTransfer nil`)
      return
    }
    if (!dataTransfer.files) {
      return
    }
    for (let i = 0; i < dataTransfer.files.length; i++) {
      this.source_.push(new UploadFile(dataTransfer.files[i]))
    }
  }
  onClose() {
    this.matDialogRef.close(this.result_)
  }
  onClickClear() {
    this.source_.clear()
  }
  onAdd(evt: any) {
    if (evt.target.files) {
      for (let i = 0; i < evt.target.files.length; i++) {
        const element = evt.target.files[i]
        this.source_.push(new UploadFile(element))
      }
    }
  }
  onClickDelete(uploadFile: UploadFile) {
    this.source_.delete(uploadFile)
  }
  onClickReset(uploadFile: UploadFile) {
    if (uploadFile.isSkip()) {
      uploadFile.state = States.Nil
    }
  }
  onClickStart() {
    if (this.disabled) {
      return
    }
    this.result_ = true
    this.disabled = true
    this._run().finally(() => {
      this.disabled = false
    })
  }
  private async _run() {
    try {
      this.source_.prepare()
      let style = EventCode.Universal
      while (this.closed_.isNotClosed) {
        // get upload file
        const uploadFile = this.source_.get()
        if (!uploadFile) {
          break
        }
        const uploader = new Uploader(
          this.data, uploadFile,
          this.httpClient,
          this.matDialog,
        )
        uploader.style = style
        this.uploader_ = uploader
        await uploader.serve()
        style = uploader.style
      }
    } catch (e) {
    }
  }
}
