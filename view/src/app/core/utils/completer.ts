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