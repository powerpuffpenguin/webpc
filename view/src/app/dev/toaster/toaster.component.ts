import { Component, OnInit } from '@angular/core';
import { ToasterService, ToastType } from 'src/app/core/toaster.service';

@Component({
  selector: 'app-toaster',
  templateUrl: './toaster.component.html',
  styleUrls: ['./toaster.component.scss']
})
export class ToasterComponent implements OnInit {
  title = false
  type: ToastType = 'success'
  constructor(private readonly toasterService: ToasterService) { }

  ngOnInit(): void {
  }

  onSubmit() {
    let title: string | undefined
    if (this.title) {
      title = this.type
    }
    this.toasterService.pop(this.type, title, `${this.title}`)

  }
}
