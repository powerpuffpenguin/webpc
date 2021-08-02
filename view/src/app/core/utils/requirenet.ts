declare const __requireNet: any
export function RequireNet(...names: Array<string>): Promise<any> {
    let result: Promise<any>
    switch (names.length) {
        case 0:
            result = Promise.resolve()
            break;
        case 1:
            result = new Promise<any>(function (resolve, reject) {
                __requireNet(names, resolve, reject)
            })
            break
        default:
            result = new Promise<any>(function (resolve, reject) {
                __requireNet(names,
                    function () {
                        resolve(arguments)
                    },
                    reject,
                )
            })
            break;
    }
    return result
}