#!/bin/bash

CUR_DIR=$(cd "$(dirname $0)"; pwd)
BIN_DIR=$CUR_DIR/../bin

for ((i=0;i<3;i++)); do
    let p="26379+i"
    mkdir -p $CUR_DIR/codis/$p
    cd $CUR_DIR/codis/$p
    cat > ${p}-sentinel.conf <<EOF
protected-mode no
port $p
dir $CUR_DIR/codis/$p
daemonize yes
logfile "/$CUR_DIR/codis/$p/redis_sentinel.log"
EOF
    $BIN_DIR/redis-sentinel ${p}-sentinel.conf
done
