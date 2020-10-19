#!/usr/bin/env sh

cd docs

# abort on errors
set -e

git config --local user.email "action@github.com"
git config --local user.name "GitHub Action"

npm install

# build
npm run build

# if you are deploying to a custom domain
# echo 'www.example.com' > CNAME

# force add due to ignored in .gitignore
git add -f ./build

git commit -m 'build: 🏗️ automatically generated documentation'
git subtree push --force --prefix build origin docs
