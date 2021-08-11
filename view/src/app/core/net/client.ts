import { Observable, timer } from "rxjs";
import { takeUntil } from "rxjs/operators";
import { Closed } from "../utils/closed";
import { Completer } from "../utils/completer";
import { isPromise } from "../utils/promise";
enum State {
    Connect,
    Opened,
    Closed,
}

export interface Response {
    code: number
    emsg: string
    data: any
}
export interface ClientOption {
    /**
     * 返回延遲多久撥號
     */
    optDelay(): number | Promise<number>
    /**
     * 返回連接的服務器地址
     */
    optURL(): string | Promise<string>
    /**
     * 當 new WebSocket 後 被調用
     * @param ws 
     */
    optOnNew(ws: WebSocket): void
    /**
     * 當一個 WebSocket 的 onerror 在 onopen 前回調時 在 onerror 之後調用
     * @param ws 
     */
    optOnOpenError(ws: WebSocket): void
    /**
     * 當一個 WebSocket 的 onopen 回調後 被調用
     * @param ws 
     */
    optOnOpen(ws: WebSocket, evt: Event): void
    /**
    * 當一個 WebSocket 的 onclose/onerror 回調後 被調用 (同個 Websocket 只會回調一次 optOnClose)
    * @param ws 
    */
    optOnClose(ws: WebSocket): void
    /**
     * 當一個 WebSocket 的 onopen 回調後 被調用
     * @param ws 
     * @param counted 已經爲這個 WebSocket 執行了多少次回調
     * @param evt 
     */
    optOnMessage(ws: WebSocket, counted: number, evt: MessageEvent): void
}

class ClientImplement {
    private closed_ = new Closed()
    private ws_: WebSocket | undefined // 緩存讀取就緒的 連接
    private completer_: Completer<WebSocket> | undefined // 保證同時只有一個連接
    private at_ = Date.now() // 最後一次連接時間 用於延遲連接
    get observable(): Observable<boolean> {
        return this.closed_.observable
    }
    get isNotClosed(): boolean {
        return this.closed_.isNotClosed
    }
    get isClosed(): boolean {
        return this.closed_.isClosed
    }
    close() {
        if (this.isClosed) {
            return
        }
        this.closed_.close()
        const ws = this.ws_
        if (ws) {
            ws.close()
            this.ws_ = undefined
        }
    }
    /**
     * @returns 返回最後緩存的就緒連接
     */
    ws(): WebSocket | undefined {
        return this.ws_
    }
    /**
     * @returns 返回就緒的 WebSocket 如果沒有就緒的連接 嘗試創建一個新連接
     */
    async promise(opts: ClientOption): Promise<WebSocket> {
        if (this.completer_) {
            return this.completer_.promise
        } else if (this.isClosed) {
            throw new Error('client already closed')
        }
        const completer = new Completer<WebSocket>()
        this.completer_ = completer
        try {
            // delay
            const opt = opts.optDelay()
            let delay = 0
            if (isPromise(opt)) {
                delay = await opt
            } else if (typeof opt === "number") {
                delay = opt
            }
            if (typeof delay === "number" && delay > 0) {
                const deadline = this.at_ + delay
                const now = Date.now()
                if (deadline > now) {
                    const duration = deadline - now
                    await timer(duration).pipe(
                        takeUntil(this.closed_.observable)
                    ).toPromise()
                    if (this.isClosed) {
                        throw new Error('client already closed')
                    }
                }
            }

            // connect
            const ws = await this._connect(opts, completer)
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
    private _connect(opts: ClientOption, completer: Completer<WebSocket>): Promise<WebSocket> {
        const opt = opts.optURL()
        if (isPromise(opt)) {
            return new Promise<WebSocket>((resolve, reject) => {
                opt.then((url) => {
                    this._connectURL(opts, completer, url, resolve, reject)
                }).catch((e) => {

                    reject(e)
                })
            })
        } else {
            return new Promise<WebSocket>((resolve, reject) => {
                this._connectURL(opts, completer, opt, resolve, reject)
            })
        }
    }
    private _connectURL(opts: ClientOption, completer: Completer<WebSocket>, url: string, resolve: any, reject: any) {
        try {
            let state = State.Connect
            const ws = new WebSocket(url)
            opts.optOnNew(ws)

            ws.onopen = (evt) => {
                if (State.Connect == state) {
                    state = State.Opened
                    resolve(ws)
                }
                opts.optOnOpen(ws, evt)
                let counted = 0
                ws.onmessage = (evt) => {
                    if (State.Opened == state) {
                        opts.optOnMessage(ws, counted, evt)
                        counted++
                    }
                }
            }
            ws.onclose = (evt) => {
                ws.close()
                if (this.ws_ == ws) {
                    this.ws_ = undefined
                }

                if (State.Connect == state) {
                    state = State.Closed
                    reject(new Error('websocket close'))
                    opts.optOnClose(ws)
                } else if (State.Opened == state) {
                    state = State.Closed
                    if (this.completer_ == completer) {
                        this.completer_ = undefined
                    }
                    opts.optOnClose(ws)
                }
            }
            ws.onerror = (evt) => {
                ws.close()
                if (this.ws_ == ws) {
                    this.ws_ = undefined
                }

                if (State.Connect == state) {
                    state = State.Closed
                    reject(new Error('websocket error'))
                    opts.optOnOpenError(ws)
                } else if (State.Opened == state) {
                    state = State.Closed
                    if (this.completer_ == completer) {
                        this.completer_ = undefined
                    }
                    opts.optOnClose(ws)
                }
            }
        } catch (e) {
            reject(e)
        }
    }
}
/**
 * 包裝的 websocket 在連接斷開後 自動恢復連接
*/
export class Client implements ClientOption {
    optDelay(): number | Promise<number> {
        return 0
    }
    optURL(): string | Promise<string> {
        return ''
    }
    optOnNew(ws: WebSocket): void { }
    optOnOpenError(ws: WebSocket): void { }
    optOnOpen(ws: WebSocket, evt: Event): void { }
    optOnClose(ws: WebSocket): void { }
    optOnMessage(ws: WebSocket, counted: number, evt: MessageEvent): void { }

    constructor() { }
    private _clientImplement = new ClientImplement()
    get observable(): Observable<boolean> {
        return this._clientImplement.observable
    }
    get isNotClosed(): boolean {
        return this._clientImplement.isNotClosed
    }
    get isClosed(): boolean {
        return this._clientImplement.isClosed
    }
    close() {
        this._clientImplement.close()
    }
    /**
     * 返回就是的連接 並發生數據
     * @param data 
     * @returns 
     */
    write(data: string | ArrayBufferLike | Blob | ArrayBufferView): Promise<void> {
        return this.promise().then((ws) => ws.send(data))
    }
    /**
     * 
     * @returns 返回最後緩存的就緒連接
     */
    ws(): WebSocket | undefined {
        return this._clientImplement.ws()
    }
    /**
     * 
     * @returns 返回就緒的 WebSocket 如果沒有就緒的連接 嘗試創建一個新連接
     */
    promise(): Promise<WebSocket> {
        return this._clientImplement.promise(this)
    }
}