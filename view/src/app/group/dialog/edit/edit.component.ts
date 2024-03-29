import { HttpClient } from '@angular/common/http';
import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ToasterService } from 'src/app/core/toaster.service';
import { finalize, takeUntil } from 'rxjs/operators';
import { ServerAPI } from 'src/app/core/core/api';
import { GroupService } from 'src/app/core/group/group.service';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Closed } from 'src/app/core/utils/closed';
import { Element } from '../../../core/group/tree';
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
  constructor(@Inject(MAT_DIALOG_DATA) public readonly data: Element,
    private readonly httpClient: HttpClient,
    private readonly toasterService: ToasterService,
    private readonly matDialogRef: MatDialogRef<EditComponent>,
    private readonly i18nService: I18nService,
    private readonly groupService: GroupService,
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
    ServerAPI.v1.groups.child('change', this.data.id).post(this.httpClient, {
      name: name,
      description: description,
    }).pipe(
      takeUntil(this.closed_.observable),
      finalize(() => {
        this.disabled = false
      })
    ).subscribe(() => {
      this.toasterService.pop('success', undefined, this.i18nService.get('group properties changed'))
      let changed = false
      if (this.data.name != name) {
        this.data.name = name
        changed = true
      } else {
        this.groupService.reset()
      }
      this.data.description = description
      this.matDialogRef.close(changed)
    }, (e) => {
      this.toasterService.pop('error', undefined, e)
    })
  }
}
