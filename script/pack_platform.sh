#!/usr/bin/env bash

set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "pack platform"
    echo
    echo "Usage:"
    echo "  $0 [flags]"
    echo
    echo "Flags:"
    echo "  -d, --debug         build debug mode"
    echo "  -p, --pack          pack to compressed package (default \"gz\") [7z gz bz2 xz zip]"
    echo "  -P, --platform      build platform (default \"$(go env GOHOSTOS)/$(go env GOHOSTARCH)\")"
    echo "  -h, --help          help for $0"
}

ARGS=`getopt -o hdp:P: --long help,debug,pack:,platform: -n "$Command" -- "$@"`
eval set -- "${ARGS}"
debug=0
os="$(go env GOHOSTOS)"
arch="$(go env GOHOSTARCH)"
pack="gz"
while true
do
    case "$1" in
        -h|--help)
            help
            exit 0
        ;;
        -d|--debug)
            debug=1
            shift
        ;;
        -P|--platform)
            os=${2%\/*}
            arch=${2#*\/}
            if [[ "$os" == "" ]];then
                echo Error: unknown os "$2" for "$Command"
                echo "Run '$Command --list' for usage."
                exit 1
            fi
            if [[ "$arch" == "" ]];then
                echo Error: unknown arch "$2" for "$Command"
                echo "Run '$Command --list' for usage."
                exit 1
            fi
            shift 2
        ;;
        -p|--pack)
            case "$2" in
                7z|gz|bz2|xz|zip)
                    pack="$2"
                ;;
                *)
                    echo Error: unknown pack "$2" for "$Command --pack"
                    echo "Supported: 7z gz bz2 xz zip"
                    exit 1
                ;;
            esac
            shift 2
        ;;
        --)
            shift
            break
        ;;
        *)
            echo Error: unknown flag "$1" for "$Command"
            echo "Run '$Command --help' for usage."
            exit 1
        ;;
    esac
done

if [[ "$debug" == 1 ]];then
    target="${Target}d"
else
    target="${Target}"
fi
name="${target}_${GOOS}_$GOARCH"
if [[ "$os" == "windows" ]];then
    target="$target.exe"
fi
case "$pack" in
    7z)
        name="$name.7z"
        args=(7z a "$name")
    ;;
    zip)
        name="$name.zip"
        args=(zip -r "$name")
    ;;
    gz)
        name="$name.tar.gz"
        args=(tar -zcvf "$name")
    ;;
    bz2)
        name="$name.tar.bz2"
        args=(tar -jcvf "$name")
    ;;
    xz)
        name="$name.tar.xz"
        args=(tar -Jcvf "$name")
    ;;
esac
if [[ "$args" == "" ]];then
    exit 0
fi
cd "$Dir/bin"
if [[ -f "$name" ]];then
    rm "$name"
fi
if [[ "$os" == "windows" ]];then
    source=(
        "$target"
        etc
        winpty-agent.exe winpty.dll
        shell-windows.bat shell-bash.bat
        webpc-master.exe webpc-slave.exe webpc-master.xml webpc-slave.xml
        install-master.bat uninstall-master.bat install-slave.bat uninstall-slave.bat
    )
else
    source=(
        "$target"
        etc
        useradd.sh run-linux shell-linux
        webpc-master.service webpc-slave.service
    )
fi
exec="${args[@]} ${source[@]}"
echo $exec
eval "$exec >> /dev/null"

exec="sha256sum $name > $name.sha256"
echo $exec
eval "$exec"
