#
# Makefile for Wiregost 
#

GO ?= go
ENV = CGO_ENABLED=0
TAGS = -tags netgo
LDFLAGS = -ldflags '-s -w'

# https://stackoverflow.com/questions/5618615/check-if-a-program-exists-from-a-makefile
EXECUTABLES = protoc protoc-gen-go packr sed git zip go
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH")))

GIT_DIRTY = $(shell git diff --quiet|| echo 'Dirty')
GIT_VERSION = $(shell git rev-parse HEAD)

SED_INPLACE := sed -i

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	SED_INPLACE := sed -i ''
endif


.PHONY: macos
macos: clean version pb
	GOOS=darwin $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-server ./server
	GOOS=darwin $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-console ./client
	GOOS=darwin $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-data-service ./data-service

.PHONY: linux
linux: clean version pb
	GOOS=linux $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-server ./server
	GOOS=linux $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-console ./client
	GOOS=linux $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-data-service ./data-service

.PHONY: windows
windows: clean version pb
	GOOS=windows $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-server.exe ./server
	GOOS=windows $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-console.exe ./client
	GOOS=windows $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-data-service.exe ./data-service


#
# Static builds were we bundle everything together ---------------------------------------------------
#
# MacOS 
.PHONY: static-server-macos
static-server-macos: clean version pb packr
	packr
	$(SED_INPLACE) '/$*.windows\/go\.zip/d' ./server/assets/a_assets-packr.go
	$(SED_INPLACE) '/$*.linux\/go\.zip/d' ./server/assets/a_assets-packr.go
	GOOS=darwin $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-server ./server

.PHONY: console-macos
console-macos: clean version pb
	GOOS=darwin $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-console ./client

.PHONY: data-service-macos
data-service-macos: clean version pb
	GOOS=darwin $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-data-service ./data-service

# Windows
.PHONY: static-server-windows
static-server-windows: clean version pb packr
	packr
	$(SED_INPLACE) '/$*.darwin\/go\.zip/d' ./server/assets/a_assets-packr.go
	$(SED_INPLACE) '/$*.linux\/go\.zip/d' ./server/assets/a_assets-packr.go
	GOOS=windows $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-server.exe ./server

.PHONY: console-windows
console-windows: clean version pb
	GOOS=windows $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-console.exe ./client

.PHONY: data-service-windows
data-service-windows: clean version pb
	GOOS=windows $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-data-service.exe ./data-service

# Linux 
.PHONY: static-server-linux
static-server-linux: clean version pb packr
	$(SED_INPLACE) '/$*.darwin\/go\.zip/d' ./server/assets/a_assets-packr.go
	$(SED_INPLACE) '/$*.windows\/go\.zip/d' ./server/assets/a_assets-packr.go
	GOOS=linux $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-server ./server

.PHONY: console-linux
console-linux: clean version pb
	GOOS=linux $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-console ./client

.PHONY: data-service-linux
data-service-linux: clean version pb
	GOOS=linux $(ENV) $(GO) build $(TAGS) $(LDFLAGS) -o wiregost-data-service ./data-service


# Accessory Makes ------------------------------------------------------------------------------------------
.PHONY: version
version:
	printf "package version\n\nconst GitVersion = \"%s\"\n" $(GIT_VERSION) > ./client/version/version.go
	printf "const GitDirty = \"%s\"\n" $(GIT_DIRTY) >> ./client/version/version.go

.PHONY: packr
packr:
	cd ./server/
	packr
	cd ..

.PHONY: clean-version
clean-version:
	printf "package version\n\nconst GitVersion = \"\"\n" > ./client/version/version.go

.PHONY: clean-all
clean-all: clean clean-version
	rm -f ./assets/darwin/go.zip
	rm -f ./assets/windows/go.zip
	rm -f ./assets/linux/go.zip
	rm -f ./assets/*.zip

.PHONY: clean
clean: clean-version
	packr clean
	rm -f ./protobuf/client/*.pb.go
	rm -f ./protobuf/ghost/*.pb.go
	rm -f wiregost-console wiregost-server *.exe

# Generate Struct tags with protoc-gen-tags-go
.PHONY: tags 
tags: 
	./proto/generate-go-tags.sh


# Prototool Compilation ------------------------------------------------------------------------------------------

SHELL := /bin/bash -o pipefail

UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)

TMP_BASE := .tmp
TMP := $(TMP_BASE)/$(UNAME_OS)/$(UNAME_ARCH)
TMP_BIN = $(TMP)/bin
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
	# cd $(PROTOTOOL_TMP); go get google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_VERSION);
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
	prototool generate
