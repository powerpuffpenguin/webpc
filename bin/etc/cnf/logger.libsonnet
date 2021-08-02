local def = import "def.libsonnet";
local level = def.Level;
{
    // log file name
    Filename: 'var/logs/webpc.log',
    // Maximum size of a single log file
    // MB
    MaxSize: 100,
    // How many log files are saved
    MaxBackups: 3,
    // How many days of logs are kept
    MaxDays: 28,
    // Do you want to output the code location
    // Caller: true,
    // debug info warn error dpanic panic fatal
    FileLevel: level.Info,
    // debug info warn error dpanic panic fatal
    ConsoleLevel: level.Debug,
}