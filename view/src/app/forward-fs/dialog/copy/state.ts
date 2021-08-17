import { HttpClient } from "@angular/common/http";
import { ServerAPI } from "src/app/core/core/api";
import { SessionService } from "src/app/core/session/session.service";
import { Response } from "src/app/core/session/session_client";
import { SessionRequest } from "src/app/core/session/session_request";
import { EventCode, fromString, sendRequest } from '../event';
import { Clipboard } from '../../manager/settings'

export interface Data {
    src: Clipboard,
    dst: {
        root: string
        dir: string
    }
}
const HeartInterval = 40 * 1000
interface Message {
    event: string
    value: string
}

interface callbacks {
    onProgress(name: string): void
    onExists(name: string): Promise<EventCode>
}
export class Client extends SessionRequest {
    private init_ = '';
    constructor(
        httpClient: HttpClient,
        sessionService: SessionService,
        private readonly data: Data,
        private readonly callbacks: callbacks,
    ) {
        super(ServerAPI.forward.v1.fs.websocketURL(data.src.id, 'copy'), HeartInterval, httpClient, sessionService)
        const src = data.src
        const dst = data.dst
        this.init_ = JSON.stringify({
            event: EventCode.Init,
            srcRoot: src.root,
            srcDir: src.dir,
            dstRoot: dst.root,
            dstDir: dst.dir,
            names: src.names,
            copied: src.copied,
        })
        this._connect()
    }
    optOnOpen(ws: WebSocket, evt: Event): void {
        ws.send(this.init_)
    }
    _onMessage(ws: WebSocket, resp: any) {
        if (resp.code) {
            const msg: Response = resp
            this.reject(new Error(`${msg.code} ${msg.message}`))
            return
        }
        const msg = resp as Message
        const callbacks = this.callbacks
        switch (fromString(msg.event)) {
            case EventCode.Exists:
                callbacks.onExists(msg.value).then((evt) => {
                    sendRequest(ws, evt)
                }).catch((e) => {
                    this.reject(e)
                })
                break
            case EventCode.Progress:
                callbacks.onProgress(msg.value)
                break
            case EventCode.Success:
                this.resolve()
                break
            default:
                console.warn(`unexpected message`, msg)
                this.reject(new Error(`unexpected message: ${msg.event}`))
                break
        }
    }
    _onArrayBuffer(ws: WebSocket, data: ArrayBuffer) {
        this.reject(new Error('not supported ArrayBuffer'))
    }
}