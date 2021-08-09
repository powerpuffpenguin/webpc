import { HttpClient } from "@angular/common/http"
import { interval } from "rxjs"
import { takeUntil } from "rxjs/operators"
import { Closed } from "src/app/core/utils/closed"
import { Loader } from "src/app/core/utils/loader"
import { RequireState } from "src/app/core/utils/requirenet"
import { durationString } from "src/app/core/utils/utils"
import { VersionState, VersionResponse, StartAtState, StartAtResponse, DataState, DataResponse } from './load_state'
export interface Error {
    id: string
    err: any
}
const DefaultValue: any = {}
export class State {
    ready = false
    loader: Loader = DefaultValue
    moment: any
    data: DataResponse = DefaultValue
    version: VersionResponse = DefaultValue
    startAt: StartAtResponse = DefaultValue
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
            new RequireState('moment', (moment) => {
                this.moment = moment
            }, (e) => {
                this.errs.push({
                    id: 'moment',
                    err: e,
                })
            }),
            new DataState(opts, (data) => {
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
        ])
    }
    refresh() {
        this.hasErr = false
        this.errs = []
        this.loader.load().then(() => {
            const moment = this.moment
            const startAt = this.startAt
            startAt.at = moment.unix(startAt.result)
            const d = moment.duration(moment.unix(moment.now() / 1000).diff(startAt.at))
            startAt.started = durationString(d)


            interval(1000).pipe(
                takeUntil(this.closed.observable),
            ).subscribe(() => {
                const startAt = this.startAt
                const d = moment.duration(moment.unix(moment.now() / 1000).diff(startAt.at))
                startAt.started = durationString(d)
            })
            this.ready = true
        }).catch((_) => {
            this.hasErr = true
        })
    }
}