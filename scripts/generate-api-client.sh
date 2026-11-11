#!/usr/bin/env bash
set -euo pipefail

# Generate Kotlin Multiplatform API client from OpenAPI spec
# Run this after backend NestJS changes that affect the OpenAPI spec

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

OPENAPI_JSON="$REPO_ROOT/packages/backend-nestjs/openapi.json"
OUTPUT_DIR="$REPO_ROOT/packages/api-client"
GENERATOR_VERSION="6.6.0"

echo "=== AREA API Client Generator ==="
echo "OpenAPI spec: $OPENAPI_JSON"
echo "Output: $OUTPUT_DIR"

if [ ! -f "$OPENAPI_JSON" ]; then
  echo "ERROR: OpenAPI JSON not found at $OPENAPI_JSON"
  echo "Run the NestJS backend first to generate the spec, or run 'pnpm run build' in packages/backend-nestjs"
  exit 1
fi

mkdir -p "$OUTPUT_DIR"

# Generate using OpenAPI Generator via Docker
docker run --rm \
  -v "$REPO_ROOT:$REPO_ROOT" \
  -w "$REPO_ROOT" \
  "openapitools/openapi-generator-cli:$GENERATOR_VERSION" \
  generate \
  -i "$OPENAPI_JSON" \
  -g kotlin \
  -p kotlinPackage=area.api \
  -p apiPackage=area.api.client \
  -p modelPackage=area.api.models \
  -p dateTimeType=java.time.LocalDateTime \
  -o "$OUTPUT_DIR"

echo "=== API client generation complete ==="
echo "Generated files in $OUTPUT_DIR"
