import { Observable } from "rxjs";
import { first, takeUntil } from "rxjs/operators";
import { Completer } from "./completer";

export interface State {
    load(): Promise<any>
}
export class BasicState<T> implements State {
    private completer_: Completer<T> | undefined
    constructor(
        private readonly loadf: () => Promise<T>,
        private readonly onReady?: (data: T) => void,
        private readonly onError?: (evt: any) => void,
    ) {
    }
    load(): Promise<T> {
        if (this.completer_) {
            return this.completer_.promise
        }
        const completer = new Completer<T>()
        this.completer_ = completer
        this.loadf().then((data) => {
            completer.resolve(data)
            if (this.onReady) {
                this.onReady(data)
            }
        }).catch((e) => {
            this.completer_ = undefined
            completer.reject(e)
            if (this.onError) {
                this.onError(e)
            }
        })
        return completer.promise
    }
}
export function firstPromise<T>(observable: Observable<T>, cancel?: Observable<any>): Promise<T> {
    return new Promise<T>((resolve, reject) => {
        if (cancel) {
            let canceled = false
            observable.pipe(
                first(),
                takeUntil(cancel.pipe((v) => {
                    canceled = true
                    return v
                })),
            ).toPromise().then((data) => {
                if (canceled && data === undefined) {
                    reject(new Error('canceled'))
                } else {
                    resolve(data)
                }
            }).catch((e) => {
                reject(e)
            })
        } else {
            observable.pipe(
                first(),
            ).subscribe((data) => {
                resolve(data)
            }, (e) => {
                reject(e)
            })
        }
    })
}

export class Loader {
    private completer_: Completer<any> | undefined

    constructor(public readonly states: Array<State>) {
    }
    load(): Promise<any> {
        if (this.completer_) {
            return this.completer_.promise
        }
        const completer = new Completer<any>()
        this.completer_ = completer
        const states = this.states
        if (this.states.length == 0) {
            completer.resolve()
        } else {
            let completed = 0
            let failed = false
            let err: any
            for (let i = 0; i < states.length; i++) {
                const state = states[i]
                state.load().catch((e) => {
                    failed = true
                    if (!err) {
                        err = e
                    }
                }).finally(() => {
                    completed++
                    if (states.length == completed) {
                        if (failed) {
                            completer.reject(err)
                            this.completer_ = undefined
                        } else {
                            completer.resolve()
                        }
                    }
                })
            }
        }
        return completer.promise
    }
}