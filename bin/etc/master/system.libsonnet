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
            // web display name
            Name: "movie",
            // local filesystem path
            Root: "/home/dev/movie",
            // Set the directory to be readable. Users with read/write permissions can read files
            Read: true,
            // Set the directory to be writable. Users with write permission can write files
            // If Write is true, Read will be forcibly set to true
            Write: true,
            // Set as a shared directory to allow anyone to read the file
            // If Shared is true, Read will be forcibly set to true
            Shared: true,
        },
        {
            Name: "home",
            Root: "/home/dev",
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
        {
            Name: "media",
            Root: "/media/dev/",
            Write: false,
            Read: true,
            Shared: false,
        },
    ],
}