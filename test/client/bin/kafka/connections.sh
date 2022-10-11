#!/bin/bash
RESULT=`echo srvr | nc localhost 2181 | grep Connections`
if [ -z "$RESULT" ]
then
  # Expected "Connections: x" but result was empty
  echo "Error"
else
  echo $RESULT
fi