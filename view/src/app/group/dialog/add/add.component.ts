import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'src/app/core/toaster.service';
import { finalize, takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Closed } from 'src/app/core/utils/closed';
import { Element } from '../../../core/group/tree';

interface Response {
  id: string
}
interface InjectData {
  onAdded(data: Element): void
  parent: Element
}
@Component({
  selector: 'app-add',
  templateUrl: './add.component.html',
  styleUrls: ['./add.component.scss']
})
export class AddComponent implements OnInit, OnDestroy {
  disabled = false
  name = ''
  description = ''
  private closed_ = new Closed()
  constructor(@Inject(MAT_DIALOG_DATA) private readonly data_: InjectData,
    private httpClient: HttpClient,
    private toasterService: ToasterService,
    private matDialogRef: MatDialogRef<AddComponent>,
    private i18nService: I18nService,
  ) { }

  ngOnInit(): void {
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  onClose() {
    this.matDialogRef.close()
  }
  get parent(): string {
    return this.data_.parent.name
  }
  onSubmit() {
    if (this.disabled) {
      return
    }
    this.disabled = true
    const name = this.name.trim()
    const description = this.description.trim()
    ServerAPI.v1.groups.post<Response>(this.httpClient, {
      parent: this.data_.parent.id,
      name: name,
      description: description,
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe((response) => {
      this.data_?.onAdded(new Element(response.id, name, description))
      this.toasterService.pop('success', undefined, this.i18nService.get('add group successed'))
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
