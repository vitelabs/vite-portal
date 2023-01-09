#!/bin/bash
pgrep orchestrator | xargs kill -9
pgrep orchestrator | xargs wait