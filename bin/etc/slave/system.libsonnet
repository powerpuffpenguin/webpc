{
    //Shell : "shell-linux", // if empty, linux default value shell-linux.sh
    //Shell : "shell-windows.bat", // if empty, windows default value shell-windows.bat
    // vnc server address
    VNC: "127.0.0.1:5900",
    // mount path to web
    Mount: [
        {
            // web display name
            Name: "s_movie",
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
            Name: "s_home",
            Root: "/home/dev",
            Write: true,
            Read: true,
            Shared: false,
        },
        {
            Name: "s_root",
            Root: "/",
            Write: false,
            Read: true,
            Shared: false,
        },
        {
            Name: "s_media",
            Root: "/media/dev/",
            Write: false,
            Read: true,
            Shared: false,
        },
    ],
}