#!/bin/bash
export HERTZTOOL_VERSION=v3.4.2

CURDIR=$(cd $(dirname $0); pwd)

export MICRO_SERVICE_NAME=dk.stream2river.api

echo "$CURDIR/bin/${MICRO_SERVICE_NAME}"

exec $CURDIR/bin/${MICRO_SERVICE_NAME}