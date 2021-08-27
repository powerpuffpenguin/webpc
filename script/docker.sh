#!/usr/bin/env bash
set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "docker build helper"
    echo
    echo "Usage:"
    echo "  $Command [flags]"
    echo
    echo "Flags:"
    echo "  -p, --push           push to hub"
    echo "  -h, --help          help for $Command"
}


ARGS=`getopt -o hp --long help,push -n "$Command" -- "$@"`
eval set -- "${ARGS}"
go=0
push=0
while true
do
    case "$1" in
        -h|--help)
            help
            exit 0
        ;;
        -g|--go)
            go=1
            shift
        ;;
        -p|--push)
            push=1
            shift
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

cd "$Dir/docker"
cp ../bin/webpc_linux_amd64.tar.gz source/webpc.tar.gz
args=(
    sudo docker build -t "\"$Docker:$Version\"" .
)
exec="${args[@]}"
echo $exec
eval "$exec"

if [[ "$push" == 1 ]];then
    args=(
        sudo docker push "\"$Docker:$Version\""
    )
    exec="${args[@]}"
    echo $exec
    eval "$exec"
fi