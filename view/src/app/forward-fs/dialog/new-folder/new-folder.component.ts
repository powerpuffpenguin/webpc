import { Component, OnInit, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { FileInfo, Dir } from '../../fs';
import { finalize } from 'rxjs/operators';
@Component({
  selector: 'app-new-folder',
  templateUrl: './new-folder.component.html',
  styleUrls: ['./new-folder.component.scss']
})
export class NewFolderComponent implements OnInit {
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialogRef: MatDialogRef<NewFolderComponent>,
    @Inject(MAT_DIALOG_DATA) public data: Dir,
  ) {
  }
  disabled = false

  ngOnInit(): void {
    this.name = this.i18nService.get('New Folder')
  }
  name = ''
  onSubmit() {
    this.disabled = true
    const data = this.data
    ServerAPI.forward.v1.fs.post<FileInfo>(this.httpClient, {
      slave_id: data.id,
      root: data.root,
      dir: data.dir,
      name: this.name,
      file: false,
    }).pipe(
      finalize(() => {
        this.disabled = false
      })
    ).subscribe((data) => {
      this.toasterService.pop('success', undefined, this.i18nService.get(`New Folder Success`))
      const node = new FileInfo(this.data.root, this.data.dir, data)
      this.matDialogRef.close(node)
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
  onClose() {
    this.matDialogRef.close()
  }
}
