#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

cd retf-env
git checkout candidate
git add kubernetes.yaml
git commit -m "Deploying from retf-app @ ${SHORT_SHA}"
git push origin candidate
