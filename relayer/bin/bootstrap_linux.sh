#!/bin/bash
set -e

echo "=====================================================================" >> relayer.log

nohup ./relayer start >> relayer.log 2>&1 &