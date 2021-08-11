import { SessionService } from "src/app/core/session/session.service"
import { HttpClient } from "@angular/common/http"
import { ServerAPI } from "src/app/core/core/api"
import { ClientOption } from "src/app/core/net/client"
import { Access } from "src/app/core/net/access"
import { environment } from "src/environments/environment"
import { interval } from "rxjs"
import { Codes } from "src/app/core/core/restful"
import { SessionClient } from "src/app/core/session/session_client"
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
    writeln(text: string, log?: boolean): void
    write(text: string, log?: boolean): void
}

export class Listener extends SessionClient implements ClientOption {
    constructor(httpClient: HttpClient,
        sessionService: SessionService,
        private readonly writer: Writer,
    ) {
        super(ServerAPI.v1.logger.websocketURL('attach'),
            HeartInterval,
            httpClient, sessionService,
        )
        this._connect()
    }
    optOnOpenError(_: WebSocket): void {
        const delay = this._onclose()
        if (delay <= 0) {
            this.writer.writeln(`connect err, retrying in 0s`, true)
        } else {
            this.writer.writeln(`connect err, retrying in ${delay}s`, true)
        }
    }
    _onConnected() {
        this.writer.writeln(`logger attach`, true)
    }
    _onConnectError(code?: Codes, message?: string) {
        this.writer.writeln(`connect err: ${code} ${message}`, true)
    }
    _onMessage(resp: Response) {
        if (resp.code == Codes.OK) {
            console.log('ws recv: ', resp)
        } else {
            console.warn('ws recv: ', resp)
        }
    }
    _onArrayBuffer(data: ArrayBuffer) {
        const enc = new TextDecoder("utf-8")
        let str = enc.decode(data)
        str = str.replace(/\n/g, "\r\n")
        this.writer.write(str)
    }
}
