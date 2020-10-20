#!/usr/bin/env sh

# abort on errors
set -e

npm install
npm run build

git config --global user.email "action@github.com"
git config --global user.name "Github Action"

git clone -b docs "https://$GITHUB_ACTOR:$GITHUB_TOKEN@github.com/jexia/semaphore.git" semaphore

rm -rf ./semaphore/*
mv ./build/* ./semaphore

cd semaphore;

if [ -z $(git status --porcelain) ];
then
    echo "docs up to date nothing to commit"
else
    echo "comitting latest changes"
    git add -A
    git commit -m 'build: 🏗️ automatically generated documentation'git push
fi

cd ..
