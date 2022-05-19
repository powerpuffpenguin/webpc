import { Queue, Value } from './containers'
export class Completer<T>{
    private promise_: Promise<T> | undefined
    private resolve_: ((value?: T | PromiseLike<T>) => void) | undefined
    private reject_: ((reason?: any) => void) | undefined
    constructor() {
        this.promise_ = new Promise<T>((resolve, reject) => {
            this.resolve_ = resolve as any
            this.reject_ = reject
        });
    }

    get promise(): Promise<T> {
        return this.promise_ as Promise<T>
    }
    resolve(value?: T | PromiseLike<T>) {
        if (this.resolve_) {
            this.resolve_(value)
        }
    }
    reject(reason?: any) {
        if (this.reject_) {
            this.reject_(reason)
        }
    }
}
export class Mutex {
    private completer_: Completer<void> | undefined
    async lock(): Promise<void> {
        while (true) {
            if (!this.completer_) {
                this.completer_ = new Completer<void>()
                break
            }
            await this.completer_.promise
        }
    }
    tryLock(): boolean {
        if (!this.completer_) {
            this.completer_ = new Completer<void>()
            return true
        }
        return false
    }
    unlock() {
        if (!this.completer_) {
            throw new Error('not locked')
        }

        const completer = this.completer_
        this.completer_ = undefined
        completer.resolve()
    }
    get isLocked(): boolean {
        if (this.completer_) {
            return true
        }
        return false
    }
    get isNotLocked(): boolean {
        if (this.completer_) {
            return false
        }
        return true
    }
}
export class Channel<T>{
    private queue_: Queue<T>
    private closed_ = false;
    get isClosed(): boolean {
        return this.closed_
    }
    get isNotClosed(): boolean {
        return !this.closed_
    }
    constructor(buffer: number) {
        const n = Math.floor(buffer)
        if (n < 1) {
            throw new Error('buffer must > 0')
        }
        this.queue_ = new Queue<T>(buffer)
    }
    close(): boolean {
        if (this.closed_) {
            return false
        }
        this.closed_ = true
        this._resolveSend()
        return true
    }
    private recv_: Completer<void> | undefined
    private _resolveRecv() {
        const completer = this.recv_
        if (completer) {
            this.recv_ = undefined
            completer.resolve()
        }
    }
    tryRecv(): Value<T> {
        const result = this.queue_.pop()
        this._resolveSend()
        return result
    }
    async recv(): Promise<Value<T>> {
        while (!this.closed_) {
            const result = this.queue_.pop()
            this._resolveSend()
            if (result.done) {
                return result
            }
            if (this.recv_) {
                await this.recv_.promise
            } else {
                this.recv_ = new Completer<void>()
            }
        }
        return this.queue_.pop()
    }
    private send_: Completer<void> | undefined
    private _resolveSend() {
        const completer = this.send_
        if (completer) {
            this.send_ = undefined
            completer.resolve()
        }
    }
    trySend(value: T): boolean {
        const queue = this.queue_
        if (queue.isFull) {
            return false
        }
        queue.push(value)
        this._resolveRecv()
        return true
    }

    async send(value: T): Promise<boolean> {
        const queue = this.queue_
        while (!this.closed_) {
            if (!queue.isFull) {
                queue.push(value)
                this._resolveRecv()
                return true
            }
            if (this.send_) {
                await this.send_.promise
            } else {
                this.send_ = new Completer<void>()
            }
        }
        return false
    }
    clear() {
        if (this.closed_) {
            return
        }
        const queue = this.queue_
        if (queue.isEmpty) {
            return
        }
        queue.clear()
        this._resolveSend()
    }
}