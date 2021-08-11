import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { Authorization, ServerAPI } from '../core/api';
import { Completer } from '../utils/completer';
import { removeItem } from "../utils/local-storage";
import { getItem, setItem } from "../utils/aes-local-storage";
import { getUnix } from '../utils/utils';
import { md5String } from '../utils/md5';
import { aesDecrypt, aesEncrypt } from '../utils/aes';
import { Codes, NetError } from '../core/restful';
const Key = 'session'
const Platform = 'web'
export interface Userdata {
    readonly id: string
    readonly name?: string
    readonly nickname?: string
    readonly authorization?: Array<number>
}
export class Session {
    constructor(public access: string, public refresh: string, public userdata: Userdata) {
    }
    get who(): string {
        if (!this.userdata || !this.userdata.id) {
            return ''
        }
        let name = this.userdata.name ?? ''
        const nickname = this.userdata.nickname ?? ''
        return nickname.length == 0 ? name : `${nickname} [${name}]`
    }
    get root(): boolean {
        return this.anyAuth(Authorization.Root)
    }
    get write(): boolean {
        return this.anyAuth(Authorization.Write)
    }
    get read(): boolean {
        return this.anyAuth(Authorization.Read)
    }
    get authorization(): Array<number> {
        if (!this.userdata || !this.userdata.authorization || !Array.isArray(this.userdata.authorization)) {
            return []
        }
        return this.userdata.authorization
    }
    /**
     * if has all authorization return true
     */
    testAuth(...vals: Array<number>): boolean {
        if (!this.userdata || !this.userdata.id) {
            return false
        }
        let found: boolean
        const authorization = this.authorization
        for (let i = 0; i < vals.length; i++) {
            found = false
            const val = vals[i]
            for (let j = 0; j < authorization.length; j++) {
                if (val == authorization[j]) {
                    found = true
                    break
                }
            }
            if (!found) {
                return false
            }
        }
        return true
    }
    /**
     * if not has any authorization return true
     */
    noneAuth(...vals: Array<number>): boolean {
        if (!this.userdata || !this.userdata.id) {
            return false
        }
        const authorization = this.authorization
        for (let i = 0; i < authorization.length; i++) {
            for (let j = 0; j < vals.length; j++) {
                const val = vals[j]
                if (authorization[i] == val) {
                    return false
                }
            }
        }
        return true
    }
    /**
     * if has any authorization return true
     */
    anyAuth(...vals: Array<number>): boolean {
        if (!this.userdata || !this.userdata.id) {
            return false
        }
        const authorization = this.authorization
        for (let i = 0; i < authorization.length; i++) {
            for (let j = 0; j < vals.length; j++) {
                const val = vals[j]
                if (authorization[i] == val) {
                    return true
                }
            }
        }
        return false
    }
}
interface Store {
    userdata: Userdata
    access: string
    refresh: string
}
interface SigninResponse {
    access: string
    refresh: string
    id: string
    name: string
    nickname: string
    authorization: Array<number>
}
interface RefreshResponse {
    access: string
    refresh: string
}
export class Manager {
    static instance_ = new Manager()
    static get instance(): Manager {
        return Manager.instance_
    }
    private constructor() {
    }
    private remember_ = false
    get session(): Session | undefined {
        return this.subject_.value
    }
    private readonly subject_ = new BehaviorSubject<Session | undefined>(undefined)
    get observable(): Observable<Session | undefined> {
        return this.subject_
    }
    private _load(): Session | undefined {
        const str = getItem(Key)
        if (typeof str !== "string") {
            return
        }
        try {
            const obj: Store = JSON.parse(aesDecrypt(str))
            if (obj !== null && typeof obj === "object") {
                const access = obj.access
                const refresh = obj.refresh
                const userdata = obj.userdata
                if (typeof access === "string" && access.length > 0
                    && typeof refresh === "string" && refresh.length > 0
                    && userdata !== null && typeof userdata === "object" && userdata.id) {
                    this.remember_ = true
                    return new Session(access, refresh, userdata)
                }
            }
        } catch (e) {
            console.warn(`load token error`, e)
        }
        return
    }
    load() {
        this.subject_.next(this._load())
    }
    private _save(session: Session) {
        try {
            const data = JSON.stringify({
                userdata: session.userdata,
                access: session.access,
                refresh: session.refresh,
            })
            console.log(`save token`, data)
            setItem(Key, aesEncrypt(data))
        } catch (e) {
            console.log('save token error', e)
        }
    }
    refresh_: Completer<Session | undefined> | undefined
    private readonly signining_ = new BehaviorSubject<boolean>(false)
    get signining(): Observable<boolean> {
        return this.signining_
    }
    async signin(httpClient: HttpClient,
        name: string, password: string, remember: boolean,
    ): Promise<Session | undefined> {
        if (this.signining_.value) {
            console.warn('wait signing completed')
            return
        }
        this.signining_.next(true)
        this.remember_ = remember
        let completer: Completer<Session | undefined> | undefined
        let session: Session | undefined
        try {
            // wait refresh completed
            while (this.refresh_) {
                const completer = this.refresh_
                try {
                    await completer.promise
                } catch (error) {
                }
                if (completer == this.refresh_) {
                    this.refresh_ = undefined
                }
            }
            completer = new Completer<Session | undefined>()
            this.refresh_ = completer
            const unix = getUnix()
            password = md5String(password)
            password = md5String(`${Platform}.${password}.${unix}`)
            const response = await ServerAPI.v1.sessions.post<SigninResponse>(httpClient,
                {
                    platform: Platform,
                    name: name,
                    password: password,
                    unix: unix,
                },
                {
                    headers: {
                        'Interceptor': 'none',
                    },
                },
            ).toPromise()
            session = new Session(response.access, response.refresh, {
                id: response.id,
                name: response.name,
                nickname: response.nickname,
                authorization: response.authorization,
            })
            if (remember) {
                this._save(session)
            }
            this.subject_.next(session)
        } finally {
            if (completer) {
                completer.resolve(session)
                if (completer == this.refresh_) {
                    this.refresh_ = undefined
                }
            }
            this.signining_.next(false)
        }
        return
    }
    signout(httpClient: HttpClient) {
        const session = this.subject_.value
        if (session) {
            if (this.remember_) {
                removeItem(Key)
            }
            this.subject_.next(undefined)
            ServerAPI.v1.sessions.child('access').delete(httpClient, {
                headers: {
                    Interceptor: 'none',
                    Authorization: `Bearer ${session.access}`
                }
            }).toPromise().then(() => {
                console.info(`signout who=${session.who}`)
            }, (e) => {
                console.warn(`signout who=${session.who} error=${e}`)
            })
        }
    }
    async refresh(httpClient: HttpClient, session: Session, err?: NetError): Promise<Session | undefined> {
        if (this.refresh_) { // refreshing
            return this.refresh_.promise
        }

        const current = this.subject_.value
        if (!current) { // already signout
            return
        } else if (session != current) { // already refresh
            return current
        }
        if (err && err.grpc != Codes.Unauthenticated) {
            throw err
        }

        // refresh
        const completer = new Completer<Session | undefined>()
        this.refresh_ = completer
        ServerAPI.v1.sessions.child('refresh').post<RefreshResponse>(httpClient,
            {
                access: session.access,
                refresh: session.refresh,
            },
            {
                headers: {
                    Interceptor: 'none',
                }
            },
        ).toPromise().then((resp) => {
            const s = new Session(resp.access, resp.refresh, session.userdata)
            if (this.remember_) {
                this._save(s)
            }
            this.subject_.next(s)
            completer.resolve(s)
        }, (e) => {
            completer.reject(e)
        }).finally(() => {
            this.refresh_ = undefined
        })
        return completer.promise
    }
    clear(session: Session) {
        if (session == this.subject_.value) {
            this.subject_.next(undefined)
            if (this.remember_) {
                const current = this._load()
                if (current && current.access == session.access) {
                    removeItem(Key)
                }
            }
        }
    }
}