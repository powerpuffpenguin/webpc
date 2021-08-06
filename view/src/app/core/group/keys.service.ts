import { Injectable } from '@angular/core';
import { Completer } from '../utils/completer';
import { GroupService } from './group.service';
import { NetElement } from './tree';

@Injectable({
  providedIn: 'root'
})
export class KeysService {
  private completer_: Completer<Map<string, NetElement>> | undefined
  constructor(private readonly groupService: GroupService) {
    groupService.resetObservable.subscribe(() => {
      if (this.completer_) {
        this.completer_ = undefined
      }
    })
  }
  get promise(): Promise<Map<string, NetElement>> {
    if (this.completer_) {
      return this.completer_.promise
    }
    const completer = new Completer<Map<string, NetElement>>()
    this.completer_ = completer
    this.groupService.promise.then((items) => {
      const keys = new Map<string, NetElement>()
      items.forEach((v) => {
        keys.set(v.id, v)
      })
      completer.resolve(keys)
    }).catch((e) => {
      completer.reject(e)
      if (completer == this.completer_) {
        this.completer_ = undefined
      }
    })
    return completer.promise
  }
  get(id: string): Promise<NetElement | undefined> {
    return this.promise.then((keys) => {
      return keys.get(id)
    })
  }
  parentName(id: string, name?: string, clear?: boolean): string {
    if (typeof name === "string" && name != "") {
      if (clear) {
        return name
      }
      return `${name} -> ${id}`
    }
    return `${id}`
  }
}
