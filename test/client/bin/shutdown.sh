#!/bin/bash
pgrep relayer | xargs kill -9
pgrep relayer | xargs wait
# keep the debug logs util the next start
