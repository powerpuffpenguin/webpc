import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ServerAPI } from 'src/app/core/core/api';
import { BasicState, firstPromise } from 'src/app/core/utils/loader';
export interface SharedResponse {
    name: Array<string>
}
export interface Options {
    readonly target: string
    readonly httpClient: HttpClient
    readonly cancel: Observable<any>
}
export class SharedState extends BasicState<SharedResponse>  {
    constructor(readonly opts: Options,
        onReady?: (data: SharedResponse) => void,
        onError?: (e: any) => void,
    ) {
        super(
            () => {
                const observable = ServerAPI.forward.v1.fs.child('shared').get<SharedResponse>(opts.httpClient, {
                    params: {
                        slave_id: opts.target,
                    },
                })
                return firstPromise(observable, opts.cancel)
            }, onReady, onError)
    }
}