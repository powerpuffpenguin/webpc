import { TestBed } from '@angular/core/testing';

import { FsGuard } from './fs.guard';

describe('FsGuard', () => {
  let guard: FsGuard;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    guard = TestBed.inject(FsGuard);
  });

  it('should be created', () => {
    expect(guard).toBeTruthy();
  });
});
