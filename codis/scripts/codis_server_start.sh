#!/bin/bash

CUR_DIR=$(cd "$(dirname $0)"; pwd)
BIN_DIR=$CUR_DIR/../bin

for ((i=0;i<8;i++)); do
    let p="16379+i"
    mkdir -p $CUR_DIR/codis/$p
    cd $CUR_DIR/codis/$p
    nohup $BIN_DIR/codis-server --port ${p} &>redis-${p}.log &
    lastpid=$!
    pidlist="$pidlist $lastpid"
    echo "codis-server-${p}.pid=$lastpid"
done
