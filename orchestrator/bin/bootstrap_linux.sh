#!/bin/bash
set -e

echo "=====================================================================" >> orchestrator.log

nohup ./orchestrator start >> orchestrator.log 2>&1 &