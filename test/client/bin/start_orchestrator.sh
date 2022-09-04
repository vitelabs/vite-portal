#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
cd $SCRIPT_DIR
cd ../../../orchestrator
go build cmd/orchestrator/main.go 
mv main $SCRIPT_DIR/orchestrator
cd $SCRIPT_DIR
echo "--config-overrides: $1";
./orchestrator start --debug > logs/orchestrator_test.log