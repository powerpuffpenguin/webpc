import { HttpClient } from "@angular/common/http";
import { MatDialog } from "@angular/material/dialog";
import { Completer } from "src/app/core/utils/completer";
import { Closed } from "src/app/core/utils/closed";
import { EventCode } from "../event";
import { Loader } from "src/app/core/utils/loader";
import { FileState, HashState } from "./load_state";
import { Chunk, States, UploadFile } from "./state";
import { Dir } from "./state";
import { ExistsChoiceComponent } from "../exists-choice/exists-choice.component";
import { ServerAPI } from "src/app/core/core/api";
import { takeUntil } from "rxjs/operators";

export class Uploader {
    private completer_ = new Completer<undefined>()
    private closed_ = new Closed()
    get isClosed(): boolean {
        return this.closed_.isClosed
    }
    get isNotClosed(): boolean {
        return this.closed_.isNotClosed
    }
    close() {
        if (this.isNotClosed) {
            this.completer_.reject(new Error('closed'))
            this.closed_.close()
            this.task_?.close()
        }
    }
    private _resolve() {
        if (this.isNotClosed) {
            this.completer_.resolve()
            this.closed_.close()
        }
    }
    style = EventCode.Universal
    constructor(private readonly dir: Dir,
        private readonly file: UploadFile,
        private readonly httpClient: HttpClient,
        private readonly matDialog: MatDialog,
    ) { }

    serve(): Promise<undefined> {
        const promise = this.completer_.promise
        this._serve().then(() => {
            if (this.isClosed) {
                return
            }
            const file = this.file
            if (file.state == States.Working) {
                file.progress = 100
                file.state = States.Ok
            }
            this._resolve()
        }).catch((e) => {
            if (this.isClosed) {
                return
            }
            const file = this.file
            if (file.state == States.Working) {
                file.error = e
                file.state = States.Error
            }
            this._resolve()
        })
        return promise
    }
    private async _serve() {
        const file = this.file
        file.state = States.Working
        file.progress = 0.10
        file.error = null
        let hash = ''
        let exists = false
        const loder = new Loader([
            new HashState({
                httpClient: this.httpClient,
                cancel: this.closed_.observable,
                dir: this.dir,
                file: file,
            }, (resp) => {
                file.progress += 0.4
                if (typeof resp.hash === "string") {
                    hash = resp.hash
                }
                if (typeof resp.exists === "boolean" && resp.exists) {
                    exists = true
                }
            }, (e) => {
                console.warn('hash state error: ', e)
            }),
            new FileState(file, (resp) => {
                file.progress += 0.3
                file.hash = resp
            }, (e) => {
                console.warn('file state error: ', e)
            })
        ])
        await loder.load()
        if (this.isClosed) {
            return
        }

        if (exists) {
            if (hash == file.hash) {
                return
            }
            if (this.style == EventCode.SkipAll) {
                file.state = States.Skip
                return
            } else if (this.style != EventCode.YesAll) {
                const style = await this.matDialog.open(ExistsChoiceComponent, {
                    data: this.file.file.name,
                    disableClose: true,
                }).afterClosed().toPromise()
                if (typeof style === "number") {
                    switch (style) {
                        case EventCode.YesAll:
                            this.style = style
                            break
                        case EventCode.Yes:
                            break
                        case EventCode.SkipAll:
                            this.style = style
                            file.state = States.Skip
                            return
                        case EventCode.Skip:
                            file.state = States.Skip
                            return
                        default:
                            file.state = States.Error
                            file.error = 'already exists'
                            return
                    }
                }
            }
        }

        const task = new Task(3,
            this.dir, file, this.httpClient,
            file.chunks, file.hash,
        )
        this.task_ = task
        await task.serve()
    }
    private task_: Task | undefined
}

