local def = import "def.libsonnet";
{
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
            local deadline = def.Duration.Day*28,

            Backend: provider.Bolt,
            Memory: {
                    Access: access,
                    Refresh: refresh,
                    Deadline: deadline,
                    MaxSize: 1000*100,
            },
            Redis: {
                    URL: 'redis://redis:6379/0',
                    Access: access,
                    Refresh: refresh,
                    Deadline: deadline,
            },
            Bolt: {
                Filename: 'var/sessionid.db',

                Access: access,
                Refresh: refresh,
                Deadline: deadline,
                MaxSize: 1000*1000*10,
            },
        },
    },
}