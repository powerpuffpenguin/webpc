import { TestBed } from '@angular/core/testing';

import { TextGuard } from './text.guard';

describe('TextGuard', () => {
  let guard: TextGuard;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    guard = TestBed.inject(TextGuard);
  });

  it('should be created', () => {
    expect(guard).toBeTruthy();
  });
});
