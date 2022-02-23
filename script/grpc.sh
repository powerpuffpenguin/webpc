#!/usr/bin/env bash

set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "grpc protoc helper"
    echo
    echo "Usage:"
    echo "  $Command [flags]"
    echo
    echo "Flags:"
    echo "  -l, --lang          generate grpc code for language (default \"go\") [go]"
    echo "  -h, --help          help for $Command"
}

ARGS=`getopt -o hl: --long help,lang: -n "$Command" -- "$@"`
eval set -- "${ARGS}"
lang="go"
while true
do
    case "$1" in
        -h|--help)
            help
            exit 0
        ;;
        -l|--lang)
            lang="$2"
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

function buildGo(){
    cd "$Dir"
    local output=protocol
    if [[ -d "$output" ]];then
        rm "$output" -rf
    fi
    if [[ -d "$UUID" ]];then
        rm "$UUID" -rf
    fi

    local document="static/document/api"
    if [[ -d "$document" ]];then
        rm "$document" -rf
        mkdir "$document"
    else
        mkdir "$document"
    fi

    local command=(
        protoc -I "pb" -I "third_party/googleapis" \
    )

    local protos=()
    for i in ${!Protos[@]};do
        protos[i]="$UUID/${Protos[i]}"
    done
    # generate grpc code
    local opts=(
        --go_out="'.'" --go_opt=paths=source_relative
        --go-grpc_out="'.'" --go-grpc_opt=paths=source_relative
        --grpc-gateway_out="'.'" --grpc-gateway_opt=paths=source_relative
        --openapiv2_out="'$document'" --openapiv2_opt logtostderr=true --openapiv2_opt use_go_templates=true --openapiv2_opt allow_merge=false
    )
    local exec="${command[@]} ${opts[@]} ${protos[@]}"
    echo $exec
    eval "$exec"

    # urls
    mv "$document/$UUID" "$document/docs"

    local url="$document/urls.js"
    echo "var URLS=[" > "$url"
    for i in ${!Protos[@]};do
        echo  "'${Protos[i]}'," >> "$url"
    done
    echo "];" >> "$url"

    echo mv $UUID protocol
    mv "$UUID" "protocol"
}

case "$lang" in
    go)
        buildGo
    ;;
    *)
        echo Error: unknown language "$1" for "$Command"
        echo "Run '$Command --help' for usage."
        exit 1
    ;;
esac