#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
cd $SCRIPT_DIR
gnome-terminal -x ./kafka/start_zookeeper.sh
sleep 2
gnome-terminal -x ./kafka/start_broker.sh