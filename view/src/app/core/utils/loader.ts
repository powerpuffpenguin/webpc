import { Completer } from "./completer";

export interface State {
    load(): Promise<void>
}

export class Loader {
    private completer_: Completer<void> | undefined

    constructor(public readonly states: Array<State>) {
    }
    load(): Promise<void> {
        if (this.completer_) {
            return this.completer_.promise
        }
        const completer = new Completer<void>()
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