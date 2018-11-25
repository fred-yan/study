#!/bin/bash

CUR_DIR=$(cd "$(dirname $0)"; pwd)
BIN_DIR=$CUR_DIR/../bin

mkdir -p $CUR_DIR/codis/dashboard
cd $CUR_DIR/codis/dashboard

cat > dashboard.toml <<EOF
coordinator_name = "etcd"
coordinator_addr = "127.0.0.1:2379"
product_name = "codis-test"
product_auth = ""
admin_addr = "0.0.0.0:18080"
EOF

nohup $BIN_DIR/codis-dashboard -c dashboard.toml &> dashboard.log &
lastpid=$!
pidlist="$pidlist $lastpid"
echo "dashboard.pid=$lastpid"

