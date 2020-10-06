#!/usr/bin/env sh

# abort on errors
set -e

npm install

# build
npm run build

# navigate into the build output directory
cd ./docs

# if you are deploying to a custom domain
# echo 'www.example.com' > CNAME

git add .
git commit -m 'build: 🏗️ automatically generated documentation'
git push
