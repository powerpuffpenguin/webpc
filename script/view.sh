#!/usr/bin/env bash

set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "static build helper"
    echo
    echo "Usage:"
    echo "  $Command [flags]"
    echo
    echo "Flags:"
    echo "  -i, --i18n          ng extract-i18n [en en-US zh-Hant zh-Hans]"
    echo "  -s, --static        build static"
    echo "  -h, --help          help for $Command"
}

ARGS=`getopt -o hi:s --long help,i18n:,static -n "$Command" -- "$@"`
eval set -- "${ARGS}"
i18n="none"
static=0
while true
do
    case "$1" in
        -h|--help)
            help
            exit 0
        ;;
        -i|--i18n)
            i18n="$2"
            shift 2
        ;;
        -s|--static)
            static=1
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


function build_en(){
    cd "$Dir/view"
    local args=(
        ng extract-i18n
    )

    local exec="${args[@]}"
    echo $exec
    eval "$exec"
}
function build_lang(){
    cd "$Dir/view"
    local args=(
        ng-xi18n update -s messages.xlf -l "$1"
    )

    local exec="${args[@]}"
    echo $exec
    eval "$exec"
}
function build_hans(){
    cd "$Dir/view"
    local args=(
        opencc -i src/locale/zh-Hant.xlf -o src/locale/zh-Hans.xlf -c t2s.json
    )

    local exec="${args[@]}"
    echo $exec
    eval "$exec"
}
case "$i18n" in
    none)
    ;;
    en)
        build_en
        exit $?
    ;;
    en-US)
        build_lang en-US
        exit $?
    ;;
    zh-Hant)
        build_lang zh-Hant
        exit $?
    ;;
    zh-Hans)
        build_hans
        exit $?
    ;;
    *)
        echo Error: unknown i18n "$i18n" for "$Command"
        echo "Run '$Command --help' for usage."
        exit 1
    ;;
esac

function build_static(){
    cd "$Dir"
    local args=(
        cp 
        view/dist/view/en/3rdpartylicenses.txt
        static/public/3rdpartylicenses.txt
    )
    local exec="${args[@]}"
    echo $exec
    eval "$exec"

    args=(
        cp 
        LICENSE
        static/public/LICENSE.txt
    )
    exec="${args[@]}"
    echo $exec
    eval "$exec"

    args=(
        statik
        -src=static/public
        -dest=assets/public
        -ns static -f
    )
    exec="${args[@]}"
    echo $exec
    eval "$exec"
    
    local items=(
        en-US
        zh-Hant
        zh-Hans
    )
    for i in ${!items[@]};do
        item=${items[i]}
        args=(
            statik -src="view/dist/view/$item" -dest "assets/$item"  -ns "$item" -f
        )
        exec="${args[@]}"
        echo $exec
        eval "$exec"
    done
}

if [[ $static == 1 ]];then
    build_static
    exit $?
fi

cd "$Dir/view"
args=(
 ng build 
 --configuration production
 --base-href /view/
 --localize
)
exec="${args[@]}"
echo $exec
eval "$exec"