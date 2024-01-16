import { Component, OnInit, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'src/app/core/toaster.service';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { FileInfo, Dir } from '../../fs';
import { finalize } from 'rxjs/operators';

@Component({
  selector: 'app-new-file',
  templateUrl: './new-file.component.html',
  styleUrls: ['./new-file.component.scss']
})
export class NewFileComponent implements OnInit {
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialogRef: MatDialogRef<NewFileComponent>,
    @Inject(MAT_DIALOG_DATA) public data: Dir,
  ) {
  }
  disabled = false

  ngOnInit(): void {
    this.name = this.i18nService.get('New File')
  }
  name = ''
  onSubmit() {
    this.disabled = true
    const data = this.data
    ServerAPI.forward.v1.fs.post<FileInfo>(this.httpClient, {
      root: data.root,
      dir: data.dir,
      name: this.name,
      file: true,
    }, {
      params: {
        slave_id: data.id,
      },
    }).pipe(
      finalize(() => {
        this.disabled = false
      })
    ).subscribe((data) => {
      this.toasterService.pop('success', undefined, this.i18nService.get(`New File Success`))
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
