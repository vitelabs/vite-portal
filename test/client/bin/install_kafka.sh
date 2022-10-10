#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
cd $SCRIPT_DIR
KAFKA_VERSION=kafka_2.13-3.3.1
TAR_NAME=$KAFKA_VERSION.tgz
DOWNLOAD_URL=https://dlcdn.apache.org/kafka/3.3.1/$TAR_NAME
DATA_DIR=$SCRIPT_DIR/kafka_data
echo "Apache Kafka installing: $TAR_NAME"
download_kafka () {
  if command -v curl >/dev/null 2>&1
  then
      echo "using curl to download Apache Kafka"
      curl -O $DOWNLOAD_URL
  else
      if command -v wget >/dev/null 2>&1
      then
          echo "using wget to download Apache Kafka"
          wget -c $DOWNLOAD_URL
      else
          echo "can't download Apache Kafka"
          exit 1
      fi
  fi
  echo "Apache Kafka download finished"
}
if [ ! -f $TAR_NAME ]
then
  download_kafka
fi
if [ -d $KAFKA_VERSION ]
then
  rm -rf $KAFKA_VERSION
fi
echo "Apache Kafka extracting..."
tar -xzf $TAR_NAME
echo "Apache Kafka configuring..."
rm -rf $DATA_DIR
sed -i "s|/tmp/kafka-logs|$DATA_DIR/kafka-logs|g" $SCRIPT_DIR/$KAFKA_VERSION/config/server.properties
sed -i "s|/tmp/zookeeper|$DATA_DIR/zookeeper|g" $SCRIPT_DIR/$KAFKA_VERSION/config/zookeeper.properties
mkdir $DATA_DIR
echo "Apache Kafka installed"
