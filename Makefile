# Borrowed from: 
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = kanosdk
VET_REPORT = vet.report
GOARCH = amd64

VERSION=0.0.1
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=anzellai
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})

# Build the project
all: clean vet linux darwin windows

linux: 
	cd ${BUILD_DIR}/server; \
	GOOS=linux GOARCH=${GOARCH} go build -o ${BUILD_DIR}/bin/linux/${BINARY}-server . ; \
	cd ${BUILD_DIR}/client; \
	GOOS=linux GOARCH=${GOARCH} go build -o ${BUILD_DIR}/bin/linux/${BINARY}-client . ; \
	cd - >/dev/null

darwin:
	cd ${BUILD_DIR}/server; \
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BUILD_DIR}/bin/darwin/${BINARY}-server . ; \
	cd ${BUILD_DIR}/client; \
	GOOS=darwin GOARCH=${GOARCH} go build -o ${BUILD_DIR}/bin/darwin/${BINARY}-client . ; \
	cd - >/dev/null

windows:
	cd ${BUILD_DIR}/server; \
	GOOS=windows GOARCH=${GOARCH} go build -o ${BUILD_DIR}/bin/windows/${BINARY}-server.exe . ; \
	cd ${BUILD_DIR}/client; \
	GOOS=windows GOARCH=${GOARCH} go build -o ${BUILD_DIR}/bin/windows/${BINARY}-client.exe . ; \
	cd - >/dev/null

vet:
	-cd ${BUILD_DIR}; \
	go vet ./... > ${VET_REPORT} 2>&1 ; \
	cd - >/dev/null

fmt:
	cd ${BUILD_DIR}; \
	go fmt $$(go list ./... | grep -v /vendor/) ; \
	cd - >/dev/null

build:
	codegen; \
	make clean; \
	make fmt; \
	make vet; \
	make linux; \
	make darwin; \
	make windows; \

clean:
	-rm -f ${BUILD_DIR}/${VET_REPORT}
	-rm -rf ${BUILD_DIR}/bin

.PHONY: linux darwin windows vet fmt clean build
