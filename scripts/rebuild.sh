#!/usr/bin/env bash
set -euo pipefail

CHECK_CONFIG_SCRIPT="${CHECK_CONFIG_SCRIPT:-./scripts/check-config.sh}"
BINARY_NAME="${BINARY_NAME:-buc-server}"
COMMAND_PATH="${COMMAND_PATH:-./cmd/buc-server}"
SERVICE_NAME="${SERVICE_NAME:-$BINARY_NAME}"
OUTPUT_PATH="${OUTPUT_PATH:-bin/$BINARY_NAME}"

"$CHECK_CONFIG_SCRIPT"
go build -o "$OUTPUT_PATH" "$COMMAND_PATH"
sudo systemctl restart "$SERVICE_NAME"
sudo systemctl --no-pager --full status "$SERVICE_NAME"
