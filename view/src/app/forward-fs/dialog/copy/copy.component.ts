import { HttpClient } from '@angular/common/http';
import { Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'src/app/core/toaster.service';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { SessionService } from 'src/app/core/session/session.service';
import { ExistsChoiceComponent } from '../exists-choice/exists-choice.component';
import { Client, Data } from './state';

@Component({
  selector: 'app-copy',
  templateUrl: './copy.component.html',
  styleUrls: ['./copy.component.scss']
})
export class CopyComponent implements OnInit, OnDestroy {
  constructor(private readonly toasterService: ToasterService,
    private readonly i18nService: I18nService,
    private readonly matDialog: MatDialog,
    private readonly httpClient: HttpClient,
    private readonly sessionService: SessionService,
    private readonly matDialogRef: MatDialogRef<CopyComponent>,
    @Inject(MAT_DIALOG_DATA) private readonly data: Data,) { }
  private client_: Client | undefined
  progress = ''
  ngOnInit(): void {
    this.onSubmit()
  }
  get copied(): boolean {
    return this.data.src.copied
  }
  onClose() {
    this.matDialogRef.close(false)
  }
  ngOnDestroy() {
    if (this.client_) {
      this.client_.close()
      this.client_ = undefined
    }
  }
  onSubmit() {
    const data = this.data
    const client = new Client(this.httpClient, this.sessionService,
      data, {
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
      if (this.copied) {
        this.toasterService.pop('success', undefined, this.i18nService.get('Copy file done'))
      } else {
        this.toasterService.pop('success', undefined, this.i18nService.get('Move file done'))
      }
      this.matDialogRef.close(true)
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
      client.close()
      if (this.client_ == client) {
        this.client_ = undefined
      }
      this.matDialogRef.close(false)
    })
  }
}
