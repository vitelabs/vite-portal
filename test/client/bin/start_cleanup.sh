#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
cd $SCRIPT_DIR
rm -rf logs
mkdir logs
rm -rf kafka_data
mkdir kafka_data