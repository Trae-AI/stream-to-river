#!/bin/bash

# exit immediately if pipeline/list/(compound command) returns non-zero status
# reference https://www.gnu.org/software/bash/manual/bash.html#The-Set-Builtin
set -e

# switch to client directory
cd client

# switch node version
source /etc/profile
echo "node version is $(node -v)"

# delete node_modules directory
rm -rf node_modules

# install dependencies
npm i -g pnpm
pnpm install

# build h5
pnpm run build:h5

cd ..
mkdir -p output
mkdir -p output_resource

# Check empty files in dist directory and add comments
find client/dist -type f -empty -exec sh -c 'echo " " > "$1"' sh {} \;

# Move build artifacts to output directory
cp -R client/dist/ output/

# Move static resources to output_resource directory
cp -R client/dist/ output_resource/
