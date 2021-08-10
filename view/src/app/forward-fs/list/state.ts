import { HttpClient } from "@angular/common/http"
import { Closed } from "src/app/core/utils/closed"
import { Loader } from "src/app/core/utils/loader"
import { FileInfo, Dir } from "../fs";
import { ListState, ListResponse } from './load_state';
export interface Error {
    id: string
    err: any
}
const DefaultValue: any = {}
export class State {
    ready = false
    loader: Loader = DefaultValue

    errs: Array<Error> = []
    hasErr = false
    closed = new Closed()

    source: Array<FileInfo> = []
    dir: Dir = {
        root: '',
        read: false,
        write: false,
        shared: false,
        dir: '',
        id: '',
    }
    constructor(readonly httpClient: HttpClient,
        public readonly target: string,
        public readonly root: string,
        public readonly path: string,
    ) {
        if (target == '') {
            return
        }
        const opts = {
            httpClient: httpClient,
            target: target,
            cancel: this.closed.observable,
        }
        this.loader = new Loader([
            new ListState(opts, root, path, (resp) => {
                const dir = resp.dir
                dir.id = target
                this.dir = dir
                if (Array.isArray(resp.items) && resp.items.length > 0) {
                    this.source = new Array<FileInfo>()
                    for (let i = 0; i < resp.items.length; i++) {
                        this.source.push(new FileInfo(dir.root, dir.dir, resp.items[i]))
                    }
                    this.source.sort(FileInfo.compare)
                } else {
                    this.source = []
                }
            }, (e) => {
                this.errs.push({
                    id: 'ListState',
                    err: e,
                })
            }),
        ])
    }
    refresh() {
        this.hasErr = false
        this.errs = []
        this.loader.load().then(() => {
            this.ready = true
        }).catch((_) => {
            this.hasErr = true
        })
    }
}