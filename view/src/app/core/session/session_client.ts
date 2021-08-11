import { HttpClient, HttpParams } from "@angular/common/http"
import { interval } from "rxjs"
import { filter, first, takeUntil } from "rxjs/operators"
import { Manager, Session } from "src/app/core/session/session"
import { environment } from "src/environments/environment"
import { Codes } from "../core/restful"
import { Client } from "../net/client"
import { SessionService } from "./session.service"
export class Access {
    private session_: Session | undefined
    private refresh_ = false
    setSession(session?: Session) {
        this.session_ = session
        this.refresh_ = false
    }
    setRefresh() {
        this.refresh_ = true
    }
    close() {
        this.session_ = undefined
        this.refresh_ = false
    }
    refresh(httpClient: HttpClient): undefined | Promise<Session | undefined> {
        if (this.refresh_ && this.session_) {
            return Manager.instance.refresh(httpClient, this.session_)
        } else {
            return
        }
    }
}
interface Response {
    code: Codes
    message: string
}
enum EventCode {
    Heart = 1,
}
function sendRequest(ws: WebSocket, evt: EventCode) {
    const msg = JSON.stringify({
        event: evt,
    })
    if (!environment.production) {
        console.log(`ws send: ${msg}`)
    }
    ws.send(msg)
}
export class SessionClient extends Client {
    readonly access = new Access()
    delay = -1
    constructor(
        public readonly baseURL: string,
        public readonly heartInterval: number,
        private readonly httpClient: HttpClient,
        private readonly sessionService: SessionService,) {
        super()
    }
    close() {
        if (this.isNotClosed) {
            this.access.close()
            super.close()
        }
    }
    optDelay(): number {
        return this.delay * 1000
    }
    async optURL(): Promise<string> {
        const access = this.access
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
    }
    optOnClose(ws: WebSocket): void {
        this._onclose()
        this._connect()
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
        } else if (data instanceof ArrayBuffer) {
            if (this._checkFirst(ws, counted)) {
                this._onArrayBuffer(data)
            }
        } else {
            ws.close()
        }
    }
    async _connect() {
        if (this.isClosed) {
            return
        }
        try {
            await this.promise()
        } catch (e) {
            if (this.isClosed) {
                return
            }
            this._connect()
        }
    }
    _onclose(): number {
        let delay = this.delay
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
        this.delay = delay
        return delay
    }
    private timer_: any
    private _checkFirst(ws: WebSocket, counted: number, code?: Codes, message?: string): boolean {
        if (counted) {
            return true
        }
        if (code === undefined || code === Codes.OK) {
            this._onConnected()
            this.delay = 0
            const heartInterval = this.heartInterval
            if (!this.timer_ && heartInterval > 1000) {
                this.timer_ = interval(heartInterval).pipe(
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
            this._onConnectError(code, message)
            if (code == Codes.Unauthenticated) {
                this.access.setRefresh()
            }
            ws.close()
            return false
        }
    }
    _onConnected() {
        console.warn(`connect success`)
    }
    _onConnectError(code?: Codes, message?: string) {
        console.warn(`connect err: ${code} ${message}`)
    }
    _onMessage(resp: Response) { }
    _onArrayBuffer(data: ArrayBuffer) { }
}