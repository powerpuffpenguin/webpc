export function getUnix(): number {
    const dateTime = Date.now()
    return Math.floor(dateTime / 1000)
}
export function durationString(d: any) {
    const result = new Array<string>()
    const days = Math.floor(d.asDays())
    if (days > 0) {
        result.push(`${days} days`)
    }
    const hours = Math.floor(d.asHours()) % 24
    if (hours > 0) {
        result.push(`${hours} hours`)
    }
    const minutes = Math.floor(d.asMinutes()) % 60
    if (minutes > 0) {
        result.push(`${minutes} minutes`)
    }
    const seconds = Math.floor(d.asSeconds()) % 60
    if (seconds > 0) {
        result.push(`${seconds} seconds`)
    }
    return result.join(' ')
}