#!/bin/bash

set -eu;

chown webpc.webpc /opt/webpc -R

if [[ "$@" == "execute-default" ]];then
    if [[ "$RUN_SLAVE" == 1 ]];then
        if [[ "$RUN_DEBUG" == 1 ]];then
            exec gosu webpc /opt/webpc/webpc slave -u "$CONNECT_URL" -d
        else
            exec gosu webpc /opt/webpc/webpc slave -u "$CONNECT_URL"
        fi
    else
        if [[ "$RUN_DEBUG" == 1 ]];then
            exec gosu webpc /opt/webpc/webpc master -d
        else
            exec gosu webpc /opt/webpc/webpc master
        fi    
    fi
else
    exec gosu webpc "$@"
fi