#!/bin/bash
set -euxo pipefail
IFS=$'\n\t'

INPUT_FILE="cloudbuild/kubernetes.tpl.yaml"
OUTPUT_FILE="/tmp/kubernetes.yaml"
IMAGE="gcr.io/$PROJECT_ID/retf:$SHORT_SHA"

sed "s!CONTAINER_IMAGE!$IMAGE!g" "$INPUT_FILE" > "$OUTPUT_FILE"
cat "$OUTPUT_FILE"
kubectl apply -f "$OUTPUT_FILE"
