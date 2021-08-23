import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
interface Data {
  fontFamily: string
  fontSize: number
  onFontFamily(val: string): void
  onFontSize(val: number): void
}
@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {

  constructor(@Inject(MAT_DIALOG_DATA) public data: Data,) { }
  options = [
    `Lucida Console`, `sans-serif `, `serif`, `monospace`, `cursive`, `fantasy`, `新宋体`
  ]
  ngOnInit(): void {
  }
  onClickFontFamily() {
    this.data.onFontFamily(this.data.fontFamily)
  }
  onClickFontSize() {
    this.data.onFontSize(this.data.fontSize)
  }
}
