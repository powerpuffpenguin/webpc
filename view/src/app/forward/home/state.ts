import { HttpClient } from "@angular/common/http"
import { interval, timer } from "rxjs"
import { takeUntil } from "rxjs/operators"
import { Closed } from "src/app/core/utils/closed"
import { datetimeString, usedString } from "src/app/core/utils/datetime"
import { Loader } from "src/app/core/utils/loader"
import { durationString } from "src/app/core/utils/utils"
import { VersionState, VersionResponse, StartAtState, StartAtResponse, DataState, DataResponse, UpgradedResponse, UpgradedState } from './load_state'
export interface Error {
    id: string
    err: any
}
const DefaultValue: any = {}
export class State {
    ready = false
    loader: Loader = DefaultValue
    data: DataResponse = DefaultValue
    version: VersionResponse = DefaultValue
    startAt: StartAtResponse = DefaultValue
    upgraded: UpgradedResponse = DefaultValue
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
            new DataState(opts, (data) => {
                console.log(data)
                this.data = data
            }, (e) => {
                this.errs.push({
                    id: 'data',
                    err: e,
                })
            }),
            new VersionState(opts, (data) => {
                this.version = data
            }, (e) => {
                this.errs.push({
                    id: 'VersionState',
                    err: e,
                })
            }),
            new StartAtState(opts, (data) => {
                this.startAt = data
            }, (e) => {
                this.errs.push({
                    id: 'StartAtState',
                    err: e,
                })
            }),
            new UpgradedState(opts, (data) => {
                this.upgraded = data
            }, (e) => {
                this.errs.push({
                    id: 'UpgradedState',
                    err: e,
                })
            }),
        ])
    }
    refresh() {
        this.hasErr = false
        this.errs = []
        this.loader.load().then(() => {
            const response = this.startAt
            const startAt = typeof response.result === "number" ? response.result : parseInt(response.result)
            if (Number.isSafeInteger(startAt)) {
                response.at = datetimeString(new Date(startAt * 1000))
                timer(0, 1000).pipe(takeUntil(this.closed.observable)).subscribe({
                    next: () => {
                        let used = Math.floor(Date.now() / 1000)
                        if (used > startAt) {
                            used -= startAt
                        } else {
                            used = 0
                        }
                        response.started = usedString(used)
                    },
                })
            }
            this.ready = true
        }).catch((_) => {
            this.hasErr = true
        })
    }
}