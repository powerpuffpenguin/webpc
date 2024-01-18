import { ITerminalAddon, ITerminalInitOnlyOptions, ITerminalOptions, Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { WebLinksAddon } from "xterm-addon-web-links";
import { WebglAddon } from "xterm-addon-webgl";
export const DefaultFontFamily = "monospace"
export const DefaultFontSize = 15

export class MyTerminal {
    private term_?: Terminal
    get term(): Terminal | undefined {
        return this.term_
    }
    private fit_?: FitAddon
    private addons_?: Array<ITerminalAddon>
    constructor(options?: ITerminalOptions & ITerminalInitOnlyOptions) {
        const addons = new Array<ITerminalAddon>()
        const term = new Terminal(options)
        this.term_ = term
        this.addons_ = addons
        try {
            const webLinks = new WebLinksAddon()
            term.loadAddon(webLinks)
            addons.push(webLinks)
            try {
                const addon = new WebglAddon()
                term.loadAddon(addon)
                addon.onContextLoss(e => {
                    addon.dispose();
                });
                addons.push(addon)
            } catch (e) {
                console.warn('new WebglAddon fail', e)
            }
            const fit = new FitAddon()
            this.fit_ = fit
            term.loadAddon(fit)
            addons.push(fit)
        } catch (e) {
            this.close()
            throw e
        }
    }
    close() {
        const addons = this.addons_
        if (addons) {
            this.addons_ = undefined
            for (const addon of addons) {
                addon.dispose()
            }
        }
        const term = this.term_
        if (term) {
            this.term_ = undefined
            this.fit_ = undefined
            term.dispose()
        }
    }
    focus() {
        this.term_?.focus()
    }
    fit() {
        this.fit_?.fit()
    }
    clear() {
        this.term_?.clear()
    }
    write(data: string | Uint8Array, callback?: () => void) {
        this.term_?.write(data, callback)
    }
    writeln(data: string | Uint8Array, callback?: () => void) {
        this.term_?.writeln(data, callback)
    }

}