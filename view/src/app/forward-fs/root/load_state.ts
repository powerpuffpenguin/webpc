import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { ServerAPI } from 'src/app/core/core/api';
import { BasicState, firstPromise } from 'src/app/core/utils/loader';
export interface MountResponse {
    name: Array<string>
}
export interface Options {
    readonly target: string
    readonly httpClient: HttpClient
    readonly cancel: Observable<any>
}
export class MountState extends BasicState<MountResponse>  {
    constructor(readonly opts: Options,
        onReady?: (data: MountResponse) => void,
        onError?: (e: any) => void,
    ) {
        super(
            () => {
                const observable = ServerAPI.forward.v1.fs.child('mount').get<MountResponse>(opts.httpClient, {
                    params: {
                        slave_id: opts.target,
                    },
                })
                return firstPromise(observable, opts.cancel)
            }, onReady, onError)
    }
}