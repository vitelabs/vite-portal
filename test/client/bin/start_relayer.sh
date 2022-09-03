#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
cd $SCRIPT_DIR
cd ../../../relayer
go build cmd/relayer/main.go 
mv main $SCRIPT_DIR/relayer
$SCRIPT_DIR/relayer start --debug > $SCRIPT_DIR/logs/relayer_test.log