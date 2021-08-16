export enum EventCode {
    Universal = 0,

    Heart = 1,

    Init = 2,
    Progress = 3,
    Success = 4,

    Yes = 5,
    No = 6,

    Exists = 7,
    YesAll = 8,
    Skip = 9,
    SkipAll = 10,
}
export function sendRequest(ws: WebSocket, evt: EventCode) {
    ws.send(JSON.stringify({
        event: evt,
    }))
}
export function fromString(val: string): EventCode {

    switch (val) {
        case 'Heart':
            return EventCode.Heart

        case 'Init':
            return EventCode.Init
        case 'Progress':
            return EventCode.Progress
        case 'Success':
            return EventCode.Success


        case 'Yes':
            return EventCode.Yes
        case 'No':
            return EventCode.No

        case 'Exists':
            return EventCode.Exists
        case 'YesAll':
            return EventCode.YesAll
        case 'Skip':
            return EventCode.Skip
        case 'SkipAll':
            return EventCode.SkipAll
    }
    return EventCode.Universal
}