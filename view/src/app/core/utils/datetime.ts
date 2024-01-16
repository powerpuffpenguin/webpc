export function usedString(used: number): string {
    const strs: Array<string> = []
    if (used >= 86400) {
        strs.push(`${Math.floor(used / 86400)} days`)
        used %= 86400
    }
    if (used >= 3600) {
        strs.push(`${Math.floor(used / 3600)} hours`)
        used %= 3600
    }
    if (used >= 60) {
        strs.push(`${Math.floor(used / 60)} minutes`)
        used %= 60
    }
    if (used >= 0) {
        strs.push(`${used} seconds`)
    }
    return strs.join(' ')
}
export function datetimeString(date: Date): string {
    const y = date.getFullYear().toString()
    const m = (date.getMonth() + 1).toString().padStart(2, '0')
    const d = date.getDate().toString().padStart(2, '0')
    const hh = date.getHours().toString().padStart(2, '0')
    const mm = date.getMinutes().toString().padStart(2, '0')
    const ss = date.getSeconds().toString().padStart(2, '0')
    return `${y}-${m}-${d} ${hh}:${mm}:${ss}`
}