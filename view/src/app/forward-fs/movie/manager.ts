import { HttpClient, HttpParams } from "@angular/common/http"
import { ServerAPI } from "src/app/core/core/api"
import { Channel } from "src/app/core/utils/completer"
import { FileInfo, Dir } from '../fs';
import { VideoJsPlayer } from "video.js"
export class Path {
    constructor(public readonly id: string,
        public readonly root: string,
        public readonly path: string,
    ) { }
    equal(other: Path): boolean {
        return this.id == other.id &&
            this.root == other.root &&
            this.path == other.path
    }
}
export interface ListResponse {
    dir: Dir
    items: Array<FileInfo>
}
const mimeExt = new Map<string, string | undefined>()
mimeExt.set('.webm', 'video/webm')
mimeExt.set('.mp4', 'video/mp4')
mimeExt.set('.m4v', undefined)
mimeExt.set('.mov', 'video/quicktime')
mimeExt.set('.avi', 'video/x-msvideo')
mimeExt.set('.flv', 'video/x-flv')
mimeExt.set('.wmv', 'video/x-ms-wmv')
mimeExt.set('.asf', 'video/x-ms-asf')
mimeExt.set('.mpeg', 'video/mpeg')
mimeExt.set('.mpg', 'video/mpeg')
mimeExt.set('.vob', undefined)
mimeExt.set('.mkv', 'video/x-matroska')
mimeExt.set('.rm', 'audio/x-pn-realaudio')
mimeExt.set('.rmvb', undefined)
class Source {
    public readonly textTracks = new Array<FileUrl>()
    constructor(private readonly access: string,
        private readonly path: Path,
        public readonly source: FileInfo) {
    }
    addTrack(url: FileUrl) {
        this.textTracks.push(url)
    }
    get url(): string {
        const path = this.path
        const params = new HttpParams({
            fromObject: {
                slave_id: path.id,
                root: path.root,
                path: this.source.filename,
                access_token: this.access,
            }
        })
        return ServerAPI.forward.v1.fs.httpURL('download_access') + '?' + params.toString()
    }
    get mime(): string | undefined {
        const ext = this.source.ext.toLowerCase()
        return mimeExt.get(ext)
    }
}
class FileUrl {
    constructor(
        private readonly access: string,
        public readonly label: string,
        public readonly id: string, public readonly root: string,
        public readonly filepath: string,
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
export class Manager {
    private items_ = new Array<Source>()
    private ch_ = new Channel<string>(1)
    constructor(private player: VideoJsPlayer,
        private readonly access: string,
        private readonly httpClient: HttpClient,
        public readonly path: Path,
    ) { }
    close() {
        this.ch_.close()
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
                    element.addTrack(new FileUrl(this.access, label, path.id, path.root, fileinfo.filename))
                }
            }
        }
        this._run()
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
    private async _play(name: string) {
        const items = this.items_
        let source: Source | undefined
        if (name == '') {
            if (items.length == 0) {
                return
            }
            source = items[0]
        } else {
            for (let i = 0; i < items.length; i++) {
                const element = items[i]
                if (element.source.name == name) {
                    source = element
                    break
                }
            }
        }
        if (source) {
            await this._playSource(source)
        }
    }

    next() {
        const items = this.items_
        if (items.length == 0) {
            return
        }
        const name = this.current_
        let found = 0
        for (let i = 0; i < items.length; i++) {
            const item = items[i]
            if (item.source.name == name) {
                found = i + 1
                break
            }
        }
        if (found == items.length) {
            found = 0
        }
        this.push(items[found].source.name)
    }
    private async _playSource(source: Source) {
        const player = this.player
        player.src({
            src: source.url,
            type: source.mime,
        })
        const add = source.textTracks
        for (let i = 0; i < add.length; i++) {
            const element = add[i]
            player.addRemoteTextTrack({
                kind: 'captions',
                src: element.url,
                default: i == 0,
                label: element.label == "" ? (i + 1).toString() : element.label,
            }, false)
        }
        console.log('play', source.source.name)
        this.current_ = source.source.name
        try {
            await player.play()
        } catch (e) {
            console.log(e)
        }
    }
    private current_ = ''
}