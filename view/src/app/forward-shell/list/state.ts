import { HttpClient } from "@angular/common/http"
import { Closed } from "src/app/core/utils/closed"
import { Loader } from "src/app/core/utils/loader"
import { ListResponse, ListState } from './load_state'
export interface Error {
    id: string
    err: any
}
const DefaultValue: any = {}
export class State {
    ready = false
    loader: Loader = DefaultValue
    list: ListResponse = DefaultValue
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
            new ListState(opts, (data) => {
                if (data && data.result.length) {
                    data.result.sort((l, r) => {
                        if (l.id == r.id) {
                            return 0
                        }
                        return l.id > r.id ? 1 : -1
                    })
                }
                this.list = data
            }, (e) => {
                this.errs.push({
                    id: 'ListState',
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