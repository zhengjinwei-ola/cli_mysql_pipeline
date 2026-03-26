#!/bin/bash
set -e

VERSION=${1:-$(git rev-parse --short HEAD 2>/dev/null || echo "dev")}
OUTPUT_DIR="dist"

echo "=== Building cli_mysql_pipeline @ ${VERSION} ==="

mkdir -p "${OUTPUT_DIR}"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -trimpath \
  -ldflags "-s -w -X main.Version=${VERSION}" \
  -o "${OUTPUT_DIR}/cli_mysql_pipeline" \
  ./instance/partystar/

cp configs/partystar.yaml "${OUTPUT_DIR}/"

echo "=== Build complete: ${OUTPUT_DIR}/cli_mysql_pipeline ==="
