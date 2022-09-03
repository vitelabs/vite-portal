#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
cd $SCRIPT_DIR
cd ../../../orchestrator
go build cmd/orchestrator/main.go 
mv main $SCRIPT_DIR/orchestrator
$SCRIPT_DIR/orchestrator start --debug > $SCRIPT_DIR/logs/orchestrator_test.log