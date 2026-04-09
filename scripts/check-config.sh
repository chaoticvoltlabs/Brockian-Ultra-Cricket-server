#!/usr/bin/env bash
set -euo pipefail

for f in config/*.json; do
  echo "Checking $f"
  jq empty "$f"
done

echo "All config files are valid."
