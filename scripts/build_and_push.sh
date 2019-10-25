#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

cd ~

echo "$GCP_TOKEN" > gcp_token.json
export GOOGLE_APPLICATION_CREDENTIALS="$HOME/gcp_token.json"

curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/v0.40.0/skaffold-linux-amd64
chmod +x skaffold
./skaffold run --quiet > skaffold_output.json
cat skaffold_output.json
