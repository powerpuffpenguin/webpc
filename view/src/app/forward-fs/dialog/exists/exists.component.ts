import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

@Component({
  selector: 'app-exists',
  templateUrl: './exists.component.html',
  styleUrls: ['./exists.component.scss']
})
export class ExistsComponent implements OnInit {

  constructor(
    private matDialogRef: MatDialogRef<ExistsComponent>,
    @Inject(MAT_DIALOG_DATA) public filename: string,) { }

  ngOnInit(): void {
  }
  onSure() {
    this.matDialogRef.close(true)
  }
  onClose() {
    this.matDialogRef.close()
  }
}
