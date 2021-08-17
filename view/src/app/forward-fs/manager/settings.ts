import { getItem, removeItem, setItem } from "src/app/core/utils/local-storage"
import { environment } from "src/environments/environment"
const Key = 'fs.clipboard'

export class Settings {
    private static instance_: Settings | undefined
    static get instance(): Settings {
        if (!this.instance_) {
            this.instance_ = new Settings()
        }
        return this.instance_
    }
    private constructor() { }
    ctrl = false
    shift = false
    all = false

    getClipboard(): Clipboard | undefined {
        const jsonStr = getItem(Key)
        if (typeof jsonStr === "string") {
            try {
                return Clipboard.fromJSON(jsonStr)
            } catch (e) {
                if (!environment.production) {
                    console.log('getClipboard error:', e)
                }
            }
        }
        return
    }
    setClipboard(obj: Clipboard) {
        setItem(Key, obj.toJSON())
    }
    removeClipboard(obj: Clipboard) {
        const jsonStr = getItem(Key)
        if (typeof jsonStr === "string") {
            try {
                const o = Clipboard.fromJSON(jsonStr)
                if (o.equalValue(obj)) {
                    removeItem(Key)
                }
            } catch (e) {
            }
        }
    }
}
export class Clipboard {
    constructor(
        public readonly id: string,
        public readonly root: string,
        public readonly dir: string,
        public readonly names: Array<string>,
        public readonly copied: boolean,
    ) { }
    static fromJSON(jsonStr: string): Clipboard {

        const obj = JSON.parse(jsonStr)
        if (typeof obj.at !== "number" ||
            typeof obj.id !== "string" ||
            typeof obj.root !== "string" ||
            typeof obj.dir !== "string" ||
            !Array.isArray(obj.names) ||
            obj.names.length == 0
        ) {
            throw new Error('jsonStr not match ')
        }
        if ((Date.now() - obj.at) / 1000 > 3600) {
            throw new Error('jsonStr expired')
        }
        for (let i = 0; i < obj.names.length; i++) {
            const name = obj.names[i]
            if (typeof name !== "string") {
                throw new Error('jsonStr not match ')
            }
        }
        const copied = obj.copied ? true : false
        return new Clipboard(obj.id, obj.root, obj.dir, obj.names, copied)
    }
    toJSON(): string {
        return JSON.stringify({
            id: this.id,
            root: this.root,
            dir: this.dir,
            names: this.names,
            copied: this.copied,
            at: Date.now(),
        })
    }
    equalValue(o: Clipboard): boolean {
        if (this.id != o.id ||
            this.root != o.root ||
            this.dir != o.dir ||
            this.names.length != o.names.length ||
            this.copied != o.copied) {
            return false
        }
        for (let i = 0; i < this.names.length; i++) {
            if (this.names[i] != o.names[i]) {
                return false
            }
        }
        return true
    }
}