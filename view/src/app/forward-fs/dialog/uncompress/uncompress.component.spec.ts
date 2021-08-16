import { ComponentFixture, TestBed } from '@angular/core/testing';

import { UncompressComponent } from './uncompress.component';

describe('UncompressComponent', () => {
  let component: UncompressComponent;
  let fixture: ComponentFixture<UncompressComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ UncompressComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(UncompressComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
