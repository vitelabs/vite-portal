#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
cd $SCRIPT_DIR
# Update node_config.json
cp $SCRIPT_DIR/node_config.json ../node_modules/@vite/vuilder/bin/
# Update relayer_config.json
cp $SCRIPT_DIR/relayer_config.json ../../../relayer/