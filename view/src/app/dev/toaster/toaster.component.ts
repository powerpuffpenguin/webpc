import { Component, OnInit } from '@angular/core';
import { ToasterService, ToastType } from 'angular2-toaster';
@Component({
  selector: 'app-toaster',
  templateUrl: './toaster.component.html',
  styleUrls: ['./toaster.component.scss']
})
export class ToasterComponent implements OnInit {
  title = false
  async = false
  type: ToastType = 'success'
  constructor(private readonly toasterService: ToasterService) { }

  ngOnInit(): void {
  }

  onSubmit() {
    let title: string | undefined
    if (this.title) {
      title = this.type
    }
    if (this.async) {
      this.toasterService.pop(this.type, title, `${this.title}`)
    } else {
      this.toasterService.popAsync(this.type, title, `${this.title}`)
    }
  }
}
