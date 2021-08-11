import { Manager, Session } from "src/app/core/session/session"
import { filter, first, takeUntil } from "rxjs/operators"
import { SessionService } from "src/app/core/session/session.service"
import { Closed } from "src/app/core/utils/closed"
import { HttpClient, HttpParams } from "@angular/common/http"
import { ServerAPI } from "src/app/core/core/api"
import { Client, ClientOption } from "src/app/core/net/client_w"
import { environment } from "src/environments/environment"
import { interval } from "rxjs"
import { Codes } from "src/app/core/core/restful"
interface Response {
    code: Codes
    message: string
}
const HeartInterval = 40 * 1000
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


interface Writer {
    writeln(text: string): void
    write(text: string): void
}
export class Listener extends Client implements ClientOption {
    baseURL = ServerAPI.v1.logger.websocketURL('attach')
    private session_: Session | undefined
    constructor(private readonly httpClient: HttpClient,
        private readonly writer: Writer,
        private readonly sessionService: SessionService,
    ) {
        super()
        this._connect()
    }
    private first_ = true
    private delay_ = -1
    optDelay(): number {
        return this.delay_ * 1000
    }
    async optURL(): Promise<string> {
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
        const access = session?.access ?? ''
        const baseURL = this.baseURL
        const url = baseURL + '?' + new HttpParams({
            fromObject: {
                access_token: access,
            }
        }).toString()
        this.session_ = session
        if (environment.production) {
            console.log(`connect ${baseURL}`)
        } else {
            console.log(`connect ${url}`)
        }
        return url
    }
    optOnNew(ws: WebSocket): void {
        ws.binaryType = 'arraybuffer'
        this.writer.writeln(`attach logger console`)
    }
    optOnOpenError(_: WebSocket): void {
        this._onclose()
    }
    optOnClose(ws: WebSocket): void {
        this._onclose()
        this._connect()
    }
    optOnMessage(ws: WebSocket, evt: MessageEvent): void {
        const data = evt.data
        if (typeof data === "string") {
            const resp: Response = JSON.parse(data)
            if (resp.code === undefined) {
                resp.code == Codes.OK
            }
            if (this._checkFirst(ws, resp.code, resp.message)) {
                this._onMessage(resp)
            }
        } else if (data instanceof ArrayBuffer) {
            if (this._checkFirst(ws)) {
                this._onArrayBuffer(data)
            }
        } else {
            ws.close()
        }
    }
    private async _connect() {
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
    private _onclose() {
        this.session_ = undefined
        this.first_ = true

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
        if (delay <= 0) {
            console.warn(`websocket err, retrying in 0s`)
        } else {
            console.warn(`websocket err, retrying in ${delay}s`)
        }
    }
    private timer_: any
    private _checkFirst(ws: WebSocket, code?: Codes, message?: string): boolean {
        if (!this.first_) {
            return true
        }
        this.first_ = false
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
            if (code == Codes.Unauthenticated && this.session_) {
                // Manager.instance.refresh(httpClient, readySession)
            }
            ws.close()
            return false
        }
    }
    private _onMessage(resp: Response) {
        if (resp.code == Codes.OK) {
            console.log('ws recv: ', resp)
        } else {
            console.warn('ws recv: ', resp)
        }
    }
    private _onArrayBuffer(data: ArrayBuffer) {
        const enc = new TextDecoder("utf-8")
        let str = enc.decode(data)
        str = str.replace(/\n/g, "\r\n")
        this.writer.write(str)
    }
}
