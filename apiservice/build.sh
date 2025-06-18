#!/bin/bash
RUN_NAME="dk.stream2river.api"

CURDIR=$(cd $(dirname $0); pwd)
cd ${CURDIR}

mkdir -p output/bin output/conf
cp script/bootstrap.sh output 2>/dev/null
#  find conf/ -type f ! -name "*_local.*" | xargs -I{} cp {} output/conf/
cp conf/* output/conf/
chmod +x output/bootstrap.sh

go mod tidy
go build -o output/bin/${RUN_NAME}
