import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Completer } from '../utils/completer';
import { NetElement } from './tree';

@Injectable({
  providedIn: 'root'
})
export class GroupService {
  private completer_: Completer<Array<NetElement>> | undefined
  constructor(private readonly httpClient: HttpClient) { }

  get promise(): Promise<Array<NetElement>> | undefined {
    if (this.completer_) {
      return this.completer_.promise
    }
    return
  }
}
