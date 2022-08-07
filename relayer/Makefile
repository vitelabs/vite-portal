MAIN_DIR=relayer
WORK_DIR=$(shell pwd)
MAIN=$(shell pwd)/cmd/$(MAIN_DIR)/main.go
BUILD_ROOT_DIR=$(shell pwd)/build
RELAYER_GIT_COMMIT=$(shell git rev-parse HEAD)
RELAYER_VERSION=$(shell cat version/buildversion)
RELAYER_VERSION_FILE=$(shell pwd)/version/buildversion.go

BUILD_DIR=$(BUILD_ROOT_DIR)/cmd/$(MAIN_DIR)
BUILD_BIN=$(BUILD_DIR)/relayer

build_version:
	@echo "package version" > $(RELAYER_VERSION_FILE)
	@echo "const RELAYER_COMMIT_VERSION = "\"$(RELAYER_GIT_COMMIT)\" >> $(RELAYER_VERSION_FILE)
	@echo "const RELAYER_BUILD_VERSION = "\"$(RELAYER_VERSION)\" >> $(RELAYER_VERSION_FILE)
	@echo "build version is "$(RELAYER_VERSION)", git commit is "$(RELAYER_GIT_COMMIT)"."