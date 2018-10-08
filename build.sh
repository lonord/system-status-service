#!/bin/bash

APP_NAME=system-status-service
APP_VERSION=1.0
BUILD_TIME=$(date "+%F %T %Z")

DIST_DIR=dist

gobuild() {
	echo "building $1 $2"
	ext=""
	if [ "$1" == "windows" ]; then
		ext=".exe"
	fi
	target_dir=$DIST_DIR/$1/$2
	rm -rf $target_dir
	GOOS=$1 GOARCH=$2 go build -o $target_dir/${APP_NAME}_$1_$2${ext} \
	-ldflags \
	"\
	-X 'main.appName=${APP_NAME}' \
	-X 'main.appVersion=${APP_VERSION}' \
	-X 'main.buildTime=${BUILD_TIME}' \
	" \
	.
}

showhelp() {
	echo "Usage: build.sh [-w] [-m] -[l]"
	echo "    -w  build windows executable"
	echo "    -m  build macos executable"
	echo "    -l  build linux executable"
}

cd "$( dirname "$0" )"
export CGO_ENABLED=0

if [ $# -gt 0 ]; then
	for arg in $*
	do
		case $arg in
			-w)
				build_windows=1
			;;
			-m)
				build_mac=1
			;;
			-l)
				build_linux=1
			;;
		esac
	done
else
	showhelp
	exit 0
fi

if [ -n "$build_mac" ]; then
	gobuild darwin amd64
fi
if [ -n "$build_linux" ]; then
	gobuild linux amd64
	gobuild linux 386
	gobuild linux arm
fi
if [ -n "$build_windows" ]; then
	gobuild windows amd64
	gobuild windows 386
fi
