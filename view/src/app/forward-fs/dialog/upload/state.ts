import { sizeString } from "src/app/core/utils/utils";
import { buf } from "crc-32";
import { Mutex } from "src/app/core/utils/completer";
export interface Dir {
    id: string
    root: string
    dir: string
}
export enum States {
    Nil,
    Working,
    Ok,
    Error,
    Skip,
}
export class UploadFile {
    constructor(public file: File) {
    }
    state = States.Nil
    progress: number = 0
    error: any
    get sizeString(): string {
        return sizeString(this.file.size)
    }
    get key(): string {
        const file = this.file
        return `${file.size}${file.lastModified}${file.name}`
    }
    isWorking(): boolean {
        return this.state == States.Working
    }
    isOk(): boolean {
        return this.state == States.Ok
    }
    isError(): boolean {
        return this.state == States.Error
    }
    isSkip(): boolean {
        return this.state == States.Skip
    }
    hash = ''
    chunks: Array<Chunk> = []
}
export class Chunk {
    hash = ''
    constructor(public file: File, public index: number, public start: number, public end: number) {
    }
    arrayBuffer(): Promise<ArrayBuffer> {
        return this.file.slice(this.start, this.end).arrayBuffer()
    }
}
// calculate request
interface Request {
    file: File
    start: number
    end: number
}
interface Response {
    vals: Array<string>
    hash: string
}
function crc32ToHex(val: number): string {
    const hash = new Uint8Array(4)
    const dataView = new DataView(hash.buffer)
    dataView.setUint32(0, val, false)
    const strs = new Array<string>(4)
    for (let i = 0; i < 4; i++) {
        strs[i] = dataView.getUint8(i).toString(16).padStart(2, '0')
    }
    return strs.join('')
}
async function calculateHash(reqs: Array<Request>): Promise<Response> {
    const vals = new Array<string>(reqs.length)
    const hash = new Uint8Array(reqs.length * 4)
    const dataView = new DataView(hash.buffer)
    for (let i = 0; i < reqs.length; i++) {
        const req = reqs[i]
        const buffer = await req.file.slice(req.start, req.end).arrayBuffer()
        const val = buf(new Uint8Array(buffer), 0)
        vals[i] = crc32ToHex(val)
        dataView.setUint32(i * 4, val, false)
    }
    return {
        vals: vals,
        hash: crc32ToHex(buf(hash, 0)),
    }
}
export class Workers {
    private static instance_: Workers | undefined
    static get instance() {
        if (!this.instance_) {
            this.instance_ = new Workers()
        }
        return this.instance_
    }
    private constructor() {
    }

    private mutex_ = new Mutex()
    private worker_: Worker | undefined
    private getWorker(): Worker | undefined {
        if (this.worker_) {
            return this.worker_
        }
        try {
            this.worker_ = new Worker(new URL('./hash.worker', import.meta.url), { type: 'module' })
        } catch (e) {
            console.warn("new Worker('./hash.worker') error: ", e)
        }
        return this.worker_
    }
    calculate(chunks: Array<Chunk>): Promise<string> {
        let promise: Promise<Response>
        if (typeof Worker === "undefined") {
            promise = calculateHash(chunks)
        } else {
            const worker = this.getWorker()
            if (worker) {
                promise = this._asyncCalculate(worker, chunks)
            } else {
                promise = calculateHash(chunks)
            }
        }
        return promise.then((resp) => {
            const vals = resp.vals
            for (let i = 0; i < vals.length; i++) {
                chunks[i].hash = vals[i]
            }
            return resp.hash
        })
    }
    private async _asyncCalculate(worker: Worker, chunks: Array<Chunk>): Promise<Response> {
        await this.mutex_.lock()
        try {
            return await this._calculate(worker, chunks)
        } finally {
            this.mutex_.unlock()
        }
    }
    private _calculate(worker: Worker, chunks: Array<Chunk>): Promise<Response> {
        return new Promise<Response>((resolve, reject) => {
            // recv
            worker.onmessage = ({ data }) => {
                if (data) {
                    if (data.error) {
                        reject(data.error)
                    } else if (Array.isArray(data.vals) && typeof data.hash === "string") {
                        for (let i = 0; i < data.vals.length; i++) {
                            const val = data.vals[i]
                            if (typeof val != "string") {
                                console.warn('unknow worker result', i, val)
                                reject(`unknow worker result`)
                                return
                            }
                        }
                        resolve(data)
                    } else {
                        console.warn('unknow worker result', data)
                        reject(`unknow worker result`)
                    }
                } else {
                    console.warn('unknow worker result', data)
                    reject(`unknow worker result`)
                }
            }
            // send
            try {
                worker.postMessage(
                    chunks.map((chunk) => {
                        return {
                            file: chunk.file,
                            start: chunk.start,
                            end: chunk.end,
                        }
                    })
                )
            } catch (e) {
                reject(e)
                return
            }
        })
    }
}