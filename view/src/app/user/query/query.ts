import { Params } from "@angular/router"
import { environment } from 'src/environments/environment';
export const DefaultLimit = environment.production ? 25 : 5
export enum Result {
    UNIVERSAL = 0,
    DATA = 1,
    COUNT = 2,
    DATA_COUNT = 3,
}
function parseParamsInt(params: Params, key: string): number {

    const v = params[key]
    if (typeof v === "string") {
        try {
            let val = parseInt(v)
            if (isNaN(val) || val < 0 || !isFinite(val)) {
                return -1
            }
            return Math.floor(val)
        } catch (e) {
            console.warn(`parseParamsInt ${key} error : `, e)
        }
    }
    return -1
}

export class Request {
    count = -1
    limit = DefaultLimit
    offset = 0
    name: string | undefined = ''
    nameFuzzy = true
    cloneTo(other: Request) {
        if (this == other) {
            return
        }
        other.count = this.count
        other.limit = this.limit
        other.offset = this.offset
        other.name = this.name
        other.nameFuzzy = this.nameFuzzy
    }
    constructor(params?: Params) {
        if (params) {
            const limit = parseParamsInt(params, 'limit')
            if (limit > 0 && limit <= 100) {
                this.limit = limit
            }
            const offset = parseParamsInt(params, 'offset')
            if (offset > 0) {
                this.offset = offset
            }
            const name = params['name']
            if (typeof name === "string") {
                this.name = name.trim()
            }
            if (params['nameFuzzy'] == "false") {
                this.nameFuzzy = false
            }
        }
    }
    toArgs(): {
        [param: string]: string | string[];
    } {
        return {
            result: (this.count < 0 ? Result.DATA_COUNT : Result.DATA).toString(),
            limit: this.limit.toString(),
            offset: this.offset.toString(),
            name: (this.name ?? '').trim(),
            nameFuzzy: this.nameFuzzy.toString(),
        }
    }
    toQuery(): Params {
        return {
            limit: this.limit,
            offset: this.offset,
            name: this.name,
            nameFuzzy: this.nameFuzzy,
        }
    }
}

export interface Data {
    id: string
    name: string
    nickname: string
    authorization: Array<number>
}
export interface Response {
    result: 'DATA' | 'COUNT' | 'DATA_COUNT'
    data: Array<any>
    count: string
}