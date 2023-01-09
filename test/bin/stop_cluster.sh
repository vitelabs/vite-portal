#!/bin/bash
SCRIPT_DIR=$(cd $(dirname ${BASH_SOURCE[0]}); pwd)
cd $SCRIPT_DIR
cd ..
docker-compose kill
# make sure volumes are not in use and can be removed
docker-compose down
./docker_remove_volumes.sh