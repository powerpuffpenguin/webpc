import { Component, OnInit, Inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { FileInfo, Dir } from '../../fs';
import { finalize } from 'rxjs/operators';
interface Target {
  dir: Dir
  source: Array<FileInfo>
}
@Component({
  selector: 'app-remove',
  templateUrl: './remove.component.html',
  styleUrls: ['./remove.component.scss']
})
export class RemoveComponent implements OnInit {

  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
    private matDialogRef: MatDialogRef<RemoveComponent>,
    @Inject(MAT_DIALOG_DATA) public target: Target,
  ) { }

  ngOnInit(): void {
  }
  disabled = false
  private result_ = false

  onSubmit() {
    this.disabled = true
    this.result_ = true
    const target = this.target
    const dir = target.dir
    ServerAPI.forward.v1.fs.delete(this.httpClient, {
      params: {
        slave_id: dir.id,
        root: dir.root,
        dir: dir.dir,
        names: target.source.map((v) => v.name),
      }
    }).pipe(
      finalize(() => {
        this.disabled = false
      })
    ).subscribe(() => {
      this.toasterService.pop('success', undefined, this.i18nService.get(`File deleted`))
      this.matDialogRef.close(true)
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
  onClose() {
    this.matDialogRef.close(this.result_)
  }
}
