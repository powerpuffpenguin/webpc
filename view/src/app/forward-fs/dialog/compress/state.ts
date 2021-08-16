import { HttpClient, HttpParams } from "@angular/common/http";
import { ServerAPI } from "src/app/core/core/api";
import { SessionService } from "src/app/core/session/session.service";
import { Response } from "src/app/core/session/session_client";
import { SessionRequest } from "src/app/core/session/session_request";
import { FileInfo } from "../../fs";
import { EventCode, fromString, sendRequest } from '../event';

const HeartInterval = 40 * 1000
interface Message {
    event: string
    value: string
    info: FileInfo,
}

export enum Algorithm {
    Tar = 1,
    Zip = 2,
    TarGZ = 3,
}

interface callbacks {
    onProgress(name: string): void
    onExists(name: string): Promise<boolean>
}
export class Client extends SessionRequest {
    private init_ = '';
    constructor(
        id: string,
        httpClient: HttpClient,
        sessionService: SessionService,
        root: string, dir: string,
        dst: string, source: Array<string>,
        algorithm: Algorithm,
        private readonly callbacks: callbacks,
    ) {
        super(ServerAPI.forward.v1.fs.websocketURL(id, 'compress'), HeartInterval, httpClient, sessionService)
        this.init_ = JSON.stringify({
            event: EventCode.Init,
            root: root,
            dir: dir,
            dst: dst,
            source: source,
            algorithm: algorithm,
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
                callbacks.onExists(msg.value).then((ok) => {
                    if (ok) {
                        sendRequest(ws, EventCode.Yes)
                    } else {
                        sendRequest(ws, EventCode.No)
                    }
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