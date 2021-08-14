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
}