#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
cd $SCRIPT_DIR
cp $SCRIPT_DIR/node_config.json ../node_modules/@vite/vuilder/bin/