import { HttpClient } from '@angular/common/http'
import { Codes } from '../core/restful'
import { Completer } from '../utils/completer'
import { SessionService } from './session.service'
import { SessionClient } from './session_client'
enum State {
    None,
    Retry,
    Connected,
    Completed,
}
export class SessionRequest extends SessionClient {
    private completer_ = new Completer<void>()
    private state_ = State.None
    constructor(baseURL: string,
        heartInterval: number,
        httpClient: HttpClient,
        sessionService: SessionService) {
        super(baseURL, heartInterval, httpClient, sessionService)
    }
    get result() {
        return this.completer_.promise
    }
    close() {
        if (this.isNotClosed) {
            super.close()
            this._reject('canceled')
        }
    }
    get isCompleted(): boolean {
        return this.state_ == State.Completed
    }
    get isNotCompleted(): boolean {
        return this.state_ != State.Completed
    }
    private _reject(e: string) {
        if (this.isNotCompleted) {
            this.state_ = State.Completed
            this.completer_.reject(new Error(e))
        }
    }
    reject(e: any) {
        if (this.isNotCompleted) {
            this.state_ = State.Completed
            this.completer_.reject(e)
        }
    }
    resolve() {
        if (this.isNotCompleted) {
            this.state_ = State.Completed
            this.completer_.resolve()
        }
    }
    optOnOpenError(ws: WebSocket): void {
        if (this.isCompleted) {
            return
        }
        this._reject('ws connect error')
    }
    optOnClose(ws: WebSocket): void {
        if (this.isCompleted) {
            return
        }

        if (this.state_ == State.None || this.state_ == State.Completed) {
            this.state_ = State.Retry
            super.optOnClose(ws)
        } else {
            this._reject('ws closed')
        }
    }
    _onConnected() {
        super._onConnected()
        if (this.state_ == State.None || this.state_ == State.Retry) {
            this.state_ = State.Connected
        }
    }
    _onConnectError(code?: Codes, message?: string) {
        console.log('_onConnectError', code, message)
        if (code == Codes.Unauthenticated) {
            super._onConnectError(code, message)
            return
        }
        this._reject(`code=${code} message=${message}`)
    }

}
