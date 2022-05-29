import { HttpClient, HttpParams } from "@angular/common/http"
import { ServerAPI } from "src/app/core/core/api"
import { Channel, Completer } from "src/app/core/utils/completer"
import { FileInfo, Dir } from '../../fs';
import { DB } from "./db";

export interface Pair {
    name: string
    path: string
}
export class Path {
    public readonly dirs: Array<Pair>
    constructor(public readonly id: string,
        public readonly root: string,
        public readonly path: string,
    ) {

        const strs = path.split('/')
        const dirs = new Array<Pair>()
        let dir = ''
        for (let i = 0; i < strs.length; i++) {
            const str = strs[i]
            if (str != "") {
                dir += '/' + str
                dirs.push({
                    name: str,
                    path: dir,
                })
            }
        }
        this.dirs = dirs
    }

    equal(other: Path): boolean {
        return this.id == other.id &&
            this.root == other.root &&
            this.path == other.path
    }

    get key(): string {
        const params = new HttpParams({
            fromObject: {
                slave_id: this.id,
                root: this.root,
                path: this.path,
            }
        })
        return params.toString()
    }
}
export interface ListResponse {
    dir: Dir
    items: Array<FileInfo>
}
export class Source {
    public readonly textTracks = new Array<FileUrl>()
    constructor(private readonly access: string,
        public readonly path: Path,
        public readonly source: FileInfo) {
    }
    addTrack(url: FileUrl) {
        this.textTracks.push(url)
    }
    private url_: string | undefined
    get url(): string {
        let url = this.url_
        if (url === undefined) {
            const path = this.path
            const params = new HttpParams({
                fromObject: {
                    slave_id: path.id,
                    root: path.root,
                    path: this.source.filename,
                    access_token: this.access,
                }
            })
            url = ServerAPI.forward.v1.fs.httpURL('download_access') + '?' + params.toString()
            this.url_ = url
        }
        return url
    }
    get id(): string {
        const path = this.path
        const params = new HttpParams({
            fromObject: {
                slave_id: path.id,
                root: path.root,
                path: this.source.filename,
            }
        })
        return params.toString()
    }
}
export class FileUrl {
    constructor(
        private readonly access: string,
        public readonly label: string,
        public readonly id: string, public readonly root: string,
        public readonly filepath: string,
        public readonly isDefault: boolean,
    ) { }

    get url(): string {
        const params = new HttpParams({
            fromObject: {
                slave_id: this.id,
                root: this.root,
                path: this.filepath,
                access_token: this.access,
            }
        })
        return ServerAPI.forward.v1.fs.httpURL('download_access') + '?' + params.toString()
    }
}
export class Current {
    currentTime = 0

    first = true
    save = false
    constructor(public readonly source: Source) { }
    get name(): string {
        return this.source.source.name
    }
    get url(): string {
        return this.source.url
    }
    get id(): string {
        return this.source.id
    }
}
interface Write {
    id: string
    currentTime: number
}
export class Manager {
    private items_ = new Array<Source>()
    get items(): Array<Source> {
        return this.items_
    }
    private ch_ = new Channel<string>(1)
    private chWrite_ = new Channel<Write>(1)
    constructor(
        private readonly video: HTMLVideoElement,
        private readonly access: string,
        private readonly httpClient: HttpClient,
        public readonly path: Path,
    ) { }
    close() {
        this.ch_.close()
        this.chWrite_.close()
    }
    async run(): Promise<void> {
        const path = this.path
        const resp = await ServerAPI.forward.v1.fs.child('list').get<ListResponse>(this.httpClient, {
            params: {
                slave_id: path.id,
                root: path.root,
                path: path.path,
            },
        }).toPromise()
        const dir = resp.dir
        const items = new Array(resp.items.length)
        for (let i = 0; i < resp.items.length; i++) {
            items[i] = new FileInfo(dir.root, dir.dir, resp.items[i])
        }
        items.sort(FileInfo.compare)
        for (let i = 0; i < items.length; i++) {
            const fileinfo = items[i]
            if (fileinfo.isVideo) {
                const source = new Source(this.access, path, fileinfo)
                this.items_.push(source)
            }
        }
        for (let i = 0; i < items.length; i++) {
            const fileinfo = items[i]
            if (!fileinfo.isText || fileinfo.ext.toLowerCase() != '.vtt') {
                continue
            }
            const basename = fileinfo.basename
            const arrs = this.items_
            for (let index = 0; index < arrs.length; index++) {
                const element = arrs[index]
                const name = element.source.basename
                if (basename.startsWith(name)) {
                    let label = basename.substring(name.length)
                    if (label.startsWith('.')) {
                        label = label.substring(1)
                    }
                    if (label == '') {
                        label = (element.textTracks.length + 1).toString()
                    }
                    element.addTrack(new FileUrl(this.access, label, path.id, path.root, fileinfo.filename, element.textTracks.length == 0))
                }
            }
        }
        this._run()
        if (DB.isSupported) {
            this._runWrite()
        }
    }
    push(id: string) {
        const ch = this.ch_
        while (!ch.trySend(id)) {
            ch.tryRecv()
        }
    }
    private async _run() {
        const ch = this.ch_
        if (ch.isClosed) {
            return
        }
        while (true) {
            const value = await ch.recv()
            if (value.done) {
                await this._play(value.value)
            } else {
                break
            }
        }
    }
    private completer_: Completer<string> | undefined
    private async _getLastDBName(): Promise<string> {
        let completer = this.completer_
        if (completer) {
            return completer.promise
        }
        completer = new Completer<string>()
        try {
            let name = ''
            if (DB.isSupported) {
                name = await DB.instance.getCurrentName(this.path.key)
            }
            completer.resolve(name)
        } catch (e) {
            this.completer_ = undefined
            completer.reject(e)
        }
        return completer.promise
    }
    private async _play(name: string) {
        const items = this.items_
        let source: Source | undefined
        let start = 0
        if (name == '') {
            try {
                name = await this._getLastDBName()
            } catch (e) {
                console.log(e)
            }
            if (items.length != 0) {
                source = items[0]
                start = 1
            }
        }

        for (let i = start; i < items.length; i++) {
            const element = items[i]
            if (element.source.name == name) {
                source = element
                break
            }
        }

        if (source) {
            await this._playSource(source)
        }
    }

