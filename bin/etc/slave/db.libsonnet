local def = import "def.libsonnet";
local driver = def.Driver;
{
    Driver: driver.Sqlite3,
    Source: [
        'var/webpc_slave.db',
    ],
    // max connections if < 1 not not limited
    MaxOpen: 50,
    // idle connections if < 1  not exists idle
    MaxIdle: 5,
    // ShowSQL: true,
    Cache: {
        // default cache rows
        Record: 1000,
        // disable cache table names
        Direct: [],
        // special cache
        Special: [
        ],
    },
}