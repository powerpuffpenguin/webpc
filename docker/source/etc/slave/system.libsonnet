{
    //Shell : "shell-linux", // if empty, linux default value shell-linux.sh
    //Shell : "shell-windows.bat", // if empty, windows default value shell-windows.bat
    // vnc server address
    VNC: "127.0.0.1:5900",
    // if true allow port forwarding
    PortForward: true,
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