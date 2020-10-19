#!/usr/bin/env sh

cd website

# abort on errors
set -e

# https://v2.docusaurus.io/docs/deployment#deploying-to-github-pages
export GIT_USER="github-actions"
export DEPLOYMENT_BRANCH="docs"
export CURRENT_BRANCH="master"

npm install
npm run deploy
