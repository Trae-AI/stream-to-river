#!/bin/bash

# Please execute gen_code.sh in the rpcservice

kitex_install() {
    echo "installing latest kitex tool ..."
    # TODO fix to use release version
    go install github.com/cloudwego/kitex/tool/cmd/kitex@feat/kitex_tool_http_annotation
    echo "installing kitex tool ... done"
    kitex -version
}

kitex_install

# Update the generated code in kitex_gen. When you modify the idl file, you need to update the generated code
kitex -module github.com/Trae-AI/stream-to-river -service dk.stream2river.word -thrift no_default_serdes -streamx ../idl/words.thrift
