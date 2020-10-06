#!/usr/bin/env sh

# abort on errors
set -e

remote_repo="https://${GITHUB_ACTOR}:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"

git config --local user.email "action@github.com"
git config --local user.name "GitHub Action"
git remote add publisher "${remote_repo}"

npm install

# build
npm run build

# navigate into the build output directory
cd ./docs

# if you are deploying to a custom domain
# echo 'www.example.com' > CNAME

git add .
git commit -m 'build: 🏗️ automatically generated documentation'
git push publisher ${BRANCH}
