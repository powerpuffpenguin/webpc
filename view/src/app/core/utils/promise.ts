export function isPromise<T = any>(obj: any): obj is Promise<T> {
    return !!obj &&
        (typeof obj === 'object' || typeof obj === 'function') &&
        typeof obj.then === 'function'
}