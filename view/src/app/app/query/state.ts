import { takeUntil } from "rxjs/operators"
import { ServerAPI } from "src/app/core/core/api"
import { SessionService } from "src/app/core/session/session.service"
import { Closed } from "src/app/core/utils/closed"
import { HttpClient, HttpParams } from "@angular/common/http"
import { interval } from "rxjs"
import { environment } from "src/environments/environment"
import { Client } from "src/app/core/net/client"
import { Manager, Session } from "src/app/core/session/session"

enum EventCode {
    Ping = 1,
    Subscribe = 2,
}
interface Request {
    event: EventCode.Subscribe
    targets?: Array<string>
}
interface Response {
    code: number
    emsg: string
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

export class StateManager {
    private closed_ = new Closed()
    private client_ = {} as Client
    constructor(
        httpClient: HttpClient,
        private readonly sessionService: SessionService,
        private readonly onchanged: (id: string, ready: boolean) => void,
    ) {
        const closed = this.closed_
        const baseURL = ServerAPI.v1.slaves.websocketURL('subscribe')
        let session: Session
        this.sessionService.observable.pipe(
            takeUntil(closed.observable),
        ).subscribe((data) => {
            if (data && data.userdata && data.userdata.id && data.access) {
                session = data
            }
        })
        let delay = -1
        let first = true
        let timer: any
        const ctx = this
        let readySession: Session | undefined
        this.client_ = new Client({
            get url(): string {
                const url = baseURL + '?' + new HttpParams({
                    fromObject: {
                        access_token: session.access,
                    }
                }).toString()
                readySession = session
                if (environment.production) {
                    console.log(`connect ${baseURL}`)
                } else {
                    console.log(`connect ${url}`)
                }
                return url
            },
            get delay(): number {
                return delay * 1000
            },
            onNew(ws: WebSocket) {
                ws.binaryType = 'arraybuffer'
            },
            onMessage(ws: WebSocket, data) {
                if (typeof data === "string") {
                    const resp: Response = JSON.parse(data)
                    if (resp.code === undefined) {
                        resp.code == 0
                    }
                    if (first) {
                        if (resp.code == 0) {
                            delay = 0
                            if (!timer) {
                                timer = interval(40 * 1000).pipe(
                                    takeUntil(closed.observable)
                                ).subscribe(() => {
                                    ctx.client_.promise().then((ws) => {
                                        console.log('send heart')
                                        sendRequest(ws, EventCode.Ping)
                                    }).catch((e) => {
                                        console.log('send heart error: ', e)
                                    })
                                })
                            }
                        } else {
                            console.warn('connect err: ', resp)
                            if (resp.code == 16 && readySession) {
                                // Manager.instance.refresh(httpClient, readySession)
                            }
                            ws.close()
                            return
                        }
                    }
                    ctx._onMessage(resp)
                } else {
                    ws.close()
                }
            },
            onClose(_: WebSocket) {
                readySession = undefined
                first = true
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

                if (ctx.targets_.length != 0) {
                    ctx.request_ = ctx.targets_
                    ctx._subscribe()
                }
            }
        })
    }

    private targets_: Array<string> = []
    private request_: Array<string> | undefined
    private readys_ = new Set<string>()
    get isClosed(): boolean {
        return this.closed_.isClosed
    }
    get isNotClosed(): boolean {
        return this.closed_.isNotClosed
    }
    close() {
        if (this.isNotClosed) {
            this.closed_.close()
            this.client_.close()
            this.request_ = undefined
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
            const ws = await this.client_.promise()
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
}