
import { ServerAPI } from "src/app/core/core/api"
import { SessionService } from "src/app/core/session/session.service"
import { HttpClient } from "@angular/common/http"
import { environment } from "src/environments/environment"
import { Codes } from "src/app/core/core/restful"
import { SessionClient } from "src/app/core/session/session_client"
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
export class StateManager extends SessionClient {
    private targets_: Array<string> = []
    private request_: Array<string> | undefined
    private readys_ = new Set<string>()
    constructor(httpClient: HttpClient,
        sessionService: SessionService,
        private readonly onchanged: (id: string, ready: boolean) => void,
    ) {
        super(ServerAPI.v1.slaves.websocketURL('subscribe'),
            HeartInterval,
            httpClient, sessionService,
        )
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
    _onMessage(ws: WebSocket, resp: Response) {
        if (resp.code == Codes.OK) {
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
        } else {
            console.warn('ws recv: ', resp)
        }
    }
    _onArrayBuffer(ws: WebSocket, data: ArrayBuffer) {
        console.warn('ws not supported ArrayBuffer message')
        ws.close()
    }

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
