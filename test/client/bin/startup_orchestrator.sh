#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
cd $SCRIPT_DIR
cd ../../../orchestrator
rm -rf logs
mkdir logs
go build cmd/orchestrator/main.go 
mv main $SCRIPT_DIR/orchestrator
$SCRIPT_DIR/orchestrator start --debug > logs/test.log