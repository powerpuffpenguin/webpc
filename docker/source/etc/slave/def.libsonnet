// define some constants
{
    Size:{
        B: 1,
        KB: self.B*1024,
        MB: self.KB*1024,
        GB: self.MB*1024,
    },
    Duration:{
        Nanosecond: 1,
        Microsecond: 1000 * self.Nanosecond,
        Millisecond: 1000 * self.Microsecond,
        Second: 1000 * self.Millisecond,
        Minute: 60 * self.Second,
        Hour: 60 * self.Minute,
        Day: 24 * self.Hour,
    },
    Driver: {
        Sqlite3: 'sqlite3', // only support on linux (go env GOHOSTOS)
        Postgres: 'postgres',
        Mysql: 'mysql',
        Mssql: 'mssql',
    },
    Level:{
        Debug: 'debug',
        Info: 'info',
        Warn: 'warn',
        Error: 'error',
        Dpanic: 'dpanic',
        Panic: 'panic',
        Fatal: 'fatal',
    },
}