interface ChunkResponse {
    result: Array<string>
}
enum TaskState {
    None,
    Working,
    End,
}
class Task {
    private completer_ = new Completer<void>()
    private closed_ = new Closed()
    private state_ = TaskState.None
    constructor(public readonly concurrent: number,
        private readonly dir: Dir,
        private readonly file: UploadFile,
        private readonly httpClient: HttpClient,
        private readonly chunks: Array<Chunk>,
        private readonly hash: string,
    ) {
    }
    close() {
        if (this.closed_.isNotClosed) {
            this.closed_.close()
            this._reject(new Error('closed'))
        }
    }
    get isClosed(): boolean {
        return this.closed_.isClosed
    }
    private _checkClosed() {
        if (this.isClosed) {
            new Error('closed')
        }
    }
    get isNotClosed(): boolean {
        return this.closed_.isNotClosed
    }
    get isWorking(): boolean {
        return this.state_ == TaskState.Working
    }
    get isNotWorking(): boolean {
        return this.state_ == TaskState.Working
    }
    private resolve() {
        if (this.state_ == TaskState.Working) {
            this.state_ = TaskState.End
            if (this.closed_.isClosed) {
                this.completer_.reject(new Error('closed'))
            } else {
                this.completer_.resolve()
            }
        }
    }
    private _reject(e: any) {
        if (this.state_ == TaskState.Working) {
            this.state_ = TaskState.End
            this.completer_.reject(e)
        }
    }
    async serve(): Promise<void> {
        const completer = this.completer_
        if (this.state_ != TaskState.None) {
            return completer.promise
        }
        this.state_ = TaskState.Working
        this._init().then((hash) => {
            let concurrent = this.concurrent
            if (concurrent < 1) {
                concurrent = 1
            }
            let wait = 0
            for (let i = 0; i < concurrent; i++) {
                wait++
                this._run(hash).then(() => {
                    wait--
                    if (!wait) {
                        this._merge().then(() => {
                            this.resolve()
                        }).catch((e) => {
                            this._reject(e)
                        })
                    }
                }).catch((e) => {
                    this._reject(e)
                })
            }
        }).catch((e) => {
            this._reject(e)
        })
        return completer.promise
    }
    get path(): string {
        if (this.dir.dir.endsWith("/")) {
            return this.dir.dir + this.file.file.name
        }
        return this.dir.dir + '/' + this.file.file.name
    }
    private async _init(): Promise<Array<string>> {
        const dir = this.dir
        const resp = await ServerAPI.forward.v1.fs.child('chunk').get<ChunkResponse>(this.httpClient, {
            params: {
                slave_id: dir.id,
                root: dir.root,
                path: this.path,
                offset: "0",
                count: this.chunks.length.toString(),
            },
        }).pipe(
            takeUntil(this.closed_.observable)
        ).toPromise()
        this._checkClosed()
        if (!Array.isArray(resp.result)) {
            throw new Error("unknow result")
        } else if (resp.result.length != this.chunks.length) {
            throw new Error("result not match")
        }
        return resp.result
    }
    private chunk_ = 0
    private getChunk(): FindChunk | undefined {
        const chunks = this.chunks
        let i = this.chunk_
        if (i >= chunks.length) {
            return
        }
        this.chunk_++
        return {
            i: i,
            chunk: chunks[i]
        }
    }
    private progress_ = -1
    private async _run(hash: Array<string>) {
        const file = this.file
        const chunks = this.chunks
        while (this.isNotClosed) {
            const chunk = this.getChunk()
            if (!chunk) {
                break
            }
            const i = chunk.i
            let progress = i + 1
            if (progress > this.progress_) {
                this.progress_ = progress
                progress = Math.floor(progress * 100 / chunks.length)
                if (progress < 1) {
                    progress = 1
                } else if (progress > 99) {
                    progress = 99.10
                } else if (progress == 100) {
                    progress = 99.90
                }
                file.progress = progress
            }
            if (chunks[i].hash != hash[i]) {
                await this._put(chunk.i, chunk.chunk)
                this._checkClosed()
            }
        }
    }
    private async _put(i: number, chunk: Chunk) {
        const dir = this.dir
        const buffer = await chunk.arrayBuffer()

        console.log("upload", {
            id: dir.id,
            root: dir.root,
            i: i,
            path: this.path,
        })
        console.log(dir)
        await ServerAPI.forward.v1.fs.child('upload', dir.id, dir.root, i, this.path).post(this.httpClient, buffer).pipe(
            takeUntil(this.closed_.observable)
        ).toPromise()
    }
    private _merge() {
        const dir = this.dir
        return ServerAPI.forward.v1.fs.child('merge').post(this.httpClient, {
            root: dir.root,
            path: this.path,
            hash: this.hash,
            count: this.chunks.length,
        }, {
            params: {
                slave_id: dir.id,
            },
        }).pipe(
            takeUntil(this.closed_.observable)
        ).toPromise()
    }
}
interface FindChunk {
    i: number
    chunk: Chunk
}