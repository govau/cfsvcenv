#!/usr/bin/env sh

set -eu

GOOS=linux GOARCH=amd64 go build -tags cloudfoundry
