import { HttpClient } from "@angular/common/http";
import { Subject } from "rxjs";
import { ServerAPI } from "src/app/core/core/api";
import { Codes } from "src/app/core/core/restful";
import { SessionService } from "src/app/core/session/session.service";
import { SessionRequest } from "src/app/core/session/session_request";
import { Terminal } from "xterm";
const HeartInterval = 40 * 1000
enum EventCode {
    Resize = 3,
    FontSize = 4,
    FontFamily = 5,

    Info = 7,
}
export interface Info {
    id: string
    name: string
    at: number
    fontSize: number
    fontFamily: string
}

export interface Target {
    id: string
    shellid: string
}

export class Shell extends SessionRequest {
    constructor(
        httpClient: HttpClient,
        sessionService: SessionService,
        public readonly id: string,
        public shellid: string,
        private readonly xterm: Terminal,
        private readonly subject: Subject<Info | undefined>,
    ) {
        super(
            ServerAPI.forward.v1.shell.websocketURL(id, shellid == 'new' ? '0' : shellid, xterm.cols, xterm.rows),
            HeartInterval,
            httpClient, sessionService,
        )
        this._connect()
    }
    _onConnectError(code?: Codes, message?: string) {
        if (code != Codes.Unauthenticated && typeof message === "string") {
            this.xterm.writeln(message)
        }
        super._onConnectError(code, message)
    }
    _onConnected() {
        super._onConnected()
        this.xterm.focus()
        this.xterm.setOption("cursorBlink", true)
    }
    _onMessage(ws: WebSocket, resp: any) {
        if (resp.event == "Info") {
            this.shellid = resp.id
            this.subject.next(resp)
        } else {
            console.warn('unknow msg: ', resp)
        }
    }
    _onArrayBuffer(ws: WebSocket, data: ArrayBuffer) {
        this.xterm.write(new Uint8Array(data))
    }
    send(data: string) {
        const ws = this.ws()
        if (ws) {
            ws.send(new TextEncoder().encode(data))
        }
    }
    sendResize(cols: number, rows: number) {
        const ws = this.ws()
        if (ws) {
            ws.send(JSON.stringify({
                event: EventCode.Resize,
                cols: cols,
                rows: rows,
            }))
        }
    }
    sendFontSize(fontSize: number) {
        const ws = this.ws()
        if (ws) {
            ws.send(JSON.stringify({
                event: EventCode.FontSize,
                fontSize: fontSize,
            }))
        }
    }
    sendFontFamily(fontFamily: string) {
        const ws = this.ws()
        if (ws) {
            ws.send(JSON.stringify({
                event: EventCode.FontFamily,
                fontFamily: fontFamily,
            }))
        }
    }
}