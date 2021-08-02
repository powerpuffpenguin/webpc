#!/usr/bin/env bash

set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "pack release"
    echo
    echo "Usage:"
    echo "  $Command [flags]"
    echo
    echo "Flags:"
    echo "  -p, --pack          pack to compressed package (default \"gz\") [7z gz bz2 xz zip]"
    echo "  -h, --help          help for $Command"
}

ARGS=`getopt -o hp: --long help,pack: -n "$Command" -- "$@"`
eval set -- "${ARGS}"
pack="gz"
while true
do
    case "$1" in
        -h|--help)
            help
            exit 0
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

declare -i step
steps=${#Platforms[@]}
for i in ${!Platforms[@]};do
    if [[ $i != 0 ]];then
        echo
    fi
    
    platform=${Platforms[i]}
    step=i+1
    echo "step $step/$steps"
    "$BashDir/go.sh" -p "$pack" -P "$platform"
done