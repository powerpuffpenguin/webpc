import { hash } from 'spark-md5'
export function md5String(message: string): string {
    return hash(message)
}