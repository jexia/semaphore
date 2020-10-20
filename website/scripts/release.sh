#!/usr/bin/env sh

npm install
npm run build

git config --global user.email "action@github.com"
git config --global user.name "Github Action"

git clone -b docs "https://$GITHUB_ACTOR:$GITHUB_TOKEN@github.com/jexia/semaphore.git" semaphore

rm -rf ./semaphore/*
mv ./build/* ./semaphore

cd semaphore;

git add -A
git commit -m 'build: 🏗️ automatically generated documentation'
git push
