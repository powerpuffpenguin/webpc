import { filter, first, takeUntil } from "rxjs/operators"
import { ServerAPI } from "src/app/core/core/api"
import { SessionService } from "src/app/core/session/session.service"
import { HttpClient, HttpParams } from "@angular/common/http"
import { interval } from "rxjs"
import { environment } from "src/environments/environment"
import { Client, ClientOption } from "src/app/core/net/client"
import { Access } from "src/app/core/net/access"
import { Codes } from "src/app/core/core/restful"
const HeartInterval = 40 * 1000
enum EventCode {
    Heart = 1,
    Subscribe = 2,
}
interface Response {
    code: Codes
    message: string
    items: Array<{
        id: string
        ready: boolean
    }>
}
function sendRequest(ws: WebSocket, evt: EventCode, targets?: Array<string>) {
    let msg: string
    if (targets) {
        msg = JSON.stringify({
            event: evt,
            targets: targets,
        })
    } else {
        msg = JSON.stringify({
            event: evt,
        })
    }
    if (!environment.production) {
        console.log(`ws send: ${msg}`)
    }
    ws.send(msg)
}

export class StateManager extends Client implements ClientOption {
    baseURL = ServerAPI.v1.slaves.websocketURL('subscribe')
    private access_ = new Access()
    constructor(private readonly httpClient: HttpClient,
        private readonly sessionService: SessionService,
        private readonly onchanged: (id: string, ready: boolean) => void,
    ) {
        super()
    }
    close() {
        if (this.isNotClosed) {
            this.access_.close()
            super.close()
            this.request_ = undefined
        }
    }
    private delay_ = -1
    optDelay(): number {
        return this.delay_ * 1000
    }
    async optURL(): Promise<string> {
        const access = this.access_
        const refresh = access.refresh(this.httpClient)
        if (refresh) {
            await refresh
        }

        const session = await this.sessionService.observable.pipe(
            filter((data) => {
                if (data && data.userdata && data.userdata.id && data.access) {
                    return true
                }
                return false
            }),
            first(),
            takeUntil(this.observable),
        ).toPromise()
        const token = session?.access ?? ''
        const baseURL = this.baseURL
        const url = baseURL + '?' + new HttpParams({
            fromObject: {
                access_token: token,
            }
        }).toString()
        access.setSession(session)
        if (environment.production) {
            console.log(`connect ${baseURL}`)
        } else {
            console.log(`connect ${url}`)
        }
        return url
    }
    optOnNew(ws: WebSocket): void {
        ws.binaryType = 'arraybuffer'
    }
    optOnOpenError(_: WebSocket): void {
        const delay = this._onclose()
        if (delay <= 0) {
            console.log(`connect err, retrying in 0s`)
        } else {
            console.log(`connect err, retrying in ${delay}s`)
        }
        if (this.targets_.length != 0) {
            this.request_ = this.targets_
            this._subscribe()
        }
    }
    optOnClose(ws: WebSocket): void {
        this._onclose()
        if (this.targets_.length != 0) {
            this.request_ = this.targets_
            this._subscribe()
        }
    }
    optOnMessage(ws: WebSocket, counted: number, evt: MessageEvent): void {
        const data = evt.data
        if (typeof data === "string") {
            const resp: Response = JSON.parse(data)
            if (resp.code === undefined) {
                resp.code == Codes.OK
            }
            if (this._checkFirst(ws, counted, resp.code, resp.message)) {
                this._onMessage(resp)
            }
        } else {
            ws.close()
        }
    }
    private timer_: any
    private _checkFirst(ws: WebSocket, counted: number, code?: Codes, message?: string): boolean {
        if (counted) {
            return true
        }
        if (code === undefined || code === Codes.OK) {
            this.delay_ = 0
            if (!this.timer_ && HeartInterval > 1000) {
                this.timer_ = interval(HeartInterval).pipe(
                    takeUntil(this.observable)
                ).subscribe(() => {
                    const ws = this.ws()
                    if (ws) {
                        console.log('send heart')
                        sendRequest(ws, EventCode.Heart)
                    }
                })
            }
            return true
        } else {
            console.warn(`connect err: ${code} ${message}`)
            if (code == Codes.Unauthenticated) {
                this.access_.setRefresh()
            }
            ws.close()
            return false
        }
    }
    private _onclose(): number {
        let delay = this.delay_
        if (delay == 0) {
            delay = -1
        } else if (delay < 0) {
            delay = 2
        } else {
            delay *= 2
            if (delay > 16) {
                delay = 16
            }
        }
        this.delay_ = delay
        return delay
    }
    private _onMessage(resp: Response) {
        if (!environment.production) {
            console.log('ws recv:', resp)
        }
        if (Array.isArray(resp.items)) {
            const readys = this.readys_
            const onchanged = this.onchanged
            resp.items.forEach((item) => {
                if (onchanged) {
                    onchanged(item.id, item.ready)
                }
                if (item.ready) {
                    readys.add(item.id)
                } else {
                    readys.delete(item.id)
                }
            })
        }
    }

    private targets_: Array<string> = []
    private request_: Array<string> | undefined
    private readys_ = new Set<string>()

    isReady(id: string): boolean {
        return this.readys_.has(id)
    }
    subscribe(targets: Array<string>) {
        if (this.isClosed) {
            return
        }
        // format
        const set = new Set<string>()
        const strs: Array<string> = []
        const readys = this.readys_
        targets.forEach((id) => {
            if (set.has(id)) {
                return
            }
            set.add(id)

            strs.push(id)
            readys.delete(id)
        })
        strs.sort()
        if (this._isEqual(this.targets_, strs)) {
            return
        }
        // changed
        this.targets_ = strs

        this.request_ = strs
        this._subscribe()
    }
    private async _subscribe() {
        try {
            const ws = await this.promise()
            if (this.isClosed) {
                return
            }
            if (!this.request_) {
                return
            }
            sendRequest(ws, EventCode.Subscribe, this.request_)
            this.request_ = undefined
        } catch (e) {
            if (this.isClosed || !this.request_) {
                return
            }
            this._subscribe()
        }
    }

    private _isEqual(l: Array<string>, r: Array<string>): boolean {
        if (l.length != r.length) {
            return false
        }
        for (let i = 0; i < l.length; i++) {
            if (l[i] != r[i]) {
                return false
            }
        }
        return true
    }
}
