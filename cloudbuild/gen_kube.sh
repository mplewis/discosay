#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

INPUT_FILE="cloudbuild/kubernetes.tpl.yaml"
OUTPUT_FILE="retf-env/kubernetes.yaml"
IMAGE="gcr.io/${PROJECT_ID}/retf:${SHORT_SHA}"

sed "s!CONTAINER_IMAGE!${IMAGE}!g" "$INPUT_FILE" > "$OUTPUT_FILE"
