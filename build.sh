#!/bin/bash

# Exit on error
set -e

# Hide output from pushd/popd
pushd () {
    command pushd "$@" > /dev/null
}
popd () {
    command popd "$@" > /dev/null
}

echo "* Configuring the build environment"
export VER="v2.0.5"
export ZIP="pixel-tensor-audio-decompressor-$VER.zip"
export GOOS=linux
export GOARCH=arm64
export SRCMENU="$PWD/menu"
export SRCMOD="$PWD/module"
export BINMENU="$SRCMOD/bin/menu"

echo "* Building the text menu engine"
pushd "$SRCMENU"
go build -o "$BINMENU" -ldflags="-s -w"
popd >&2

echo "* Packaging the module"
pushd "$SRCMOD" >&2
zip -r -0 -v "$ZIP" . > /dev/null
popd >&2
mv "$SRCMOD/$ZIP" "$PWD/$ZIP"

echo "* Done!"
echo "-> $ZIP"
