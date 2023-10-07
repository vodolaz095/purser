#!/usr/bin/env bash

set -ex

source cmd/shared.sh
ADDR="http://localhost:3000"

curl -v -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  --data "{\"body\":\"${BODY}\"}" \
  "${ADDR}/api/v1/secret/"
