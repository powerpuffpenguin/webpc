import { Manager, Session } from "src/app/core/session/session"
import { filter, takeUntil } from "rxjs/operators"
import { SessionService } from "src/app/core/session/session.service"
import { Closed } from "src/app/core/utils/closed"
import { HttpClient, HttpParams } from "@angular/common/http"
import { ServerAPI } from "src/app/core/core/api"

interface Writer {
    writeln(text: string): void
    write(text: string): void
}
export class Listener {
    private session_: Session | undefined
    private connectSession_: Session | undefined
    constructor(
        private readonly httpClient: HttpClient,
        public readonly writer: Writer,
        private readonly sessionService: SessionService,
    ) {
        this.sessionService.observable.pipe(
            takeUntil(this.closed_.observable),
            filter((session) => session?.access && session?.userdata.id ? true : false)
        ).subscribe((session) => {
            this.session_ = session
            if (!this.websocket_) {
                this._postConnect()
            }
        })
    }
    pause = false
    close() {
        if (this.closed_.isClosed) {
            return
        }
        this.closed_.close()
        if (this.websocket_) {
            this.websocket_.close()
            this.websocket_ = undefined
        }
        if (this.timeout_) {
            clearTimeout(this.timeout_)
            this.timeout_ = null
        }
        if (this.interval_) {
            clearInterval(this.interval_)
            this.interval_ = null
        }
    }
    private websocket_: WebSocket | undefined
    private closed_ = new Closed()
    private delay_ = 0
    private timeout_: any
    private interval_: any
    private _postConnect() {
        if (!this.pause && this.delay_) {
            this.writer.writeln(`reconnect on ${this.delay_} second after`)
        }
        this.timeout_ = setTimeout(() => {
            this.timeout_ = null
            const session = this.session_
            if (session?.access && session?.userdata?.id) {
                const query = new HttpParams({
                    fromObject: {
                        access_token: session.access,
                    }
                })
                this.connectSession_ = this.session_
                const url = ServerAPI.v1.logger.websocketURL('attach') + '?' + query.toString()
                this._connect(url)
            }
        }, this.delay_ * 1000)
        if (!this.delay_) {
            this.delay_ = 1
        } else if (this.delay_ < 16) {
            this.delay_ *= 2
        }
    }
    private _connect(url: string) {
        if (this.closed_.isClosed) {
            return
        }
        const ws = new WebSocket(url)
        ws.binaryType = "arraybuffer"
        this.websocket_ = ws
        ws.onopen = () => {
            this._onopen(ws)
        }
        ws.onerror = (evt) => {
            ws.close()
        }
        ws.onclose = (evt) => {
            this._onclose(ws, evt)
        }
    }
    private _onopen(ws: WebSocket) {
        if (ws != this.websocket_ || this.closed_.isClosed) {
            ws.close()
            return
        }
        if (!this.interval_) {
            this.interval_ = setInterval(() => {
                try {
                    if (this.websocket_) {
                        this.websocket_.send(`{}`)
                    }
                } catch (e) {
                }
            }, 1000 * 40)
        }
        this.writer.writeln(`attach logger console`)
        ws.onmessage = (evt) => {
            if (ws != this.websocket_) {
                ws.close()
                return
            }
            this.delay_ = 1
            this._onmessage(ws, evt)
        }
    }
    private async _onclose(ws: WebSocket, evt: CloseEvent) {
        if (this.websocket_ != ws) {
            return
        }
        console.log(`ws closed code=${evt.code} reason=${evt.reason} `)
        this.websocket_ = undefined
        if (evt.code == 401 && this.connectSession_) {
            try {
                await Manager.instance.refresh(this.httpClient, this.connectSession_, evt.code, evt.reason)
            } catch (e) {
                console.warn(`refresh token error`, e)
            }
        }
        this._postConnect()
    }
    private _onmessage(ws: WebSocket, evt: MessageEvent) {
        if (typeof evt.data === "string") {
            const str = evt.data.replace(/\n/g, "\r\n")
            this.writer.write(str)
        } else {
            this.writer.writeln(`not supported data type : ${typeof evt.data}`)
            console.log(`not supported data type : ${typeof evt.data}`, evt.data)
        }
    }
}
