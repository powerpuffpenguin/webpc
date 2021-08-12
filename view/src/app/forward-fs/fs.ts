import { HttpParams } from '@angular/common/http';
import { ServerAPI } from '../core/core/api';
import { sizeString } from '../core/utils/utils';

const audios = new Set();
([
    '.mp3',
    '.aac',
    '.flac',
    '.ape',
    '.wav',
    '.ogg',
]).forEach(function (str: string) {
    audios.add(str)
});
const videos = new Set();
([
    '.webm',
    '.mp4', '.m4v',
    '.mov',
    '.avi',
    '.flv',
    '.wmv', '.asf',
    '.mpeg', '.mpg', '.vob',
    '.mkv',
    '.rm', '.rmvb',
]).forEach(function (str: string) {
    videos.add(str)
});
const images = new Set();
([
    '.gif',
    '.jpeg', '.jpg',
    '.bmp',
    '.png',
    '.svg', '.ico',
    '.webp',

]).forEach(function (str: string) {
    images.add(str)
});
const texts = new Set();
([
    '.txt', '.text',
    '.json', '.xml', '.yaml', '.ini', '.html', '.md',
    '.sh',
    '.bat', '.cmd', '.vbs',
    '.go', '.dart', '.py', '.py3', '.js', '.ts', '.c', '.cc', '.h', '.hpp', '.cpp',
    '.java', '.php', '.lua', '.jsonnet', '.libsonnet',
]).forEach(function (str: string) {
    texts.add(str)
});
export enum FileType {
    Dir,
    Video,
    Audio,
    Image,
    Text,
    Binary,
}
export interface Dir {
    root: string
    read: boolean
    write: boolean
    shared: boolean
    dir: string
    id: string
}
export class FileInfo {
    name: string
    mode: number
    size: number
    isDir: boolean

    filename: string
    checked = false
    private _filetype = FileType.Binary

    constructor(public root: string, public dir: string, other: FileInfo) {
        this.name = other.name
        this.mode = other.mode
        this.size = other.size
        this.isDir = other.isDir
        if (dir.endsWith('/')) {
            this.filename = dir + other.name
        } else {
            this.filename = dir + '/' + other.name
        }

        if (this.isDir) {
            this._filetype = FileType.Dir
        } else {
            this._settype()
        }
    }
    private _settype() {
        const ext = this.ext.toLowerCase()
        if (videos.has(ext)) {
            this._filetype = FileType.Video
        } else if (audios.has(ext)) {
            this._filetype = FileType.Audio
        } else if (images.has(ext)) {
            this._filetype = FileType.Image
        } else if (texts.has(ext)) {
            this._filetype = FileType.Text
        } else {
            this._filetype = FileType.Binary
        }
    }
    setName(name: string) {
        this.name = name
        const dir = this.dir
        if (dir.endsWith('/')) {
            this.filename = dir + name
        } else {
            this.filename = dir + '/' + name
        }
        if (!this.isDir) {
            this._settype()
        }
    }
    static compare(l: FileInfo, r: FileInfo): number {
        let lv = 0
        if (l.isDir) {
            lv = 1
        }
        let rv = 0
        if (r.isDir) {
            rv = 1
        }
        if (lv == rv) {
            if (l.name == r.name) {
                return 0
            }
            return l.name < r.name ? -1 : 1
        }
        return rv - lv
    }
    get ext(): string {
        const index = this.name.lastIndexOf('.')
        if (index == -1) {
            return ''
        }
        return this.name.substring(index)
    }
    get filetype(): FileType {
        return this._filetype
    }
    get isSupportUncompress(): boolean {
        if (this.isDir) {
            return false
        }
        const name = this.name.toLowerCase()
        return name.endsWith(`.tar.gz`) || name.endsWith(`.tar.bz2`) || name.endsWith(`.tar`) || name.endsWith(`.zip`)
    }
    get url(): string {
        switch (this._filetype) {
            case FileType.Dir:
                return '/forward/fs/list'
            case FileType.Video:
                return '/forward/fs/view/video'
            case FileType.Audio:
                return '/forward/fs/view/audio'
            case FileType.Image:
                return '/forward/fs/view/image'
            case FileType.Text:
                return '/forward/fs/view/text'
        }
        return ''
    }
    downloadURL(id: string, root: string, path: string, access: string): string {
        if (this.isDir) {
            return ''
        }
        const params = new HttpParams({
            fromObject: {
                slave_id: id,
                root: root,
                path: path,
                access_token: access,
            }
        })
        return ServerAPI.forward.v1.fs.httpURL('download') + '?' + params.toString()
    }
    get icon(): string {
        switch (this.filetype) {
            case FileType.Dir:
                return 'folder'
            case FileType.Video:
                return 'movie_creation'
            case FileType.Audio:
                return 'audiotrack'
            case FileType.Image:
                return 'insert_photo'
            case FileType.Text:
                return 'event_note'
        }
        if (this.isSupportUncompress) {
            return 'unarchive'
        }
        return 'insert_drive_file'
    }
    get modeString(): string {
        const str = "dalTLDpSugct?"
        let w = 0
        const buf = new Array<string>(32)
        const m = this.mode || 0
        for (let i = 0; i < str.length; i++) {
            const c = str[i]
            if ((m & (1 << (32 - 1 - i))) != 0) {
                buf[w] = c
                w++
            }
        }
        if (w == 0) {
            buf[w] = '-'
            w++
        }
        const rwx = "rwxrwxrwx"
        for (let i = 0; i < rwx.length; i++) {
            const c = rwx[i]
            if ((m & (1 << (9 - 1 - i))) != 0) {
                buf[w] = c
            } else {
                buf[w] = '-'
            }
            w++
        }
        return buf.slice(0, w).join('')
    }
    get sizeString(): string {
        if (this.isDir) {
            return ''
        }
        return sizeString(this.size)
    }
}
export interface DirName {
    name: string
    dir: string
}
export function split(name: string): DirName {
    if (typeof name !== "string") {
        return {
            name: '',
            dir: '',
        }
    }
    const index = name.lastIndexOf('/')
    if (index != -1) {
        return {
            name: name.substring(index + 1),
            dir: name.substring(0, index),
        }
    }
    return {
        name: name,
        dir: '',
    }
}