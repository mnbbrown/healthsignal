# ########################################################## #
# Makefile for Golang Project
# Includes cross-compiling, installation, cleanup
# ########################################################## #

# Check for required command tools to build or stop immediately
EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=healthsignal-client
VERSION=1.0.1
BUILD=`git rev-parse HEAD`
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64
GOARMS=6 7

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

default: build

all: clean build_all

run: build
	__dist/healthsignal-client

setup:
	mkdir -p __dist

build: setup
	go build ${LDFLAGS} -o __dist/${BINARY}

build_all: setup
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o __dist/$(BINARY)-$(GOOS)-$(GOARCH))))
	$(foreach GOARM, $(GOARMS), $(shell export GOARCH=arm; export GOARM=$(GOARM); go build -v -o __dist/$(BINARY)-arm$(GOARM)))

deps:
	go get ./...

install:
	go install ${LDFLAGS}

# Remove only what we've created
clean:
	find __dist/${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

.PHONY: check clean install build_all all
