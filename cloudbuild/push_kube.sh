#!/bin/bash
set -euxo pipefail
IFS=$'\n\t'

cd retf-env
git add kubernetes.yaml
git commit -m "Deploying from retf-app @ ${COMMIT_SHA}"
git push origin candidate
