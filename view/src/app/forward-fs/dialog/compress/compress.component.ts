import { HttpClient } from '@angular/common/http';
import { Component, Inject, OnDestroy, OnInit } from '@angular/core';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'src/app/core/toaster.service';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { FileInfo, Dir, split } from '../../fs';
import { Algorithm, Client } from './state';
import { SessionService } from "src/app/core/session/session.service";
import { ExistsComponent } from '../exists/exists.component';
interface Target {
  dir: Dir
  source: Array<FileInfo>
}
@Component({
  selector: 'app-compress',
  templateUrl: './compress.component.html',
  styleUrls: ['./compress.component.scss']
})
export class CompressComponent implements OnInit, OnDestroy {
  constructor(private readonly toasterService: ToasterService,
    private readonly i18nService: I18nService,
    private readonly matDialog: MatDialog,
    private readonly httpClient: HttpClient,
    private readonly sessionService: SessionService,
    private readonly matDialogRef: MatDialogRef<CompressComponent>,
    @Inject(MAT_DIALOG_DATA) public target: Target,) {
    const source = target.source
    if (source.length == 1) {
      this.name = source[0].name
    } else if (target.dir.dir) {
      this.name = split(target.dir.dir).name
    }
    if (typeof this.name !== "string" || this.name == '') {
      this.name = 'archive'
    }
  }
  private client_: Client | undefined
  result = false
  get disabled(): boolean {
    return this.client_ ? true : false
  }
  name = ''
  progress = ''
  algorithm = Algorithm.TarGZ
  algorithms = [Algorithm.TarGZ, Algorithm.Tar, Algorithm.Zip]
  ngOnInit(): void {
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
    if (this.disabled) {
      return
    }
    this.result = true

    const target = this.target
    const dir = target.dir
    const client = new Client(dir.id, this.httpClient, this.sessionService,
      dir.root, dir.dir,
      this.name, target.source.map((v) => v.name),
      this.algorithm, {
      onProgress: (name: string) => {
        this.progress = name
      },
      onExists: (name) => {
        return this.matDialog.open(ExistsComponent, {
          data: name,
          disableClose: true,
        }).afterClosed().toPromise()
      },
    },
    )
    this.client_ = client
    client.result.then(() => {
      this.toasterService.pop('success', undefined, this.i18nService.get('Compress done'))
      this.matDialogRef.close(true)
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
      client.close()
      if (this.client_ == client) {
        this.client_ = undefined
      }
    })
  }
  getExt(val: Algorithm): string {
    switch (val) {
      case Algorithm.Tar:
        return '.tar'
      case Algorithm.Zip:
        return '.zip'
      default:
        return '.tar.gz'
    }
  }
}
