import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UpgradedComponent } from './upgraded.component';

describe('UpgradedComponent', () => {
  let component: UpgradedComponent;
  let fixture: ComponentFixture<UpgradedComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ UpgradedComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(UpgradedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
