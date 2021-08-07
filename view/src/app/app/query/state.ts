import { filter, takeUntil } from "rxjs/operators"
import { ServerAPI } from "src/app/core/core/api"
import { SessionService } from "src/app/core/session/session.service"
import { Closed } from "src/app/core/utils/closed"
import { Completer } from "src/app/core/utils/completer"
import { Session } from "src/app/core/session/session"
import { HttpParams } from "@angular/common/http"
import { interval, Observable } from "rxjs"
enum EventCode {
    Ping = 1,
    Subscribe = 2,
}
interface Request {
    code: EventCode.Subscribe
    targets?: Array<string>
}
interface Response {
    items: Array<{
        id: string
        ready: boolean
    }>
}
export class StateManager {
    constructor(
        private readonly sessionService: SessionService,
    ) {
        this.sessionService.observable.pipe(
            takeUntil(this.closed_.observable),
            filter((session) => session?.access && session?.userdata.id ? true : false)
        ).subscribe((session) => {
            this.session_ = session
        })
    }
    private session_: Session | undefined
    private completer_: Completer<WebSocket> | undefined
    private ws_: WebSocket | undefined
    private targets_: Array<string> = []
    private request_: Array<string> | undefined
    private readys_ = new Set<string>()
    private closed_ = new Closed()
    get isClosed(): boolean {
        return this.closed_.isClosed
    }
    get isNotClosed(): boolean {
        return !this.closed_.isNotClosed
    }
    close() {
        if (this.isNotClosed) {
            this.closed_.close()
            if (this.timeout_) {
                clearTimeout(this.timeout_)
            }
            if (this.ws_) {
                this.ws_.close()
                this.ws_ = undefined
            }
        }
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
            const ws = await this.ws()
            if (!this.request_) {
                return
            }
            ws.send(JSON.stringify({
                code: EventCode.Subscribe,
                targets: this.request_,
            }))
            this.request_ = undefined
        } catch (e) {
            if (this.request_) {
                this._subscribe()
            }
        }
    }
    private tempDelay_ = 0
    private at = 0
    private timeout_: any
    private async ws(): Promise<WebSocket> {
        if (this.completer_) {
            return this.completer_.promise
        }
        let completer = new Completer<WebSocket>()
        this.completer_ = completer

        const tempDelay = this.tempDelay_
        if (tempDelay) {
            const deadline = this.at + tempDelay * 1000
            const now = Date.now()
            if (deadline > now) {
                const duration = deadline - now
                await new Promise<void>((resolve) => {
                    console.warn(`ws connect wait ${duration / 1000}s`)
                    this.timeout_ = setTimeout(() => {
                        resolve()
                    }, duration)
                })
            }
        }
        try {
            const ws = await this._connect()
            this.tempDelay_ = 0
            completer.resolve(ws)
        } catch (e) {
            this.completer_ = undefined

            if (this.tempDelay_) {
                this.tempDelay_ *= 2
                if (this.tempDelay_ > 16) {
                    this.tempDelay_ = 16
                }
            } else {
                this.tempDelay_ = 2
            }
            this.at = Date.now()
            console.log(`ws connect error, retrying in ${this.tempDelay_}s.`)
            completer.reject(e)
            return completer.promise
        }
        return completer.promise
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
    private interval_: any
    private async _connect(): Promise<WebSocket> {
        return new Promise<WebSocket>((resolve, reject) => {
            try {
                const session = this.session_
                if (session?.access && session?.userdata?.id) {
                    const query = new HttpParams({
                        fromObject: {
                            access_token: session.access,
                        }
                    })
                    let url = ServerAPI.v1.slaves.websocketURL('subscribe')
                    console.log(`connect ${url}`)
                    url += + '?' + query.toString()
                    const ws = new WebSocket(url)
                    ws.binaryType = "arraybuffer"
                    let end = false
                    ws.onopen = () => {
                        this.ws_ = ws
                        end = true
                        if (!this.interval_) {
                            this.interval_ = interval(40 * 1000).pipe(
                                takeUntil(this.closed_.observable)
                            ).subscribe(() => {
                                console.log(`send heart`)
                                ws.send(JSON.stringify({
                                    code: EventCode.Ping,
                                }))
                            })
                        }
                        ws.onmessage = (evt) => {
                            if (ws != this.ws_) {
                                return
                            }
                            if (typeof evt.data == "string") {
                                this._onMessage(evt.data)
                            }
                        }
                        resolve(ws)
                    }
                    ws.onclose = (evt) => {
                        if (this.ws_ == ws) {
                            this.ws_ = undefined
                            this.completer_ = undefined
                        }

                        ws.close()
                        if (!end) {
                            end = true
                            reject(new Error(`ws close: ${evt.reason}`))
                        }
                    }
                    ws.onerror = (evt) => {
                        if (this.ws_ == ws) {
                            this.ws_ = undefined
                            this.completer_ = undefined
                        }

                        ws.close()
                        if (!end) {
                            end = true
                            reject(new Error(`ws err`))
                        }
                    }
                } else {
                    reject(new Error(`session nil`))
                }
            } catch (e) {
                reject(e)
            }
        })
    }
    private _onMessage(data: string) {
        const obj: Response = JSON.parse(data)
        if (Array.isArray(obj.items)) {
            const readys = this.readys_
            obj.items.forEach((item) => {
                if (item.ready) {
                    readys.add(item.id)
                } else {
                    readys.delete(item.id)
                }
            })
        }
    }
}