#!/bin/sh

set -e

function help() {
    echo "buildTool helps developers to powerful.

Usage:
  buildtool [flags]
  buildtool [command] [apkpath] [apkfile]

Available Commands:
  b|build       Build a package with debug mode
  n|native      Build jni lib with CGO
  h|help        Help about any command

Flags:
  -h, --help   help for buildTool
"
    exit 0
}


function build() {
    # GOOS=darwin|linux|windows
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o example_darwin
    # CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o example_linux
    # CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o example_windows

    swag init
    rm -f docs/docs.go
}

case "$1" in
    b|build)
    build -d
    ;;
    p|package)
    package
    ;;
    *)
    help
esac

