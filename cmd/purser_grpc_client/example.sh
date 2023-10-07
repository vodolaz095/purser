#!/usr/bin/env bash

set -e

source cmd/shared.sh
ADDR="localhost:3001"

go run cmd/purser_grpc_client/main.go --addr="${ADDR}" --token="${TOKEN}" --body="${BODY}"
