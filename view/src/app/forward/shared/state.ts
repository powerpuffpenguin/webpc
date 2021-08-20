import { HttpClient } from "@angular/common/http"
import { Closed } from "src/app/core/utils/closed"
import { Loader } from "src/app/core/utils/loader"
import { SharedResponse, SharedState } from './load_state'
export interface Error {
    id: string
    err: any
}
const DefaultValue: any = {}
export class State {
    ready = false
    loader: Loader = DefaultValue
    mount: SharedResponse = DefaultValue
    errs: Array<Error> = []
    hasErr = false
    closed = new Closed()
    constructor(readonly httpClient: HttpClient, public readonly target: string) {
        if (target == '') {
            return
        }
        const opts = {
            httpClient: httpClient,
            target: target,
            cancel: this.closed.observable,
        }
        this.loader = new Loader([
            new SharedState(opts, (data) => {
                this.mount = data
                if (Array.isArray(data.name)) {
                    data.name.sort()
                }
            }, (e) => {
                this.errs.push({
                    id: 'SharedState',
                    err: e,
                })
            }),
        ])
    }
    refresh() {
        this.hasErr = false
        this.errs = []
        this.loader.load().then(() => {
            this.ready = true
        }).catch((_) => {
            this.hasErr = true
        })
    }
}