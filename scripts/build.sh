#!/usr/bin/env bash

set -euo pipefail

GOOS="linux" go build -ldflags='-s -w' -o bin/main watchexec/cmd/main

ln -fs main bin/build
ln -fs main bin/detect