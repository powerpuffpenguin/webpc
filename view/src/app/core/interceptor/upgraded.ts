import { BehaviorSubject, Observable } from "rxjs"

export class Upgraded {
    static readonly instance = new Upgraded()
    private constructor() {
    }
    private version_ = new BehaviorSubject<string>('')
    get versionObservable(): Observable<string> {
        return this.version_
    }
    nextVersion(version: string) {
        if (this.version_.value != version) {
            this.version_.next(version)
        }
    }
}