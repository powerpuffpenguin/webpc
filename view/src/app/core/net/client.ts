import { Observable } from "rxjs";
import { Closed } from "../utils/closed";
import { Completer } from "../utils/completer";
enum State {
    Connect,
    Opened,
    Closed,
}
export interface ClientOption {
    url: string
    delay?: number
    onNew?: (ws: WebSocket) => void
    onOpen?: (ws: WebSocket) => void
    onOpenErr?: (ws: WebSocket) => void
    onClose?: (ws: WebSocket) => void
    onMessage?: (ws: WebSocket, data: any) => void
}
export interface Response {
    code: number
    emsg: string
    data: any
}
export class Client {
    private closed_ = new Closed()
    private completer_: Completer<WebSocket> | undefined
    private ws_: WebSocket | undefined
    private at_ = Date.now()
    private timeout_: any
    constructor(public readonly opts: ClientOption) { }
    close() {
        if (this.isClosed) {
            return
        }
        this.closed_.close()
        const timeout = this.timeout_
        if (timeout) {
            this.timeout_ = undefined
            clearTimeout(timeout)
        }
        if (this.ws_) {
            this.ws_.close()
            this.ws_ = undefined
        }
    }
    get isNotClosed(): boolean {
        return this.closed_.isNotClosed
    }
    get isClosed(): boolean {
        return this.closed_.isClosed
    }

    write(data: string | ArrayBufferLike | Blob | ArrayBufferView): Promise<void> {
        return this.promise().then((ws) => ws.send(data))
    }
    ws(): WebSocket | undefined {
        return this.ws_
    }
    async promise(): Promise<WebSocket> {
        if (this.completer_) {
            return this.completer_.promise
        } else if (this.isClosed) {
            throw new Error('client already closed')
        }
        const completer = new Completer<WebSocket>()
        this.completer_ = completer
        try {
            // delay
            const delay = this.opts.delay
            if (typeof delay === "number" && delay > 0) {
                const deadline = this.at_ + delay
                const now = Date.now()
                if (deadline > now) {
                    const duration = deadline - now
                    await new Promise<void>((resolve) => {
                        console.warn(`client connect wait ${duration / 1000}s`)
                        this.timeout_ = setTimeout(() => {
                            resolve()
                        }, duration)
                    })
                    if (this.isClosed) {
                        throw new Error('client already closed')
                    }
                }
            }

            // connect
            const ws = await this._connect(completer)
            if (this.isClosed) {
                ws.close()
                throw new Error('client already closed')
            }
            this.at_ = Date.now()
            this.ws_ = ws
            completer.resolve(ws)
        } catch (e) {
            this.at_ = Date.now()
            this.completer_ = undefined
            completer.reject(e)
        }
        return completer.promise
    }
    private _connect(completer: Completer<WebSocket>): Promise<WebSocket> {
        return new Promise<WebSocket>((resolve, reject) => {
            try {
                const opts = this.opts
                let state = State.Connect
                const url = this.opts.url
                const ws = new WebSocket(url)
                if (opts.onNew) {
                    opts.onNew(ws)
                }
                ws.onopen = (evt) => {
                    if (State.Connect == state) {
                        state = State.Opened
                        resolve(ws)
                    }
                    if (opts.onOpen) {
                        opts.onOpen(ws)
                    }

                    const onMessage = opts.onMessage
                    if (onMessage) {
                        ws.onmessage = (evt) => {
                            if (State.Opened == state) {
                                onMessage(ws, evt.data)
                            }
                        }
                    }
                }
                ws.onclose = (evt) => {
                    ws.close()

                    if (State.Connect == state) {
                        state = State.Closed
                        reject(new Error('websocket close'))
                    } else if (State.Opened == state) {
                        state = State.Closed
                        if (this.completer_ == completer) {
                            this.completer_ = undefined
                        }
                        if (opts.onClose) {
                            opts.onClose(ws)
                        }
                    }
                }
                ws.onerror = (evt) => {
                    ws.close()

                    if (State.Connect == state) {
                        state = State.Closed
                        reject(new Error('websocket error'))
                        if (opts.onOpenErr) {
                            opts.onOpenErr(ws)
                        }
                    } else if (State.Opened == state) {
                        state = State.Closed
                        if (this.completer_ == completer) {
                            this.completer_ = undefined
                        }
                        if (opts.onClose) {
                            opts.onClose(ws)
                        }
                    }
                }
            } catch (e) {
                reject(e)
            }
        })
    }
}