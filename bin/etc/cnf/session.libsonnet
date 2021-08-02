local def = import "def.libsonnet";
{
    Backend: def.Session.Memory,
    Coder: def.Coder.JSON,
    Client: {
        Protocol: def.Protocol.H2C,
        Addr: 'sessionid_server:80',
        Token: '',
    },
    Memory: {
        Manager: {
            // Token signature algorithm
            Method: def.Method.HMD5,
            // Signing key
            Key: 'cerberus is an idea',
        },
        Provider: {
            local provider = def.Provider,
            local access = def.Duration.Hour*1,
            local refresh = def.Duration.Hour*12*3,
            local clear = def.Duration.Minute*30,

            Backend: provider.Bolt,
            Memory: {
                    Access: access,
                    Refresh: refresh,
                    MaxSize: 1000,
                    Batch: 128,
                    Clear: clear,
            },
            Redis: {
                    URL: 'redis://redis:6379/0',
                    Access: access,
                    Refresh: refresh,
                    Batch: 128,
                    KeyPrefix: 'sessionid.provider.redis.',
                    MetadataKey: '__private_provider_redis',
            },
            Bolt: {
                Filename: 'var/sessionid.db',

                Access: access,
                Refresh: refresh,
                MaxSize: 10000,
                Batch: 128,
                Clear: clear,
            },
        },
    },
}