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
.PHONY: pb
pb:
	protoc -I protobuf/ protobuf/ghost/ghost.proto --go_out=protobuf/
	protoc -I protobuf/ protobuf/client/client.proto --go_out=protobuf/

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

