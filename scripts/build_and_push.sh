#!/bin/bash
set -euxo pipefail
IFS=$'\n\t'

cd ~

echo "$GCP_TOKEN" > gcp_token.json
export GOOGLE_APPLICATION_CREDENTIALS="$HOME/gcp_token.json"
gcloud auth activate-service-account --key-file="$HOME/gcp_token.json"
gcloud --quiet config set project "$GCP_PROJECT"
gcloud --quiet config set compute/zone "$GCP_ZONE"

curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/v0.40.0/skaffold-linux-amd64
chmod +x skaffold
(
  cd ~/repo
  ~/skaffold build --quiet > ~/skaffold_output.json
)
cat skaffold_output.json
