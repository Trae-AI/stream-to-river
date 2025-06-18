#! /bin/bash

set -e

CURDIR=$(cd $(dirname $0); pwd)
cd ${CURDIR}

# build rpc service
docker build -f Dockerfile-rpc -t rpcservice ../

# build api service
docker build -f Dockerfile-api -t apiservice ../

docker-compose up -d
