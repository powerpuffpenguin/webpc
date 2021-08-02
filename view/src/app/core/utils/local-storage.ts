
class Bridge {
    private static instance_: Bridge | undefined
    static get instance(): Bridge {
        if (!Bridge.instance_) {
            Bridge.instance_ = new Bridge()
        }
        return Bridge.instance_
    }
    private keys_: Map<string, string> | undefined
    private constructor() {
        if (typeof localStorage === undefined) {
            this.keys_ = new Map<string, string>()
        }
    }
    getItem(key: string): string | null {
        if (this.keys_) {
            const result = this.keys_.get(key)
            if (result === undefined) {
                return null
            }
            return result
        }
        return localStorage.getItem(key)
    }
    setItem(key: string, value: string) {
        if (this.keys_) {
            this.keys_.set(key, value)
        } else {
            localStorage.setItem(key, value)
        }
    }
    removeItem(key: string) {
        if (this.keys_) {
            this.keys_.delete(key)
        } else {
            localStorage.removeItem(key)
        }
    }
    clear() {
        if (this.keys_) {
            this.keys_.clear()
        } else {
            localStorage.clear()
        }
    }
}
export function expired(key: string, value: string) {
    const instance = Bridge.instance
    if (instance.getItem(key) === value) {
        instance.removeItem(key)
    }
}
export function getItem(key: string, def?: string): string | null {
    const result = Bridge.instance.getItem(key)
    if (typeof def === "string" && result === null) {
        return def
    }
    return Bridge.instance.getItem(key)
}
export function setItem(key: string, value: string) {
    Bridge.instance.setItem(key, value)
}
export function removeItem(key: string) {
    Bridge.instance.removeItem(key)
}
export function clear() {
    Bridge.instance.clear()
}