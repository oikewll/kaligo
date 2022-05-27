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
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -trimpath -o example_darwin_arm64
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o example_darwin_amd64
    # CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o example_linux_amd64
    # CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o example_windows_amd64

    # UPX压缩
    # brew install --build-from-source upx
    upx -9 example_darwin_arm64
    upx -9 example_darwin_amd64

    swag init
    rm -f docs/docs.go
    rm -f docs/swagger.yaml
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

