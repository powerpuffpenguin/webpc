import { HttpClient } from "@angular/common/http"
import { Observable } from "rxjs"
import { ServerAPI } from "src/app/core/core/api"
import { BasicState, firstPromise } from "src/app/core/utils/loader"
import { MB } from "src/app/core/utils/utils"
import { Dir } from "./state"
import { Chunk, UploadFile, Workers } from './state'
const ChunkSize = 5 * MB

export interface HashOptions {
    readonly httpClient: HttpClient
    readonly cancel: Observable<any>
    readonly dir: Dir
    readonly file: UploadFile
}
export interface HashResponse {
    exists: boolean
    hash: string
}
export class HashState extends BasicState<HashResponse>  {
    constructor(readonly opts: HashOptions,
        onReady?: (data: HashResponse) => void,
        onError?: (e: any) => void,
    ) {
        super(
            () => {
                const dir = opts.dir
                const file = opts.file.file
                const observable = ServerAPI.forward.v1.fs.child('hash').get<HashResponse>(opts.httpClient, {
                    params: {
                        slave_id: dir.id,
                        root: dir.root,
                        path: dir.dir + '/' + file.name,
                        size: file.size.toString(),
                        chunk: ChunkSize.toString(),
                    },
                })
                return firstPromise(observable, opts.cancel)
            }, onReady, onError)
    }
}

function makeChunk(file: File): Array<Chunk> {
    const size = file.size
    const chunks = new Array<Chunk>()
    let start = 0
    let index = -1
    while (start != size) {
        let end = start + ChunkSize
        if (end > size) {
            end = size
        }
        ++index
        const chunk = new Chunk(file, index, start, end)
        chunks.push(chunk)
        start = end
    }
    return chunks
}
export class FileState extends BasicState<string>  {
    constructor(readonly file: UploadFile,
        onReady?: (data: string) => void,
        onError?: (e: any) => void,
    ) {
        super(
            () => {
                const val = file.hash
                if (val != '') {
                    return Promise.resolve(val)
                }
                const chunks = makeChunk(file.file)
                return Workers.instance.calculate(chunks).then((resp) => {
                    file.chunks = chunks
                    return resp
                })
            }, onReady, onError)
    }
}