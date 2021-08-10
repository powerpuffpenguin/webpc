import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { FileInfo, FileType } from '../fs';
import { Router } from '@angular/router';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
export interface NativeEvent extends Event {
  ctrlKey: boolean
  shiftKey: boolean
}
export interface CheckEvent {
  target: FileInfo
  event: NativeEvent
}
@Component({
  selector: 'fs-file',
  templateUrl: './file.component.html',
  styleUrls: ['./file.component.scss']
})
export class FileComponent implements OnInit {
  constructor(private router: Router,
    private toasterService: ToasterService,
    private i18nService: I18nService,
  ) { }
  @Input()
  source = {} as FileInfo
  @Input()
  target = ''
  @Output()
  checkChange = new EventEmitter<CheckEvent>()
  @Output()
  menuChange = new EventEmitter<CheckEvent>()
  ngOnInit(): void {
  }
  get icon(): string {
    if (this.source) {
      return this.source.icon
    }
    return 'insert_drive_file'
  }
  onContextmenu(evt: NativeEvent) {
    evt.stopPropagation()
    this.menuChange.emit({
      event: evt,
      target: this.source,
    })
    return false
  }
  onClick(evt: NativeEvent) {
    evt.stopPropagation()
    this.checkChange.emit({
      event: evt,
      target: this.source,
    })
    return false
  }
  onDbclick() {
    switch (this.source.filetype) {
      case FileType.Dir:
        this.router.navigate(['forward', 'fs', 'list'], {
          queryParams: {
            id: this.target,
            root: this.source.root,
            path: this.source.filename,
          }
        })
        return
      case FileType.Video:
        this.router.navigate(['forward', 'fs', 'view', 'video'], {
          queryParams: {
            id: this.target,
            root: this.source.root,
            path: this.source.filename,
          }
        })
        return
      case FileType.Audio:
        this.router.navigate(['forward', 'fs', 'view', 'audio'], {
          queryParams: {
            id: this.target,
            root: this.source.root,
            path: this.source.filename,
          }
        })
        return
      case FileType.Image:
        this.router.navigate(['forward', 'fs', 'view', 'image'], {
          queryParams: {
            id: this.target,
            root: this.source.root,
            path: this.source.filename,
          }
        })
        return
      case FileType.Text:
        this.router.navigate(['forward', 'fs', 'view', 'text'], {
          queryParams: {
            id: this.target,
            root: this.source.root,
            path: this.source.filename,
          }
        })
        return
      case FileType.Binary:
        this.toasterService.pop('warning',
          this.i18nService.get('Unsupported file type'),
          this.source.filename,
        )
        return
      default:
        this.toasterService.pop('warning',
          this.i18nService.get('Unsupported file type'),
          this.source.filename,
        )
    }
  }
}
