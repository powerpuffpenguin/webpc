import { HttpClient } from '@angular/common/http';
import { Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { FileInfo, Dir } from '../../fs';
import { Client } from './state';
import { SessionService } from "src/app/core/session/session.service";
import { ExistsChoiceComponent } from '../exists-choice/exists-choice.component'
interface Target {
  dir: Dir
  source: FileInfo
}
@Component({
  selector: 'app-uncompress',
  templateUrl: './uncompress.component.html',
  styleUrls: ['./uncompress.component.scss']
})
export class UncompressComponent implements OnInit, OnDestroy {
  constructor(private readonly toasterService: ToasterService,
    private readonly i18nService: I18nService,
    private readonly matDialog: MatDialog,
    private readonly httpClient: HttpClient,
    private readonly sessionService: SessionService,
    private readonly matDialogRef: MatDialogRef<UncompressComponent>,
    @Inject(MAT_DIALOG_DATA) public target: Target,) {
  }
  private client_: Client | undefined
  result = false
  progress = ''
  ngOnInit(): void {
    this.result = true
    this.onSubmit()
  }
  onClose() {
    this.matDialogRef.close(this.result)
  }
  ngOnDestroy() {
    if (this.client_) {
      this.client_.close()
      this.client_ = undefined
    }
  }
  onSubmit() {
    const target = this.target
    const dir = target.dir
    const client = new Client(dir.id, this.httpClient, this.sessionService,
      dir.root, dir.dir,
      target.source.name, {
      onProgress: (name: string) => {
        this.progress = name
      },
      onExists: (name) => {
        return this.matDialog.open(ExistsChoiceComponent, {
          data: name,
          disableClose: true,
        }).afterClosed().toPromise()
      },
    },
    )
    this.client_ = client
    client.result.then(() => {
      this.toasterService.pop('success', undefined, this.i18nService.get('Uncompress done'))
      this.matDialogRef.close(true)
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
      client.close()
      if (this.client_ == client) {
        this.client_ = undefined
      }
      this.matDialogRef.close(true)
    })
  }
}
