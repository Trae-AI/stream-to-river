#!/usr/bin/env bash
RUN_NAME="dk.stream2river.word"

CURDIR=$(cd $(dirname $0); pwd)
cd ${CURDIR}

mkdir -p output/bin output/conf output/data
cp script/* output/
cp -r conf/* output/conf/
chmod +x output/bootstrap.sh

go mod tidy

cp tools/stream_sqlite.sql output/data/
cd tools/init_default_user && go build && cd - && cp tools/init_default_user/init_default_user output/data/

if [ "$IS_SYSTEM_TEST_ENV" != "1" ]; then
    go build -o output/bin/${RUN_NAME}
else
    go test -c -covermode=set -o output/bin/${RUN_NAME} -coverpkg=./...
fi
