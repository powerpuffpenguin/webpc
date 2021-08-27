{
    // Register the master server as a virtual slave to allow web access
    Enable: true,
    //Shell : "shell-linux", // if empty, linux default value shell-linux.sh
    //Shell : "shell-windows.bat", // if empty, windows default value shell-windows.bat
    // vnc server address
    VNC: "127.0.0.1:5900",
    // mount path to web
    Mount: [
        {
            Name: "fs",
            Root: "/opt/webpc/fs",
            Write: true,
            Read: true,
            Shared: false,
        },
        {
            Name: "root",
            Root: "/",
            Write: false,
            Read: true,
            Shared: false,
        },
    ],
}