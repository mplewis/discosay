#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

git config --global user.name "$(cat cloudbuild/committer_name)"
git config --global user.email "$(cat cloudbuild/committer_email)"
mv cloudbuild/known_hosts /root/.ssh/known_hosts
chmod 600 /root/.ssh/id_rsa
cat << EOF > /root/.ssh/config
Hostname github.com
IdentityFile /root/.ssh/id_rsa
EOF

git clone git@github.com:mplewis/retf-env.git
cd retf-env
git checkout candidate
