import { Component, OnInit, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { EventCode } from '../event';

@Component({
  selector: 'app-exists-choice',
  templateUrl: './exists-choice.component.html',
  styleUrls: ['./exists-choice.component.scss']
})
export class ExistsChoiceComponent implements OnInit {
  constructor(
    private matDialogRef: MatDialogRef<ExistsChoiceComponent>,
    @Inject(MAT_DIALOG_DATA) public filename: string) { }

  ngOnInit(): void {
  }
  onYes() {
    this.matDialogRef.close(EventCode.Yes)
  }
  onYesAll() {
    this.matDialogRef.close(EventCode.YesAll)
  }
  onNo() {
    this.matDialogRef.close(EventCode.No)
  }
  onSkip() {
    this.matDialogRef.close(EventCode.Skip)
  }
  onSkipAll() {
    this.matDialogRef.close(EventCode.SkipAll)
  }
}
