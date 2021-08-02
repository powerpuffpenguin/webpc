#!/usr/bin/env bash
set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "go build helper"
    echo
    echo "Usage:"
    echo "  $Command [flags]"
    echo
    echo "Flags:"
    echo "  -c, --clear         clear output"
    echo "  -d, --debug         build debug mode"
    echo "  -l, --list          list all supported platforms"
    echo "  -p, --pack          pack to compressed package [7z gz bz2 xz zip]"
    echo "  -P, --platform      build platform (default \"$(go env GOHOSTOS)/$(go env GOHOSTARCH)\")"
    echo "  -u, --upx           use upx to compress executable programs"
    echo "  -h, --help          help for $Command"
}


ARGS=`getopt -o hldp:P:cu --long help,list,os:,arch:,debug,pack:,platform:,clear,upx -n "$Command" -- "$@"`
eval set -- "${ARGS}"
list=0
debug=0
clear=0
os="$(go env GOHOSTOS)"
arch="$(go env GOHOSTARCH)"
pack=""
upx=0
while true
do
    case "$1" in
        -h|--help)
            help
            exit 0
        ;;
        -l|--list)
            list=1
            shift
        ;;
        -u|--upx)
            upx=1
            shift
        ;;
        -d|--debug)
            debug=1
            shift
        ;;
        -c|--clear)
            clear=1
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
            pack="$2"
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

if [[ "$list" == 1 ]];then
    go tool dist list
    exit $?
fi
if [[ "$clear" == 1 ]];then
    "$BashDir/clear.sh"
    exit $?
fi

export GOOS="$os"
export GOARCH="$arch"
target="$Target"
if [[ "$debug" == 1 ]];then
    target="$target"d
fi
if [[ "$os" == "windows" ]];then
    target="$target.exe"
fi
if [[ "$debug" == 1 ]];then
    args=(
        go build 
        -o "bin/$target"
    )
else
    args=(
        go build 
        -ldflags '"-s -w"'
        -o "bin/$target"
    )
fi

cd "$Dir"
# version
"$BashDir/version.sh"

# build
echo "build for \"$GOOS/$GOARCH\""
exec="${args[@]}"
echo $exec
eval "$exec"

# upx 
if [[ $upx == 1 ]];then
    upx "bin/$target"
fi

# pack
if [[ "$pack" == "" ]];then
    exit 0
fi
"$BashDir/pack_platform.sh" -p "$pack" -P "$os/$arch"