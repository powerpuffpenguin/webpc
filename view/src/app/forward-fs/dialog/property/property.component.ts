import { Component, OnInit, Inject } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { FileInfo } from '../../fs';
@Component({
  selector: 'app-property',
  templateUrl: './property.component.html',
  styleUrls: ['./property.component.scss']
})
export class PropertyComponent implements OnInit {

  constructor(@Inject(MAT_DIALOG_DATA) public source: Array<FileInfo>,
  ) {
  }
  ngOnInit(): void {
  }

}
