#!/bin/bash
set -euxo pipefail
IFS=$'\n\t'

export TEMPLATE="$HOME/repo/kubernetes.tpl.yaml"

cd ~

curl -Lo jq https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64
chmod a+x jq

export IMAGE
IMAGE="$(./jq '.builds[0].tag' -r < skaffold_output.json)"
sed "s/CONTAINER_IMAGE/$IMAGE/" < "$TEMPLATE" > retf.yaml
cat retf.yaml
