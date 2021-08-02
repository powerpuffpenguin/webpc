import { environment } from 'src/environments/environment';
import { ModeOfOperation, utils } from 'aes-js';
import { md5String } from './md5';

const Key = utils.hex.toBytes(md5String(environment.aesKey))
const IV = utils.hex.toBytes(md5String(environment.aesIV))
export function aesDecrypt(data: string, key?: string, iv?: string): string {
    let keySource = Key
    let ivSource = IV
    if (typeof key === "string" && key.length > 0) {
        keySource = utils.hex.toBytes(md5String(key))
    }
    if (typeof iv === "string" && iv.length > 0) {
        ivSource = utils.hex.toBytes(md5String(iv))
    }
    const aes = new ModeOfOperation.ofb(keySource, ivSource)
    return utils.utf8.fromBytes(aes.decrypt(utils.hex.toBytes(data)))
}
export function aesEncrypt(data: string, key?: string, iv?: string): string {
    let keySource = Key
    let ivSource = IV
    if (typeof key === "string" && key.length > 0) {
        keySource = utils.hex.toBytes(md5String(key))
    }
    if (typeof iv === "string" && iv.length > 0) {
        ivSource = utils.hex.toBytes(md5String(iv))
    }
    const aes = new ModeOfOperation.ofb(keySource, ivSource)
    return utils.hex.fromBytes(aes.encrypt(utils.utf8.toBytes(data)))
}