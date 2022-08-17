PROJECT_GIT_COMMIT=$(shell git rev-parse HEAD)
PROJECT_VERSION=$(shell cd .. && cat shared/pkg/version/buildversion)
PROJECT_VERSION_FILE=$(shell cd .. && pwd)/shared/pkg/version/buildversion.go

test:
	GO111MODULE=on go test ./...

build_version:
	@echo "package version" > $(PROJECT_VERSION_FILE)
	@echo "const PROJECT_COMMIT_VERSION = "\"$(PROJECT_GIT_COMMIT)\" >> $(PROJECT_VERSION_FILE)
	@echo "const PROJECT_BUILD_VERSION = "\"$(PROJECT_VERSION)\" >> $(PROJECT_VERSION_FILE)
	@echo "build version is "$(PROJECT_VERSION)", git commit is "$(PROJECT_GIT_COMMIT)"."