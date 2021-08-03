import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatCheckboxChange } from '@angular/material/checkbox';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'angular2-toaster';
import { finalize, takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Closed } from 'src/app/core/utils/closed';
import { Data } from '../../query/query';

@Component({
  selector: 'app-edit',
  templateUrl: './edit.component.html',
  styleUrls: ['./edit.component.scss']
})
export class EditComponent implements OnInit, OnDestroy {
  disabled = false
  name = ''
  description = ''
  private closed_ = new Closed()
  constructor(@Inject(MAT_DIALOG_DATA) public readonly data: Data,
    private httpClient: HttpClient,
    private toasterService: ToasterService,
    private matDialogRef: MatDialogRef<EditComponent>,
    private i18nService: I18nService,
  ) {
    this.name = data.name
    this.description = data.description
  }

  ngOnInit(): void {
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  onClose() {
    this.matDialogRef.close()
  }
  get isNotChanged(): boolean {
    return this.name.trim() == this.data.name.trim() &&
      this.description.trim() == this.data.description.trim()
  }
  onSubmit() {
    if (this.disabled || this.isNotChanged) {
      return
    }
    this.disabled = true
    const name = this.name.trim()
    const description = this.description.trim()
    ServerAPI.v1.slaves.child('change', this.data.id).post(this.httpClient, {
      id: this.data.id,
      name: name,
      description: description,
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe(() => {
      this.toasterService.pop('success', undefined, this.i18nService.get('device properties changed'))
      this.data.name = name
      this.data.description = description
      this.matDialogRef.close(true)
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
