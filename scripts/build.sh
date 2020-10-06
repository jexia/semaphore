#!/usr/bin/env sh

# abort on errors
set -e

# build
npm run build

# navigate into the build output directory
cd ./dist

# if you are deploying to a custom domain
# echo 'www.example.com' > CNAME

git config --global user.name '${GITHUB_ACTOR}'
git config --global user.email '${GITHUB_ACTOR}@users.noreply.github.com'

git add .
git commit -m 'build: 🏗️ automatically generated documentation'
git push
