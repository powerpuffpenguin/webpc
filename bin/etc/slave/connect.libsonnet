local def = import "def.libsonnet";
local size = def.Size;
local duration = def.Duration;
{
    // http addr
    URL: 'ws://127.0.0.1:9000/api/v1/dialer/64048031f73a11eba3890242ac120064',
    // if true allow insecure server connections when using SSL
    // Insecure: true, 
    // grpc server option
    Option: {
        WriteBufferSize: 32*size.KB,
        ReadBufferSize: 32*size.KB,
        InitialWindowSize: 0*size.KB, // < 64k ignored
        InitialConnWindowSize: 0*size.KB, // < 64k ignored
        MaxRecvMsgSize: 0, // <1 4mb
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