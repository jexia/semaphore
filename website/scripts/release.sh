#!/usr/bin/env sh

# abort on errors
set -e

npm install
npm run build

git config --global user.email "action@github.com"
git config --global user.name "GitHub Action"

git clone -b docs https://github.com/jexia/semaphore.git semaphore

(cd semaphore; rm -r *; mv ../build/* ./)
(cd semaphore; git add -A; git commit -m 'build: 🏗️ automatically generated documentation'; git push)
