import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ExistsChoiceComponent } from './exists-choice.component';

describe('ExistsChoiceComponent', () => {
  let component: ExistsChoiceComponent;
  let fixture: ComponentFixture<ExistsChoiceComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ExistsChoiceComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ExistsChoiceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
