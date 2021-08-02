#!/usr/bin/env bash

set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "document build helper"
    echo
    echo "Usage:"
    echo "  $Command [flags]"
    echo
    echo "Flags:"
    echo "  -h, --help          help for $Command"
}

ARGS=`getopt -o h --long help -n "$Command" -- "$@"`
eval set -- "${ARGS}"
while true
do
    case "$1" in
        -h|--help)
            help
            exit 0
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

cd "$Dir"
args=(
    statik
    -src=static/document
    -dest=assets/document
    -ns document -f
)
exec="${args[@]}"
echo $exec
eval "$exec"