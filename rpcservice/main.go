// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package main

import (
	"log"

	"github.com/cloudwego/kitex/server"

	"github.com/Trae-AI/stream-to-river/internal/suite/kitex"
	"github.com/Trae-AI/stream-to-river/rpcservice/biz/config"
	"github.com/Trae-AI/stream-to-river/rpcservice/biz/words/vocapi"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/mysql"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/redis"
	words "github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words/wordservice"
)

// main is the entry point of the RPC service application.
func main() {
	// Initialize resource like local config, mysql
	initResource()

	// Initialize the RPC server
	svr := words.NewServer(new(WordServiceImpl), server.WithSuite(kitex.NewServerSuite()))
	err := svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}

func initResource() {
	// load config
	dbConfig, lingoConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("loadConfig error: err=%v", err)
	}

	// Initialize the database with the loaded configuration
	err = mysql.InitDBWithConfig(dbConfig)
	if err != nil {
		log.Fatalf("InitDBWithConfig error: err=%v", err)
	}
	log.Printf("InitDBWithConfig success. use:%s dbConfig: %v", dbConfig.DBType, dbConfig)

	// Initialize the Lingo configuration
	vocapi.InitLingoConfig(lingoConfig)
	log.Printf("InitLingoConfig success. lingoConfig: %v", lingoConfig)

	// Initialize the cache
	redis.InitCache()
}
