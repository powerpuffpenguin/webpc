import { Completer } from "src/app/core/utils/completer"
const version = 2
export class DB {
    static readonly instance = new DB()
    static get isSupported(): boolean {
        if (typeof window.indexedDB === "object" && typeof window.indexedDB.open === "function") {
            return true
        } else {
            return false
        }
    }
    private db_: Completer<IDBDatabase> | undefined
    async getDB(): Promise<IDBDatabase> {
        if (this.db_) {
            return this.db_.promise
        }
        let completer = new Completer<IDBDatabase>()
        this.db_ = completer
        try {
            const request = window.indexedDB.open(`videojs`, version)
            request.onupgradeneeded = (_) => {
                const db = request.result
                if (!db.objectStoreNames.contains('currentTime')) {
                    db.createObjectStore(
                        'currentTime',
                        {
                            keyPath: 'id',
                            autoIncrement: false,
                        }
                    )
                }
                if (!db.objectStoreNames.contains('currentName')) {
                    db.createObjectStore(
                        'currentName',
                        {
                            keyPath: 'id',
                            autoIncrement: false,
                        }
                    )
                }
            };
            request.onsuccess = (_) => {
                completer.resolve(request.result)
            }
            request.onerror = (e) => {
                completer.reject(e)
            }

        } catch (e) {
            this.db_ = undefined
            completer.reject(e)
        }
        return completer.promise
    }
    async putCurrentTime(id: string, val: number) {
        const db = await this.getDB()
        await new Promise<void>((resolve, reject) => {
            try {
                // 增加數據
                const request = db.transaction(['currentTime'], 'readwrite')
                    .objectStore('currentTime')
                    .put({
                        id: id,
                        val: val,
                    });
                request.onsuccess = (evt) => {
                    resolve()
                }
                request.onerror = (e) => {
                    reject(e)
                }
            } catch (e) {
                reject(e)
            }
        })
    }
    async getCurrentTime(id: string): Promise<number> {
        const db = await this.getDB()
        return new Promise<number>((resolve, reject) => {
            try {
                // 增加數據
                const request = db.transaction(['currentTime'], 'readonly')
                    .objectStore('currentTime')
                    .get(id)
                request.onsuccess = (evt) => {
                    if (request.result && typeof request.result.val === "number") {
                        resolve(request.result.val)
                    } else {
                        resolve(0)
                    }
                }
                request.onerror = (e) => {
                    reject(e)
                }
            } catch (e) {
                reject(e)
            }
        })
    }
    async putCurrentName(id: string, name: string) {
        const db = await this.getDB()
        await new Promise<void>((resolve, reject) => {
            try {
                // 增加數據
                const request = db.transaction(['currentName'], 'readwrite')
                    .objectStore('currentName')
                    .put({
                        id: id,
                        name: name,
                    });
                request.onsuccess = (evt) => {
                    resolve()
                }
                request.onerror = (e) => {
                    reject(e)
                }
            } catch (e) {
                reject(e)
            }
        })
    }
    async getCurrentName(id: string): Promise<string> {
        const db = await this.getDB()
        return new Promise<string>((resolve, reject) => {
            try {
                // 增加數據
                const request = db.transaction(['currentName'], 'readonly')
                    .objectStore('currentName')
                    .get(id)
                request.onsuccess = (evt) => {
                    if (request.result && typeof request.result.name === "string") {
                        resolve(request.result.name)
                    } else {
                        resolve("")
                    }
                }
                request.onerror = (e) => {
                    reject(e)
                }
            } catch (e) {
                reject(e)
            }
        })
    }
}