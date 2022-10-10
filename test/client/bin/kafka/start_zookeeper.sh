#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
KAFKA_VERSION="$(cat ./$SCRIPTDIR/kafka/version.txt)"
cd $SCRIPT_DIR
cd ..
if [ ! -d $KAFKA_VERSION ]
then
  echo "folder '$KAFKA_VERSION' does not exist"
  exit 1
fi
cd $KAFKA_VERSION
./bin/zookeeper-server-start.sh config/zookeeper.properties