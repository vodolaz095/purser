#!/usr/bin/env bash

set -e

ADDR="localhost:3000"
TOKEN="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2b2RvbGF6MDk1IiwiaWF0IjoxNTE2MjM5MDIzLCJleHAiOjE4MTYyMzkwMjN9.ga-y5GLMIJaoO6-UPY3eTztUQjGNj1QD_7yS0oggse7MKoyUaEZDMcO8ADRm7m6F1oWvSWBpu6hoPfLCQ64Emg"
BODY="example"


curl -v -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  --data "{\"body\":\"${BODY}\"}" \
  "${ADDR}/api/v1/secret/"
