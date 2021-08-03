import { Component, OnInit, OnDestroy } from '@angular/core';
import { filter, first, takeUntil } from 'rxjs/operators';
import { SessionService } from 'src/app/core/session/session.service';
import { Closed } from 'src/app/core/utils/closed';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit, OnDestroy {
  private closed_ = new Closed()
  session = false
  constructor(public sessionService: SessionService) {
    this.sessionService.observable.pipe(
      filter((s) => s ? true : false),
      first(),
      takeUntil(this.closed_.observable),
    ).subscribe(() => {
      this.session = true
    })
  }

  ngOnInit(): void {
  }
  ngOnDestroy() {
    this.closed_.close()
  }

}
