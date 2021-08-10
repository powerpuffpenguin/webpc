import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ServerAPI } from 'src/app/core/core/api';
import { BasicState, firstPromise } from 'src/app/core/utils/loader';
import { FileInfo, Dir } from '../fs';
export interface ListResponse {
    dir: Dir
    items: Array<FileInfo>
}

export interface Options {
    readonly target: string
    readonly httpClient: HttpClient
    readonly cancel: Observable<any>
}
export class ListState extends BasicState<ListResponse>  {
    constructor(readonly opts: Options,
        root: string,
        path: string,
        onReady?: (data: ListResponse) => void,
        onError?: (e: any) => void,
    ) {
        super(
            () => {
                const observable = ServerAPI.forward.v1.fs.child('list').get<ListResponse>(opts.httpClient, {
                    params: {
                        slave_id: opts.target,
                        root: root,
                        path: path,
                    },
                })
                return firstPromise(observable, opts.cancel)
            }, onReady, onError)
    }
}