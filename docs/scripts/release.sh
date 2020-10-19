#!/usr/bin/env sh

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
git add -f ./dist

git commit -m 'build: 🏗️ automatically generated documentation'
git subtree push --prefix dist origin gh-pages