    next(): string | undefined {
        const items = this.items_
        if (items.length == 0) {
            return
        }
        const current = this.current
        if (!current) {
            return
        }
        const name = current.name
        let found = 0
        for (let i = 0; i < items.length; i++) {
            const item = items[i]
            if (item.source.name == name) {
                found = i + 1
                break
            }
        }
        if (found == items.length) { // list end
            // found = 0
            return
        }
        return items[found].source.name
    }
    private async _playSource(source: Source) {
        if (source == this.current_?.source) {
            return
        }
        console.log('play', source.source.name)
        const current = new Current(source)
        this.current_ = current
        try {
            if (DB.isSupported) {
                try {
                    const currentTime = await DB.instance.getCurrentTime(source.id)
                    if (this.current_ == current) {
                        current.currentTime = currentTime
                        console.log('last time', currentTime)
                    }
                } catch (e) {
                    console.log(e)
                }
            }
        } catch (e) {
            console.log(e)
        }
    }
    private current_: Current | undefined
    get current(): Current | undefined {
        return this.current_
    }
    async timeupdate() {
        const current = this.current_
        if (!current) {
            return
        }
        // restore progress
        if (current.currentTime > 0 && current.first) {
            current.first = false
            if (current.currentTime > 0) {
                this._skipTo(current)
                return
            }
        }
        if (DB.isSupported) {
            const player = this.video
            // 進度提前 2 秒，以播放一段觀看過的內容以便 銜接喚醒 用戶觀看記憶
            const wait = 2
            const currentTime = player.currentTime - wait
            if (currentTime > 2 && currentTime + 2 + wait < player.duration) {
                // 每 2 秒 記錄進度
                const last = this.last_
                if (currentTime < last || currentTime > last + wait) {
                    this.last_ = currentTime
                    this.chWrite_.trySend({
                        id: current.id,
                        currentTime: currentTime,
                    })
                }
            }
        }
    }
    private last_ = 0
    private async _runWrite() {
        const ch = this.chWrite_
        if (ch.isClosed) {
            return
        }
        while (true) {
            const value = await ch.recv()
            if (value.done) {
                await this._write(value.value)
            } else {
                break
            }
        }
    }
    private async _write(w: Write) {
        const id = w.id
        try {
            await DB.instance.putCurrentTime(id, w.currentTime)
        } catch (e) {
            console.log(e)
        }
    }
    private save_ = false
    async saveName() {
        const current = this.current_
        if (!current || this.save_) {
            return
        }
        this.save_ = true

        const path = this.path
        const key = path.key
        try {
            await DB.instance.putCurrentName(key, current.name)
        } catch (e) {
            console.log(e)
        } finally {
            this.save_ = false
        }
    }
    private _skipTo(current: Current) {
        if (this.current_ != current) {
            return
        }
        const player = this.video
        const skipTo = current.currentTime
        if (skipTo > player.currentTime + 1 && skipTo + 1 < player.duration) {
            if (player.fastSeek) {
                player.fastSeek(skipTo)
            } else {
                player.currentTime = skipTo
            }
        }
    }
}