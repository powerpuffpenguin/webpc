import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { ToasterService } from 'angular2-toaster';
import { takeUntil } from 'rxjs/operators';
import { SessionService } from 'src/app/core/session/session.service';
import { Session } from 'src/app/core/session/session';
import { Closed } from 'src/app/core/utils/closed';

@Component({
  selector: 'shared-sign-in',
  templateUrl: './sign-in.component.html',
  styleUrls: ['./sign-in.component.scss']
})
export class SignInComponent implements OnInit, OnDestroy {
  constructor(private readonly sessionService: SessionService,
    private readonly matDialogRef: MatDialogRef<SignInComponent>,
    private readonly toasterService: ToasterService,
  ) { }
  disabled = true

  name = ''
  password = ''
  remember = true
  visibility = false
  session: Session | undefined
  private closed_ = new Closed()
  ngOnInit(): void {
    const sessionService = this.sessionService

    sessionService.signining.pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((disabled) => {
      this.disabled = disabled
    })

    sessionService.observable.pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((session) => {
      if (session != this.session) {
        this.session = session
      }
    })
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  onClose() {
    this.matDialogRef.close()
  }
  onSubmit() {
    this.closed_.watchPromise(
      this.sessionService.signin(this.name, this.password, this.remember),
      undefined,
      (e) => {
        this.toasterService.pop('error', undefined, e)
      },
    )
  }
}
