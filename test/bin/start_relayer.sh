#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
cd $SCRIPT_DIR
cd ../../relayer
go build cmd/relayer/main.go 
mv main $SCRIPT_DIR/relayer
cd $SCRIPT_DIR
echo "--config-overrides: $1";
./relayer start --debug --config-overrides=$1 > logs/relayer_test.log