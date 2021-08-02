import { getItem as nativeGetItem, setItem as nativeSetItem } from './local-storage'
import { environment } from 'src/environments/environment';
import { aesDecrypt, aesEncrypt } from './aes';

export function getItem(key: string, def?: string): string | null {
    const result = nativeGetItem(key)
    if (typeof def === "string" && typeof result !== "string") {
        return def
    }
    if (typeof result !== "string") {
        return null
    }
    const obj = JSON.parse(result ?? '{}')
    const at = obj.at
    return aesDecrypt(obj.value, `${environment.aesKey}.${at}`, `${environment.aesIV}.${at}`)
}
export function setItem(key: string, value: string) {
    const at = new Date().toUTCString()
    value = aesEncrypt(value, `${environment.aesKey}.${at}`, `${environment.aesIV}.${at}`)
    nativeSetItem(key,
        JSON.stringify({
            value: value,
            at: at,
        }),
    )
}
