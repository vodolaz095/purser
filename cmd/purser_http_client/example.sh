#!/usr/bin/env bash

set -e

source cmd/shared.sh
ADDR="localhost:3000"

go run cmd/purser_http_client/main.go --addr="http://${ADDR}" --token="${TOKEN}" --body="${BODY}"
