// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/Trae-AI/stream-to-river/apiservice/biz/rpcclient"
	"github.com/Trae-AI/stream-to-river/internal/suite/hertz"
)

// main is the entry point of the API service application.
// It initializes the configuration, observability tools, HTTP server, and RPC client,
// then registers API handlers and starts the server.
func main() {
	// Initialize the application configuration.
	initWithConfig("")

	// init HTTP server and tracing
	opts, mws := hertz.NewServerSuite().Options()
	hz := server.Default(opts...)
	hz.Use(mws...)

	// Initialize the Kitex RPC client.
	if err := rpcclient.InitRPCClient(); err != nil {
		// Log a fatal error and terminate the application if RPC client initialization fails.
		hlog.Fatal(err)
	}

	// Register API handlers with the Hertz server.
	RegisterAPI(hz)

	// Start the Hertz server and block the main goroutine.
	hz.Spin()
}
