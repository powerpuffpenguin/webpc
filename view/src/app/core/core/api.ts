import { MakeRESTful } from './restful';
const root = 'api'

export const ServerAPI = {
    v1: {
        sessions: MakeRESTful(root, 'v1', 'sessions'),
        system: MakeRESTful(root, 'v1', 'system'),
        groups: MakeRESTful(root, 'v1', 'groups'),
        users: MakeRESTful(root, 'v1', 'users'),
        slaves: MakeRESTful(root, 'v1', 'slaves'),
        logger: MakeRESTful(root, 'v1', 'logger'),
    },
    forward: {
        v1: {
            system: MakeRESTful(root, 'forward', 'v1', 'system'),
            fs: MakeRESTful(root, 'forward', 'v1', 'fs'),
            static: MakeRESTful(root, 'forward', 'v1', 'static'),
            logger: MakeRESTful(root, 'forward', 'v1', 'logger'),
        }
    },
    static: {
        licenses: MakeRESTful('static', '3rdpartylicenses.txt'),
        license: MakeRESTful('static', 'LICENSE.txt'),
    },
}
export enum Authorization {
    // Super administrator
    Root = 1,
    // access server
    Server = 2,
    // web shell
    Shell = 3,
    // filesystem read
    Read = 4,
    // filesystem write
    Write = 5,
    // vnc view
    VNC = 6,
}
export const Authorizations = [
    Authorization.Root,
    Authorization.Server,
    Authorization.Shell,
    Authorization.Read,
    Authorization.Write,
    Authorization.VNC,
]
export function AuthorizationName(authorization: Authorization): string {
    switch (authorization) {
        case Authorization.Root:
            return 'root'
        case Authorization.Server:
            return 'server'
        case Authorization.Shell:
            return 'shell'
        case Authorization.Read:
            return 'read'
        case Authorization.Write:
            return 'write'
        case Authorization.VNC:
            return 'vnc'
        default:
            return `${authorization}`
    }
}