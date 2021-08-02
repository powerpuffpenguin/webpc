import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RequireComponent } from './require.component';

describe('RequireComponent', () => {
  let component: RequireComponent;
  let fixture: ComponentFixture<RequireComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RequireComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(RequireComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
