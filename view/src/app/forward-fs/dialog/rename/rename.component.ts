import { Component, OnInit, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Dir, FileInfo } from '../../fs';
import { finalize } from 'rxjs/operators';
interface Target {
  dir: Dir
  target: FileInfo
}
@Component({
  selector: 'app-rename',
  templateUrl: './rename.component.html',
  styleUrls: ['./rename.component.scss']
})
export class RenameComponent implements OnInit {
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialogRef: MatDialogRef<RenameComponent>,
    @Inject(MAT_DIALOG_DATA) public data: Target,
  ) {
  }
  disabled = false
  ngOnInit(): void {
    this.name = this.data.target.name
  }
  name = ''
  get isNotChanged(): boolean {
    return this.name == this.data.target.name
  }
  onSubmit() {
    this.disabled = true
    const dir = this.data.dir
    const data = this.data
    const target = this.data.target
    ServerAPI.forward.v1.fs.child('rename').post(this.httpClient, {
      slave_id: dir.id,
      root: dir.root,
      dir: dir.dir,
      old: target.name,
      current: this.name,
    }).pipe(
      finalize(() => {
        this.disabled = false
      })
    ).subscribe(() => {
      this.toasterService.pop('success', undefined, this.i18nService.get(`Rename Success`))
      target.setName(this.name)
      this.matDialogRef.close()
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
  onClose() {
    this.matDialogRef.close()
  }
}
