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