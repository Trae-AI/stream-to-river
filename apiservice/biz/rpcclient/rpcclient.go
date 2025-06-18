// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package rpcclient

import (
	"github.com/cloudwego/kitex/client"

	"github.com/Trae-AI/stream-to-river/internal/suite/kitex"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words/wordservice"
)

// dest is the destination service name for the RPC client.
// It specifies the target service that the RPC client will communicate with.
var dest = "dk.stream2river.word"

// WordsRPCCli is a global client instance for the wordservice.
// It is used to make RPC calls to the word - related services defined in the wordservice.
var WordsRPCCli wordservice.Client

// InitRPCClient initializes the RPC client for the wordservice.
// It first sets up the service governance configuration by creating a client suite,
// then creates a new client instance with the specified destination and client suite.
// Finally, it assigns the newly created client instance to the global variable WordsRPCCli.
//
// Returns:
//   - error: An error object if an unexpected error occurs during client initialization.
func InitRPCClient() error {
	// Create a new client instance for the wordservice with the specified destination and client suite.
	c, err := wordservice.NewClient(dest, client.WithSuite(kitex.NewClientSuite(dest)))
	if err != nil {
		// Return the error if client creation fails.
		return err
	}

	// Assign the new client instance to the global variable
	WordsRPCCli = c
	return nil
}
