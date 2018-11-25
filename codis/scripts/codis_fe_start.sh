#!/bin/bash

CUR_DIR=$(cd "$(dirname $0)"; pwd)
BIN_DIR=$CUR_DIR/../bin

mkdir -p $CUR_DIR/codis/fe
cd $CUR_DIR/codis/fe

cat > codis.json <<EOF
[
    {
        "name": "codis-test",
        "dashboard": "127.0.0.1:18080"
    }
]
EOF

nohup $BIN_DIR/codis-fe -d codis.json --listen 0.0.0.0:8080 &> fe.log &
lastpid=$!
pidlist="$pidlist $lastpid"
echo "fe.pid=$lastpid"
