import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ServerAPI } from 'src/app/core/core/api';
import { BasicState, firstPromise } from 'src/app/core/utils/loader';
export interface ListResult {
    id: string
    name: string
    attached: boolean
}
export interface ListResponse {
    result: Array<ListResult>
}
export interface Options {
    readonly target: string
    readonly httpClient: HttpClient
    readonly cancel: Observable<any>
}
export class ListState extends BasicState<ListResponse>  {
    constructor(readonly opts: Options,
        onReady?: (data: ListResponse) => void,
        onError?: (e: any) => void,
    ) {
        super(
            () => {
                const observable = ServerAPI.forward.v1.shell.child('list').get<ListResponse>(opts.httpClient, {
                    params: {
                        slave_id: opts.target,
                    },
                })
                return firstPromise(observable, opts.cancel)
            }, onReady, onError)
    }
}