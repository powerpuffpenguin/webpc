import { Component, OnDestroy, ViewChild, ElementRef, } from '@angular/core';
import { MatIconRegistry } from '@angular/material/icon';
import { ToasterConfig } from 'angular2-toaster';
import { filter, takeUntil } from 'rxjs/operators';
import { I18nService } from './core/i18n/i18n.service';
import { NavigationService } from './core/navigation/navigation.service';
import { SettingsService } from './core/settings/settings.service';
import { Closed } from './core/utils/closed';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnDestroy {
  private closed_ = new Closed()
  theme = ''
  config = new ToasterConfig({
    positionClass: "toast-bottom-right"
  })
  fullscreen = false
  constructor(readonly settingsService: SettingsService,
    readonly navigationService: NavigationService,
    private readonly matIconRegistry: MatIconRegistry,
    private readonly i18nService: I18nService,
  ) {
    // theme
    this.theme = settingsService.getTheme()
    settingsService.theme.pipe(
      filter((v) => v != this.theme),
      takeUntil(this.closed_.observable),
    ).subscribe((theme) => {
      this.theme = theme
    })

    this.navigationService.fullscreenObservable.pipe(
      takeUntil(this.closed_.observable),
    ).subscribe((ok) => {
      this.fullscreen = ok
    })

    // fontawesome
    this.matIconRegistry.registerFontClassAlias(
      'fontawesome-fa',
      'fa'
    ).registerFontClassAlias(
      'fontawesome-fab',
      'fab'
    ).registerFontClassAlias(
      'fontawesome-fal',
      'fal'
    ).registerFontClassAlias(
      'fontawesome-far',
      'far'
    ).registerFontClassAlias(
      'fontawesome-fas',
      'fas'
    )
  }
  @ViewChild("xi18n")
  private xi18nRef: ElementRef | undefined
  ngAfterViewInit() {
    this.i18nService.init(this.xi18nRef?.nativeElement)
  }
  ngOnDestroy() {
    this.closed_.close()
  }
}
