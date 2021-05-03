#!/usr/bin/env sh

set -e

npm install
npm run build

git clone -b docs "https://$GITHUB_ACTOR:$GITHUB_TOKEN@github.com/jexia/semaphore.git" semaphore

rm -rf ./semaphore/*
mv ./build/* ./semaphore

cd semaphore;

git config user.email "action@github.com"
git config user.name "Github Action"

git add -A
git commit -m 'build: 🏗️ automatically generated documentation'
git push
