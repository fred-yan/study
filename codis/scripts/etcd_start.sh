#!/bin/bash

CUR_DIR=$(cd "$(dirname $0)"; pwd)
BIN_DIR=$CUR_DIR/../bin

mkdir -p $CUR_DIR/codis/etcd
cd $CUR_DIR/codis/etcd

nohup $BIN_DIR/etcd --name=codis-test &>etcd.log &
lastpid=$!
pidlist=$lastpid
echo "etcd.pid=$lastpid"
