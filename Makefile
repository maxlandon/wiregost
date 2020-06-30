# -------------------------------------------------------------------------------------------------
# 					Makefile for Wiregost 
# -------------------------------------------------------------------------------------------------


GO ?= go
ENV = CGO_ENABLED=0
TAGS = -tags netgo

# https://stackoverflow.com/questions/5618615/check-if-a-program-exists-from-a-makefile
EXECUTABLES = protoc protoc-gen-go packr sed git zip go
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

SED_INPLACE := sed -i

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	SED_INPLACE := sed -i ''
endif

#
# Version Information ------------------------------------------------------------------------------
 
VERSION = 1.0.0
COMPILED_AT = $(shell date +%s)
RELEASES_URL = https://api.github.com/repos/maxlandon/wiregost/releases
GIT_DIRTY = $(shell git diff --quiet|| echo 'Dirty')
GIT_COMMIT = $(shell git rev-parse HEAD)

# Client 
CLIENT_PKG = github.com/maxlandon/wiregost/client/version
CLIENT_LDFLAGS = -ldflags "-s -w \
	-X $(CLIENT_PKG).Version=$(VERSION) \
	-X $(CLIENT_PKG).CompiledAt=$(COMPILED_AT) \
	-X $(CLIENT_PKG).GithubReleasesURL=$(RELEASES_URL) \
	-X $(CLIENT_PKG).GitCommit=$(GIT_COMMIT) \
	-X $(CLIENT_PKG).GitDirty=$(GIT_DIRTY)"

# Server 
SERVER_PKG = github.com/maxlandon/wiregost/server/version
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
	GOOS=darwin $(ENV) $(GO) build $(TAGS) $(SERVER_LDFLAGS) -o wiregost-server ./server
	GOOS=darwin $(ENV) $(GO) build $(TAGS) $(CLIENT_LDFLAGS) -o wiregost-console ./client

.PHONY: linux
linux: clean proto 
	GOOS=linux $(ENV) $(GO) build $(TAGS) $(SERVER_LDFLAGS) -o wiregost-server ./server
	GOOS=linux $(ENV) $(GO) build $(TAGS) $(CLIENT_LDFLAGS) -o wiregost-console ./client

.PHONY: windows
windows: clean proto 
	GOOS=windows $(ENV) $(GO) build $(TAGS) $(SERVER_LDFLAGS) -o wiregost-server.exe ./server
	GOOS=windows $(ENV) $(GO) build $(TAGS) $(CLIENT_LDFLAGS) -o wiregost-console.exe ./client


#
# Static builds were we bundle everything together 
#
# MacOS 
.PHONY: static-server-macos
static-server-macos: clean proto packr
	packr
	$(SED_INPLACE) '/$*.windows\/go\.zip/d' ./server/assets/a_assets-packr.go
	$(SED_INPLACE) '/$*.linux\/go\.zip/d' ./server/assets/a_assets-packr.go
	GOOS=darwin $(ENV) $(GO) build $(TAGS) $(SERVER_LDFLAGS) -o wiregost-server ./server

.PHONY: console-macos
console-macos: clean proto
	GOOS=darwin $(ENV) $(GO) build $(TAGS) $(CLIENT_LDFLAGS) -o wiregost-console ./client

# Windows
.PHONY: static-server-windows
static-server-windows: clean proto packr
	packr
	$(SED_INPLACE) '/$*.darwin\/go\.zip/d' ./server/assets/a_assets-packr.go
	$(SED_INPLACE) '/$*.linux\/go\.zip/d' ./server/assets/a_assets-packr.go
	GOOS=windows $(ENV) $(GO) build $(TAGS) $(SERVER_LDFLAGS) -o wiregost-server.exe ./server

.PHONY: console-windows
console-windows: clean proto
	GOOS=windows $(ENV) $(GO) build $(TAGS) $(CLIENT_LDFLAGS) -o wiregost-console.exe ./client

