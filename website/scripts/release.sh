#!/usr/bin/env sh

# abort on errors
set -e

git config --global user.email "actions@github.com"
git config --global user.name "Github Actions"

# https://v2.docusaurus.io/docs/deployment#deploying-to-github-pages
export GIT_USER="github-actions"
export DEPLOYMENT_BRANCH="docs"
export CURRENT_BRANCH="master"

npm install
npm run deploy
