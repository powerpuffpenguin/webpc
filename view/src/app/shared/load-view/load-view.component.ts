import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
export interface LoadError {
  id: string
  err: any
}
@Component({
  selector: 'shared-load-view',
  templateUrl: './load-view.component.html',
  styleUrls: ['./load-view.component.scss']
})
export class LoadViewComponent implements OnInit {
  @Input()
  hasErr = false
  @Input()
  errs: Array<LoadError> = []
  @Output()
  valChange = new EventEmitter<Event>()
  constructor() { }
  ngOnInit(): void {
  }
  onClickRefresh(evt: Event) {
    this.valChange.emit(evt)
  }
}
