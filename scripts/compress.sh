#!/bin/bash

# We will add upx support for macOS 13, after the https://github.com/upx/upx/issues/612 issue get fixed.

if [[ $1 == *"linux"* ]]; then
  upx --lzma "$1"
elif [[ $1 == *"window"* ]]; then
  upx -9 "$1"
fi
