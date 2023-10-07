#!/usr/bin/env bash

set -e

ADDR="localhost:3001"
TOKEN="eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ2b2RvbGF6MDk1IiwiaWF0IjoxNTE2MjM5MDIzLCJleHAiOjE4MTYyMzkwMjN9.ga-y5GLMIJaoO6-UPY3eTztUQjGNj1QD_7yS0oggse7MKoyUaEZDMcO8ADRm7m6F1oWvSWBpu6hoPfLCQ64Emg"
BODY="example"

go run main.go --addr="${ADDR}" --token="${TOKEN}" --body="${BODY}"
