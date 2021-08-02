import { Component, OnInit, OnDestroy } from '@angular/core';
import { filter, takeUntil } from 'rxjs/operators';
import { SettingsService, Theme } from 'src/app/core/settings/settings.service';
import { Closed } from 'src/app/core/utils/closed';
import { environment } from 'src/environments/environment';
import { MatDialog } from '@angular/material/dialog';
import { SignInComponent } from '../sign-in/sign-in.component';
import { SessionService } from 'src/app/core/session/session.service';
import { Session } from 'src/app/core/session/session';
import { PasswordComponent } from '../password/password.component';
interface Data {
  id: Theme,
  name: string
}
const Themes: Array<Data> = [
  {
    id: 'deeppurple-amber',
    name: 'Deep Purple & Amber'
  },
  {
    id: 'indigo-pink',
    name: 'Indigo & Pink',
  },
  {
    id: 'pink-bluegrey',
    name: 'Pink & Blue-grey',
  },
  {
    id: 'purple-green',
    name: 'Purple & Green',
  },
]
@Component({
  selector: 'shared-navigation-bar',
  templateUrl: './navigation-bar.component.html',
  styleUrls: ['./navigation-bar.component.scss']
})
export class NavigationBarComponent implements OnInit, OnDestroy {
  constructor(private readonly settingsService: SettingsService,
    private readonly matDialog: MatDialog,
    private readonly sessionService: SessionService,
  ) {
  }
  session: Session | undefined
  private closed_ = new Closed()
  themes = Themes
  theme = ''
  get production(): boolean {
    return environment.production
  }
  ngOnInit(): void {
    const settingsService = this.settingsService
    this.theme = settingsService.getTheme()
    settingsService.theme.pipe(
      filter((v) => v != this.theme),
      takeUntil(this.closed_.observable),
    ).subscribe((theme) => {
      this.theme = theme
    })
    this.sessionService.observable.pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((session) => {
      this.session = session
    })
  }
  ngOnDestroy() {
    this.closed_.close()
  }
  trackByThemeId(index: number, item: any): string {
    return item.id
  }
  themeIcon(id: string) {
    return id === this.theme ? 'radio_button_checked' : 'radio_button_unchecked'
  }
  themeColor(id: string) {
    return id === this.theme ? 'accent' : ''
  }
  onClickTheme(id: Theme) {
    this.settingsService.nextTheme(id)
  }
  onClickSignin() {
    this.matDialog.open(SignInComponent)
  }
  onClickSignout() {
    this.sessionService.signout()
  }
  onClickPassword() {
    this.matDialog.open(PasswordComponent, {
      disableClose: true,
    })
  }
}
