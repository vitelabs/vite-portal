#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
KAFKA_VERSION="$(cat ./$SCRIPTDIR/kafka/version.txt)"
cd $SCRIPT_DIR
cd ..
rm -rf kafka_data/kafka-logs
mkdir -p kafka_data/kafka-logs
if [ ! -d $KAFKA_VERSION ]
then
  echo "folder '$KAFKA_VERSION' does not exist"
  exit 1
fi
cd $KAFKA_VERSION
./bin/kafka-server-start.sh config/server.properties