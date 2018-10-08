#!/bin/bash

APP_NAME=system-status-service
APP_VERSION=1.0
DIST_DIR=dist
BUILD_TIME=$(date "+%F %T %Z")

WINDOWS_ARCH=386,amd64
LINUX_ARCH=386,amd64,arm

gobuild() {
	echo "building $1 $2"
	ext=""
	append_suffix=""
	if [ "$1" == "windows" ]; then
		ext=".exe"
	fi
	if [ -n "$exe_suffix" ]; then
		append_suffix=_$1_$2
	fi
	target_dir=$DIST_DIR/$1/$2
	rm -rf $target_dir
	GOOS=$1 GOARCH=$2 go build -o $target_dir/${APP_NAME}${append_suffix}${ext} \
	-ldflags \
	"\
	-X 'main.appName=${APP_NAME}' \
	-X 'main.appVersion=${APP_VERSION}' \
	-X 'main.buildTime=${BUILD_TIME}' \
	" \
	.
}

showhelp() {
	echo "Usage: build.sh [-m] [-w] [-l] [-s]"
	echo "    -m  build macos executable of amd64"
	echo "    -w  build windows executables of all arch ($WINDOWS_ARCH)"
	echo "    -w[=<arch>,...]  build windows executables of specific arch"
	echo "    -l  build linux executables of all arch ($LINUX_ARCH)"
	echo "    -l[=<arch>,...]  build linux executables of specific arch"
	echo "    -s  append os type and arch suffix of executable name (use 'foo_linux_amd64' instead of 'foo')"
}

archContains() {
	str=$1
	array=(${str//,/ })
	for var in ${array[@]}
	do
		if [ "$var" == "$2" ]; then
			echo "true"
			return
		fi
	done
}

cd "$( dirname "$0" )"
export CGO_ENABLED=0

if [ $# -gt 0 ]; then
	for arg in $*
	do
		case $arg in
			-m)
				build_mac=1
			;;
			-w)
				build_windows=$WINDOWS_ARCH
			;;
			-w=*)
				build_windows=${arg#*-w=}
			;;
			-l)
				build_linux=$LINUX_ARCH
			;;
			-l=*)
				build_linux=${arg#*-l=}
			;;
			-s)
				exe_suffix=1
			;;
			*)
				echo "unknow arg: $arg"
				showhelp
				exit 1
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
	array=(${build_linux//,/ })
	for var in ${array[@]}
	do
		cont=$(archContains $LINUX_ARCH $var)
		if [ -n "$cont" ]; then
			gobuild linux $var
		else
			echo "unknow arch $var, skip"
		fi
	done
fi
if [ -n "$build_windows" ]; then
	array=(${build_windows//,/ })
	for var in ${array[@]}
	do
		cont=$(archContains $WINDOWS_ARCH $var)
		if [ -n "$cont" ]; then
			gobuild windows $var
		else
			echo "unknow arch $var, skip"
		fi
	done
	gobuild windows amd64
	gobuild windows 386
fi
