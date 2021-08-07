import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ServerAPI } from 'src/app/core/core/api';
import { BasicState, firstPromise } from 'src/app/core/utils/loader';
export interface DataResponse {
    id: string
    name: string
    description: string
}
export interface VersionResponse {
    platform: string
    version: string
}
export interface StartAtResponse {
    result: number

    at: any
    started: string
}
export interface Options {
    readonly target: string
    readonly httpClient: HttpClient
    readonly cancel: Observable<any>
}
export class DataState extends BasicState<DataResponse>  {
    constructor(readonly opts: Options,
        onReady?: (data: DataResponse) => void,
        onError?: (e: any) => void,
    ) {
        super(
            () => {
                if (opts.target == '0') {
                    return Promise.resolve({
                        id: '0',
                        name: 'Server',
                        description: 'Master Server',
                    })
                }
                const observable = ServerAPI.v1.slaves.child('id', opts.target).get<DataResponse>(opts.httpClient,)
                return firstPromise(observable, opts.cancel)
            }, onReady, onError)
    }

}
export class VersionState extends BasicState<VersionResponse>  {
    constructor(readonly opts: Options,
        onReady?: (data: VersionResponse) => void,
        onError?: (e: any) => void,
    ) {
        super(
            () => {
                const observable = ServerAPI.forward.v1.system.child('version').get<VersionResponse>(opts.httpClient, {
                    params: {
                        slave_id: opts.target,
                    },
                })
                return firstPromise(observable, opts.cancel)
            }, onReady, onError)
    }

}
export class StartAtState extends BasicState<StartAtResponse>  {
    constructor(readonly opts: Options,
        onReady?: (data: StartAtResponse) => void,
        onError?: (e: any) => void,
    ) {
        super(
            () => {
                const observable = ServerAPI.forward.v1.system.child('start_at').get<StartAtResponse>(opts.httpClient, {
                    params: {
                        slave_id: opts.target,
                    },
                })
                return firstPromise(observable, opts.cancel)
            }, onReady, onError)
    }

}
