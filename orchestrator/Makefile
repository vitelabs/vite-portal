NAME=orchestrator
CONFIG=orchestrator_config.json
WORK_DIR=$(shell pwd)
MAIN=$(shell pwd)/cmd/$(NAME)/main.go
BUILD_ROOT_DIR=$(shell pwd)/build
PROJECT_GIT_COMMIT=$(shell git rev-parse HEAD)
PROJECT_VERSION=$(shell cat version/buildversion)
PROJECT_VERSION_FILE=$(shell pwd)/version/buildversion.go

BUILD_DIR=$(BUILD_ROOT_DIR)/cmd/$(NAME)
BUILD_BIN=$(BUILD_DIR)/$(NAME)

build:
	GO111MODULE=on go build -o $(BUILD_BIN) $(MAIN)
	@echo "Build $(NAME) done."
	@echo "Run $(BUILD_DIR)/$(NAME) to start $(NAME)."

test:
	GO111MODULE=on go test ./...

create_config:
	GO111MODULE=on go run ./config/main.go

build_version: create_config
	@echo "package version" > $(PROJECT_VERSION_FILE)
	@echo "const PROJECT_COMMIT_VERSION = "\"$(PROJECT_GIT_COMMIT)\" >> $(PROJECT_VERSION_FILE)
	@echo "const PROJECT_BUILD_VERSION = "\"$(PROJECT_VERSION)\" >> $(PROJECT_VERSION_FILE)
	@echo "build version is "$(PROJECT_VERSION)", git commit is "$(PROJECT_GIT_COMMIT)"."

build_linux_amd64:
	$(eval BUILD_DIR_LOCAL := $(BUILD_DIR)/$(NAME)-$(PROJECT_VERSION)-linux-amd64)
	env GOOS=linux CGO_ENABLED=0 GO111MODULE=on GOARCH=amd64 go build -o $(BUILD_DIR_LOCAL)/$(NAME) $(MAIN)

	@cp $(shell pwd)/config/$(CONFIG) $(BUILD_DIR_LOCAL)/$(CONFIG)
	@cp $(shell pwd)/bin/bootstrap_linux.sh $(BUILD_DIR_LOCAL)/bootstrap.sh

	@ls -d $(BUILD_DIR_LOCAL)/$(NAME)
	@echo "Build $(NAME)-linux-amd64 version done."

build_linux_arm64:
	$(eval BUILD_DIR_LOCAL := $(BUILD_DIR)/$(NAME)-$(PROJECT_VERSION)-linux-arm64)
	env GOOS=linux CGO_ENABLED=0 GO111MODULE=on GOARCH=arm64 go build -o $(BUILD_DIR_LOCAL)/$(NAME) $(MAIN)

	@cp $(shell pwd)/config/$(CONFIG) $(BUILD_DIR_LOCAL)/$(CONFIG)
	@cp $(shell pwd)/bin/bootstrap_linux.sh $(BUILD_DIR_LOCAL)/bootstrap.sh

	@ls -d $(BUILD_DIR_LOCAL)/$(NAME)
	@echo "Build $(NAME)-linux-arm64 version done."

app-linux-amd64: build_version build_linux_amd64

app-linux-arm64: build_version build_linux_arm64

app: build_version build

all: app-linux-amd64 app-linux-arm64

clean:
	rm -r $(BUILD_ROOT_DIR)/