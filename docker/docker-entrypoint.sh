#!/bin/bash

set -eu;

chown webpc.webpc /opt/webpc -R

if [[ "$@" == "execute-default" ]];then
    if [[ "$RUN_SLAVE" == 1 ]];then
        if [[ "$RUN_DEBUG" == 1 ]];then
            if [[ "$CONNECT_URL" == "" ]];then
                exec gosu webpc /opt/webpc/webpc slave --no-upgrade -d
            else
                exec gosu webpc /opt/webpc/webpc slave --no-upgrade -u "$CONNECT_URL" -d
            fi
        else
            if [[ "$CONNECT_URL" == "" ]];then
                exec gosu webpc /opt/webpc/webpc slave --no-upgrade
            else
                exec gosu webpc /opt/webpc/webpc slave --no-upgrade -u "$CONNECT_URL"
            fi
        fi
    else
        if [[ "$RUN_DEBUG" == 1 ]];then
            exec gosu webpc /opt/webpc/webpc master --no-upgrade -d
        else
            exec gosu webpc /opt/webpc/webpc master --no-upgrade
        fi    
    fi
else
    exec gosu webpc "$@"
fi