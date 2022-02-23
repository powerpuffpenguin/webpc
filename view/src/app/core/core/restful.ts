import { HttpHeaders, HttpParams, HttpClient, HttpUrlEncodingCodec } from '@angular/common/http';
import { Observable } from 'rxjs';
export enum Codes {
    OK = 0,
    Canceled = 1,
    Unknown = 2,
    InvalidArgument = 3,
    DeadlineExceeded = 4,
    NotFound = 5,
    AlreadyExists = 6,
    PermissionDenied = 7,
    ResourceExhausted = 8,
    FailedPrecondition = 9,
    Aborted = 10,
    OutOfRange = 11,
    Unimplemented = 12,
    Internal = 13,
    Unavailable = 14,
    DataLoss = 15,
    Unauthenticated = 16,
    _maxCode = 17
}
export function codesString(code: Codes): string {
    switch (code) {
        case Codes.OK:
            return 'Ok'
        case Codes.Canceled:
            return 'Canceled'
        case Codes.Unknown:
            return 'Unknown'
        case Codes.InvalidArgument:
            return 'InvalidArgument'
        case Codes.DeadlineExceeded:
            return 'DeadlineExceeded'
        case Codes.NotFound:
            return 'NotFound'
        case Codes.AlreadyExists:
            return 'AlreadyExists'
        case Codes.PermissionDenied:
            return 'PermissionDenied'
        case Codes.ResourceExhausted:
            return 'ResourceExhausted'
        case Codes.FailedPrecondition:
            return 'FailedPrecondition'
        case Codes.Aborted:
            return 'Aborted'
        case Codes.OutOfRange:
            return 'OutOfRange'
        case Codes.Unimplemented:
            return 'Unimplemented'
        case Codes.Internal:
            return 'Internal'
        case Codes.Unavailable:
            return 'Unavailable'
        case Codes.DataLoss:
            return 'DataLoss'
        case Codes.Unauthenticated:
            return 'Unauthenticated'
    }
    return `Unknown`
}
export class NetError {
    constructor(
        public readonly status: number, // http status code
        public readonly grpc: Codes,// grpc code
        public readonly message: string, // string message
    ) {
    }
    private str_: string | undefined
    toString(): string {
        if (typeof this.str_ !== "string") {
            const result = new Array<string>()
            if (typeof this.status === "number") {
                result.push(`http=${this.status}`)
            }
            if (typeof this.grpc === "number" && this.grpc != 2) {
                result.push(`grpc=${this.grpc}`)
            }
            if (typeof this.message === "string" && this.message.length > 0) {
                if (result.length != 0) {
                    result.push(`msg=${this.message}`)
                }
            }
            this.str_ = result.join(',')
        }
        return this.str_
    }
}
interface Err {
    status: number
    statusText: string
    message?: string
    error?: string | {
        code?: Codes
        message?: string
        description?: string
        error?: {
            message?: string
        }
    }
}
export function resolveError(e: any): NetError {
    console.warn(e)
    if (!e) {
        return new NetError(200, Codes.OK, 'success')
    }
    if (typeof e === "string") {
        return new NetError(500, Codes.Unknown, e)
    }
    if (e !== null && typeof e === 'object' && typeof e.status === "number") {
        return resolveHttpError(e)
    }
    return new NetError(500, Codes.Unknown, e.toString())
}
export function resolveHttpError(e: Err): NetError {
    let error = e.error
    const status = e.status
    let grpc = Codes.Unknown
    if (typeof error === "string") {
        return new NetError(status, grpc, error)
    }
    if (error) {
        if (typeof error.code === "number") {
            return new NetError(status, error.code, error.message ?? `rpc error: code = ${codesString(error.code)}`)
        } else if (typeof error.message === "string") {
            if (typeof error.code === "number") {
                grpc = error.code
            }
            return new NetError(status, grpc, error.message)
        } else if (typeof error.description === "string") {
            return new NetError(status, grpc, error.description)
        } else if (error.error) {
            error = error.error
            if (typeof error.message === "string") {
                return new NetError(status, grpc, error.message)
            }
            return new NetError(status, grpc, error.toString())
        }
    } else if (e.message) {
        return new NetError(status, grpc, e.message)
    }
    return new NetError(status, grpc, e.statusText)
}
function wrapObservable<T>() {
    return (source: Observable<T>) => {
        return new Observable<T>(subscriber => {
            source.subscribe({
                next(v) {
                    try {
                        subscriber.next(v)
                    } catch (e) {
                        subscriber.error(resolveError(e))
                    }
                },
                error(e) {
                    subscriber.error(resolveError(e))
                },
                complete() {
                    subscriber.complete()
                },
            })
        })
    }
}
export function MakeRESTful(...path: Array<string | number | boolean>): RESTful {
    let url = ''
    if (path && path.length > 0) {
        for (let i = 0; i < path.length; i++) {
            const codec = new HttpUrlEncodingCodec()
            path[i] = codec.encodeKey(path[i].toString())
        }
        url += '/' + path.join('/')
    }
    return new RESTful(url)
}
export class RESTful {
    constructor(public readonly baseURL: string) {
    }
    httpURL(...path: Array<string | number | boolean>): string {
        let url = this.baseURL
        if (path && path.length > 0) {
            for (let i = 0; i < path.length; i++) {
                const codec = new HttpUrlEncodingCodec()
                path[i] = codec.encodeKey(path[i].toString())
            }
            url += '/' + path.join('/')
        }
        return url
    }
    websocketURL(
        ...path: Array<string | number | boolean>
    ): string {
        const location = document.location
        let addr: string
        // console.log(location.protocol)
        if (location.protocol == "https:") {
            addr = `wss://${location.hostname}`
            if (location.port == "" || location.port == "443") {
                // addr += ":443"
            } else {
                addr += `:${location.port}`
            }
        } else {
            addr = `ws://${location.hostname}`
            if (location.port == "" || location.port == "80") {
                // addr += ":80"
            } else {
                addr += `:${location.port}`
            }
        }
        return `${addr}${this.httpURL(...path)}`
    }
    child(...path: Array<string | number | boolean>): RESTful {
        return new RESTful(this.httpURL(...path))
    }
    get<T>(client: HttpClient,
        options?: {
            headers?: HttpHeaders | {
                [header: string]: string | string[];
            };
            observe?: 'body';
            params?: HttpParams | {
                [param: string]: string | string[];
            };
            reportProgress?: boolean;
            responseType?: 'json';
            withCredentials?: boolean;
        },
    ): Observable<T>;
    get(client: HttpClient,
        options: {
            headers?: HttpHeaders | {
                [header: string]: string | string[];
            };
            observe?: 'body';
            params?: HttpParams | {
                [param: string]: string | string[];
            };
            reportProgress?: boolean;
            responseType: 'text';
            withCredentials?: boolean;
        },
    ): Observable<string>;
    get(client: HttpClient, options?: any): any {
        return client.get(this.baseURL, options).pipe(wrapObservable())
    }
    post<T>(client: HttpClient, body: any | null,
        options?: {
            headers?: HttpHeaders | {
                [header: string]: string | string[];
            };
            observe?: 'body';
            params?: HttpParams | {
                [param: string]: string | string[];
            };
            reportProgress?: boolean;
            responseType?: 'json';
            withCredentials?: boolean;
        },
    ): Observable<T> {
        return client.post<T>(this.baseURL, body, options).pipe(wrapObservable())
    }
    delete<T>(client: HttpClient,
        options?: {
            headers?: HttpHeaders | {
                [header: string]: string | string[];
            };
            observe?: 'body';
            params?: HttpParams | {
                [param: string]: string | string[];
            };
            reportProgress?: boolean;
            responseType?: 'json';
            withCredentials?: boolean;
        },
    ): Observable<T> {
        return client.delete<T>(this.baseURL, options).pipe(wrapObservable())
    }
    put<T>(client: HttpClient, body: any | null,
        options?: {
            headers?: HttpHeaders | {
                [header: string]: string | string[];
            };
            observe?: 'body';
            params?: HttpParams | {
                [param: string]: string | string[];
            };
            reportProgress?: boolean;
            responseType?: 'json';
            withCredentials?: boolean;
        },
    ): Observable<T> {
        return client.put<T>(this.baseURL, body, options).pipe(wrapObservable())
    }
    patch<T>(client: HttpClient, body: any | null,
        options?: {
            headers?: HttpHeaders | {
                [header: string]: string | string[];
            };
            observe?: 'body';
            params?: HttpParams | {
                [param: string]: string | string[];
            };
            reportProgress?: boolean;
            responseType?: 'json';
            withCredentials?: boolean;
        },
    ): Observable<T> {
        return client.patch<T>(this.baseURL, body, options).pipe(wrapObservable())
    }
}