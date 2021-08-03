import { MakeRESTful } from './restful';
const root = 'api'

export const ServerAPI = {
    v1: {
        sessions: MakeRESTful(root, 'v1', 'sessions'),
        system: MakeRESTful(root, 'v1', 'system'),
        users: MakeRESTful(root, 'v1', 'users'),
        slaves: MakeRESTful(root, 'v1', 'slaves'),
        logger: MakeRESTful(root, 'v1', 'logger'),
    },
    static: {
        licenses: MakeRESTful('static', '3rdpartylicenses.txt'),
        license: MakeRESTful('static', 'LICENSE.txt'),
    },
}
export enum Authorization {
    // super administrator
    Root = 1,
}
export const Authorizations = [Authorization.Root]
export function AuthorizationName(authorization: Authorization): string {
    switch (authorization) {
        case Authorization.Root:
            return 'root'
        default:
            return `${authorization}`
    }
}