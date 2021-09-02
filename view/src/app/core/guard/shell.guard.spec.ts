import { TestBed } from '@angular/core/testing';

import { ShellGuard } from './shell.guard';

describe('ShellGuard', () => {
  let guard: ShellGuard;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    guard = TestBed.inject(ShellGuard);
  });

  it('should be created', () => {
    expect(guard).toBeTruthy();
  });
});
