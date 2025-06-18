#!/bin/bash

kitex_install() {
    echo "installing latest kitex tool ..."
    # TODO fix to use release version
    go install github.com/cloudwego/kitex/tool/cmd/kitex@feat/kitex_tool_http_annotation
    echo "installing kitex tool ... done"
    kitex -version
}

kitex_install

# update scaffolding code for rpcservice
cd rpcservice
kitex -module github.com/coze-dev/stream-to-river -service dk.stream2river.word -thrift no_default_serdes -streamx ../idl/words.thrift
