import { TestBed } from '@angular/core/testing';

import { FullscreenService } from './fullscreen.service';

describe('FullscreenService', () => {
  let service: FullscreenService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(FullscreenService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
