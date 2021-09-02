import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'app-upgraded',
  templateUrl: './upgraded.component.html',
  styleUrls: ['./upgraded.component.scss']
})
export class UpgradedComponent implements OnInit {

  constructor(@Inject(MAT_DIALOG_DATA) public readonly data: string) { }

  ngOnInit(): void {
  }

}
