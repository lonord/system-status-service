#!/bin/bash

VERSION=1.0
APP_NAME=sss
PACKAGE_NAME=sys-status-service
BUILD_TIME=$(date "+%F %T")

cd $(dirname $0)

if [ ! -e "dist" ]; then
	mkdir dist
fi

gobuild() {
	if [ -e "dist/tmp" ]; then
		rm -rf dist/tmp
		mkdir dist/tmp
	fi
	go build -o dist/tmp/usr/local/bin/$APP_NAME \
	-ldflags \
	"\
	-X 'main.appVersion=${VERSION}' \
	-X 'main.buildTime=${BUILD_TIME}' \
	" \
	.
}

build_deb() {
	gobuild
	rm -rf dist/linux/$1
	mkdir -p dist/linux/$1
	fpm -s dir -t deb -a $1 -p dist/linux/$1/${PACKAGE_NAME}_v${VERSION}_$1.deb -n $APP_NAME -v ${VERSION} -C dist/tmp .
}

# build deb x64
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
build_deb amd64

# build deb armhf
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=arm
build_deb armhf
