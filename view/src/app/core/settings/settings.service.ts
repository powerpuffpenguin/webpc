import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { getItem, setItem } from "../utils/local-storage";
export type Theme = 'deeppurple-amber' | 'indigo-pink' | 'pink-bluegrey' | 'purple-green'
const ThemeKey = 'settings.theme'
export const DefaultTheme: Theme = 'deeppurple-amber'
@Injectable({
  providedIn: 'root'
})
export class SettingsService {
  private readonly theme_ = new BehaviorSubject<Theme>(loadTheme())
  constructor() { }
  get theme(): Observable<Theme> {
    return this.theme_
  }
  getTheme(): Theme {
    return this.theme_.value
  }
  nextTheme(theme: Theme) {
    if (theme === this.theme_.value) {
      return
    }
    if (theme !== 'deeppurple-amber' &&
      theme !== 'indigo-pink' &&
      theme !== 'pink-bluegrey' &&
      theme !== 'purple-green') {
      throw new Error(`theme not supported : ${theme}`);
    }
    this.theme_.next(theme)
    setItem(ThemeKey, theme)
  }
}
function loadTheme(): Theme {
  let result = getItem(ThemeKey) as Theme
  if (result !== 'deeppurple-amber' &&
    result !== 'indigo-pink' &&
    result !== 'pink-bluegrey' &&
    result !== 'purple-green') {
    result = DefaultTheme
  }
  return result
}