# Linux 
.PHONY: static-server-linux
static-server-linux: clean proto packr
	$(SED_INPLACE) '/$*.darwin\/go\.zip/d' ./server/assets/a_assets-packr.go
	$(SED_INPLACE) '/$*.windows\/go\.zip/d' ./server/assets/a_assets-packr.go
	GOOS=linux $(ENV) $(GO) build $(TAGS) $(SERVER_LDFLAGS) -o wiregost-server ./server

.PHONY: console-linux
console-linux: clean proto 
	GOOS=linux $(ENV) $(GO) build $(TAGS) $(CLIENT_LDFLAGS) -o wiregost-console ./client

# All-In-One Release compilation & compression
.PHONY: release
release:
	mkdir -p release-${VERSION}/linux
	mkdir -p release-${VERSION}/macos
	mkdir -p release-${VERSION}/windows

	$(MAKE) console-linux
	zip release-${VERSION}/linux/wiregost-console_linux.zip ./client/wiregost-console.go
	$(MAKE) static-linux
	zip release-${VERSION}/linux/wiregost-server_linux.zip ./server/wiregost-server.go

	$(MAKE) macos
	zip release-${VERSION}/macos/wiregost-console_macos.zip ./client/wiregost-console.go
	$(MAKE) static-macos
	zip release-${VERSION}/macos/wiregost-server_macos.zip ./server/wiregost-server.go

	$(MAKE) windows
	zip release-${VERSION}/windows/wiregost-console_windows.zip ./sliver-client.exe
	$(MAKE) static-windows
	zip release-${VERSION}/windows/wiregost-server_windows.zip ./server/wiregost-server.exe


# Accessory Makes ------------------------------------------------------------------------------------------

.PHONY: packr
packr:
	cd ./server/
	packr
	cd ..

.PHONY: clean-all
clean-all: clean
	rm -f ./assets/darwin/go.zip
	rm -f ./assets/windows/go.zip
	rm -f ./assets/linux/go.zip
	rm -f ./assets/*.zip

.PHONY: clean
clean: 
	packr clean
	rm -f ./protobuf/client/*.proto.go
	rm -f ./protobuf/ghost/*.proto.go
	rm -f wiregost-console wiregost-server *.exe

# Generate Struct tags with protoc-gen-tags-go
.PHONY: tags 
tags: 
	./proto/generate-go-tags.sh


# Prototool Compilation ------------------------------------------------------------------------------------------

SHELL := /bin/bash -o pipefail

UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)

TMP_BASE := ./proto/bin
TMP := $(TMP_BASE)/$(UNAME_OS)/$(UNAME_ARCH)
TMP_BIN = $(TMP)
TMP_VERSIONS := $(TMP)/versions

export GO111MODULE := on
export GOBIN := $(abspath $(TMP_BIN))
export PATH := $(GOBIN):$(PATH)

# This is the only variable that ever should change.
# This can be a branch, tag, or commit.
# When changed, the given version of Prototool will be installed to
# .tmp/$(uname -s)/(uname -m)/bin/prototool
PROTOTOOL_VERSION := v1.9.0
PROTOC_GEN_VERSION := v1.4.0

PROTOTOOL := $(TMP_VERSIONS)/prototool/$(PROTOTOOL_VERSION)
$(PROTOTOOL):
	$(eval PROTOTOOL_TMP := $(shell mktemp -d))
	cd $(PROTOTOOL_TMP); go get github.com/uber/prototool/cmd/prototool@$(PROTOTOOL_VERSION); 
	cd $(PROTOTOOL_TMP); go get google.golang.org/protobuf/cmd/protoc-gen-go
	# cd $(PROTOTOOL_TMP); go install github.com/grpc/grpc-go/cmd/protoc-gen-go-grpc       // Should not be needed because go.mod already has it
	@rm -rf $(PROTOTOOL_TMP)
	@rm -rf $(dir $(PROTOTOOL))
	@mkdir -p $(dir $(PROTOTOOL))
	@touch $(PROTOTOOL)

# proto is a target that uses prototool.
# By depending on $(PROTOTOOL), prototool will be installed on the Makefile's path.
# Since the path above has the temporary GOBIN at the front, this will use the
# locally installed prototool.
.PHONY: proto
proto: $(PROTOTOOL)
	cd proto/; prototool generate
