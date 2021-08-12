import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'fs-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss']
})
export class NavigationComponent implements OnInit {
  @Input()
  target = ''

  @Input()
  root = ''
  constructor() { }

  ngOnInit(): void {
  }

}
