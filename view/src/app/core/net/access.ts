import { HttpClient } from "@angular/common/http"
import { Manager, Session } from "src/app/core/session/session"
export class Access {
    private session_: Session | undefined
    private refresh_ = false
    setSession(session?: Session) {
        this.session_ = session
        this.refresh_ = false
    }
    setRefresh() {
        this.refresh_ = true
    }
    close() {
        this.session_ = undefined
        this.refresh_ = false
    }
    refresh(httpClient: HttpClient): undefined | Promise<Session | undefined> {
        if (this.refresh_ && this.session_) {
            return Manager.instance.refresh(httpClient, this.session_)
        } else {
            return
        }
    }
}