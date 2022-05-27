local def = import "def.libsonnet";
local size = def.Size;
local duration = def.Duration;
{
    // http addr
    Addr: ':9000',
    // if not empty use https
    CertFile: '',
    KeyFile: '',
    // enable swagger-ui on /document/
    Swagger: true,
    // grpc server option
    Option: {
        WriteBufferSize: 32*size.KB,
        ReadBufferSize: 32*size.KB,
        InitialWindowSize: 0*size.KB, // < 64k ignored
        InitialConnWindowSize: 0*size.KB, // < 64k ignored
        MaxRecvMsgSize: 0, // <1 6mb
        MaxSendMsgSize: 0, // <1 math.MaxInt32
        MaxConcurrentStreams: 0,
        ConnectionTimeout: 120 * duration.Second,
        Keepalive: {
            MaxConnectionIdle: 0,
            MaxConnectionAge: 0,
            MaxConnectionAgeGrace: 0,
            Time: 0,
            Timeout: 0,
        },
    },
}