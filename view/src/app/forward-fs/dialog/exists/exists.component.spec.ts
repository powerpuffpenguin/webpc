import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ExistsComponent } from './exists.component';

describe('ExistsComponent', () => {
  let component: ExistsComponent;
  let fixture: ComponentFixture<ExistsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ExistsComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ExistsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
