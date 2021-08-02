#!/usr/bin/env bash
set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi
filename="$Dir/version/version.go"
function write(){
    echo package version > "$filename"
    echo  >> "$filename"
    echo  "var (" >> "$filename"
    echo  "	Version = \`$Version\`" >> "$filename"
    echo  ")" >> "$filename"
}
if [[ -f "$filename" ]];then
    version=$(grep Version "$filename" | awk '{print $3}')
    if [[ "\`$Version\`" !=  "$version" ]];then
        write
    fi
else
    write
fi