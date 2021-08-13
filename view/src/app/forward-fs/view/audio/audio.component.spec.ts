import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AudioComponent } from './audio.component';

describe('AudioComponent', () => {
  let component: AudioComponent;
  let fixture: ComponentFixture<AudioComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ AudioComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(AudioComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
