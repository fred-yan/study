#!/bin/bash

CUR_DIR=$(cd "$(dirname $0)"; pwd)
BIN_DIR=$CUR_DIR/../bin

for ((i=0;i<4;i++)); do
    let p1="11080+i"
    let p2="19000+i"
    mkdir -p $CUR_DIR/codis/$p1
    cd $CUR_DIR/codis/$p1
    cat > ${p1}.toml <<EOF
product_name = "codis-test"
product_auth = ""
proto_type = "tcp4"
admin_addr = "0.0.0.0:${p1}"
proxy_addr = "0.0.0.0:${p2}"
EOF
    nohup $BIN_DIR/codis-proxy -c ${p1}.toml &>${p1}.log &
    lastpid=$!
    pidlist="$pidlist $lastpid"
    echo "proxy-${p1}x${p2}.pid=$lastpid"
done
