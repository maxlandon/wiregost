# -------------------------------------------------------------------------------------------------
# 					Wiregost Makefile 
# -------------------------------------------------------------------------------------------------


GO ?= go
ENV = CGO_ENABLED=0
TAGS = -tags netgo,go_sqlite

# https://stackoverflow.com/questions/5618615/check-if-a-program-exists-from-a-makefile
EXECUTABLES = protoc protoc-gen-go sed git zip go
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

SED_INPLACE := sed -i

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	SED_INPLACE := sed -i ''
endif

#
# Version Information ------------------------------------------------------------------------------
 
VERSION ?= $(shell git describe --abbrev=0)
COMPILED_AT = $(shell date +%s)
RELEASES_URL = https://api.github.com/repos/maxlandon/wiregost/releases
GIT_DIRTY = $(shell git diff --quiet|| echo 'Dirty')
GIT_COMMIT = $(shell git rev-parse HEAD)

# Client 
CLIENT_PKG = github.com/maxlandon/wiregost/internal/client/version
CLIENT_LDFLAGS = -ldflags "-s -w \
	-X $(CLIENT_PKG).Version=$(VERSION) \
	-X $(CLIENT_PKG).CompiledAt=$(COMPILED_AT) \
	-X $(CLIENT_PKG).GithubReleasesURL=$(RELEASES_URL) \
	-X $(CLIENT_PKG).GitCommit=$(GIT_COMMIT) \
	-X $(CLIENT_PKG).GitDirty=$(GIT_DIRTY)"

# Server 
SERVER_PKG = github.com/maxlandon/wiregost/internal/server/version
SERVER_LDFLAGS = -ldflags "-s -w \
	-X $(SERVER_PKG).Version=$(VERSION) \
	-X $(SERVER_PKG).CompiledAt=$(COMPILED_AT) \
	-X $(SERVER_PKG).GithubReleasesURL=$(RELEASES_URL) \
	-X $(SERVER_PKG).GitCommit=$(GIT_COMMIT) \
	-X $(SERVER_PKG).GitDirty=$(GIT_DIRTY)"


# TARGETS ------------------------------------------------------------------------------------------
#
.PHONY: macos
macos: clean proto
	GOOS=darwin $(ENV) $(GO) build $(TAGS) $(SERVER_LDFLAGS) -o ./cmd/server/wiregost-server ./cmd/server
	GOOS=darwin $(ENV) $(GO) build $(TAGS) $(CLIENT_LDFLAGS) -o ./cmd/client/wiregost-client ./cmd/client

.PHONY: linux
linux: clean proto
	GOOS=linux $(ENV) $(GO) build $(TAGS) $(SERVER_LDFLAGS) -o ./cmd/server/wiregost-server ./cmd/server
	GOOS=linux $(ENV) $(GO) build $(TAGS) $(CLIENT_LDFLAGS) -o ./cmd/client/wiregost-client ./cmd/client

.PHONY: windows
windows: clean proto
	GOOS=windows $(ENV) $(GO) build $(TAGS) $(SERVER_LDFLAGS) -o ./cmd/server/wiregost-server.exe ./cmd/server
	GOOS=windows $(ENV) $(GO) build $(TAGS) $(CLIENT_LDFLAGS) -o ./cmd/client/wiregost-client.exe ./cmd/client


# Accessory Makes ----------------------------------------------------------------------------------

.PHONY: deps
deps:
	# Install protoc plugins
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# proto is a target that uses prototool.
# By depending on $(PROTOTOOL), prototool will be installed on the Makefile's path.
# Since the path above has the temporary GOBIN at the front, this will use the
# locally installed prototool.
.PHONY: proto
proto: $(PROTOTOOL)
	cd internal/proto; buf generate

.PHONY: clean
clean:
	rm -f ./cmd/client/wiregost-client ./cmd/client/wiregost-client_* ./cmd/server/wiregost-server ./cmd/server/wiregost-server_*
