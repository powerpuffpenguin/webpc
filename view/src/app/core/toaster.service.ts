import { Injectable } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
export type ToastType = 'success' | 'error' | 'warning' | 'info'
@Injectable({
  providedIn: 'root'
})
export class ToasterService {

  constructor(private readonly toastrService: ToastrService) { }
  pop(type: ToastType, title?: string, message?: any) {
    switch (type) {
      case "success":
        this.toastrService.success(message, title)
        break
      case "error":
        this.toastrService.error(message, title)
        break
      case "warning":
        this.toastrService.warning(message, title)
        break
      case "info":
        this.toastrService.info(message, title)
        break
      default:
        throw new Error(`unknow toaster type: ${type}`)
    }
  }
}
