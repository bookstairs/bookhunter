#!/bin/bash

if [[ $1 == *"linux"* ]]; then
  upx --lzma "$1"
elif [[ $1 == *"window"* ]]; then
  upx -9 "$1"
fi
