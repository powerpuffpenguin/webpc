import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
interface Dir {
  name: string
  path: string
}
@Component({
  selector: 'fs-path',
  templateUrl: './path.component.html',
  styleUrls: ['./path.component.scss']
})
export class PathComponent implements OnInit {

  @Output() pathChange = new EventEmitter<string>()
  dirs: Array<Dir> = []
  private path_: string = ''
  @Input()
  disabled: boolean = false
  @Input()
  set path(val: string) {
    if (typeof val !== "string") {
      val = ''
    }
    if (val == this.path_) {
      return
    }
    this.path_ = val
    this.val = val
    const strs = this.path_.split('/')
    const dirs = new Array<Dir>()
    let path = ''
    for (let i = 0; i < strs.length; i++) {
      const str = strs[i]
      if (str != "") {
        path += '/' + str
        dirs.push({
          name: str,
          path: path,
        })
      }
    }
    this.dirs = dirs
  }
  get path(): string {
    return this.path_
  }
  edit = false
  val: string = ''
  ngOnInit(): void {
  }
  onClickDone() {
    this.pathChange.emit(this.val)
  }
  onClickDir(node: Dir) {
    this.pathChange.emit(node.path)
  }
}
