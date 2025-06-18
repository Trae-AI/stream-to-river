// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package kitex

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/server"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/config-file/filewatcher"
	fileserver "github.com/kitex-contrib/config-file/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"gopkg.in/yaml.v3"

	"github.com/Trae-AI/stream-to-river/internal/suite/env"
	"github.com/Trae-AI/stream-to-river/internal/suite/log"
)

const (
	serverDynamicConfigFile = "kitex_server.json"
)

type LoggerConfig struct {
	Level   log.Level `yaml:"Level"`
	Outputs []string  `yaml:"Outputs"`
}

type OpenTelemetryConfig struct {
	Endpoint string `yaml:"Endpoint"`
}

type ServerConfig struct {
	Address       string              `yaml:"Address"`
	Log           LoggerConfig        `yaml:"Log"`
	OpenTelemetry OpenTelemetryConfig `yaml:"OpenTelemetry"`
}

type serverSuite struct {
	serviceName string
	config      *ServerConfig
	fileWatcher filewatcher.FileWatcher
}

func NewServerSuite() server.Suite {
	yamlConfigFile := utils.GetConfFile()
	fd, err := os.Open(yamlConfigFile)
	if err != nil {
		// log the pid since there may be other process (e.g. AGW sidecar) writing to the same log file
		panic(fmt.Errorf("[pid=%v] open '%s' failed: %w", os.Getpid(), yamlConfigFile, err))
	}
	defer fd.Close()

	b, err := ioutil.ReadAll(fd)
	if err != nil {
		panic(err)
	}

	var cfg ServerConfig
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		panic(fmt.Errorf("yaml.Unmarshal '%s' failed: %w", yamlConfigFile, err))
	}

	// dynamic config
	fw, err := filewatcher.NewFileWatcher(filepath.Join(utils.GetConfDir(), serverDynamicConfigFile))
	if err != nil {
		panic(err)
	}
	if err = fw.StartWatching(); err != nil {
		panic(err)
	}

	return &serverSuite{
		serviceName: env.ServiceName(),
		config:      &cfg,
		fileWatcher: fw,
	}
}

func (s *serverSuite) Options() []server.Option {
	var opts []server.Option
	addr, _ := net.ResolveTCPAddr("tcp", s.config.Address)
	opts = append(opts, server.WithServiceAddr(addr))
	opts = append(opts, server.WithSuite(tracing.NewServerSuite()))
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: s.serviceName}))
	opts = append(opts, server.WithSuite(fileserver.NewSuite(s.serviceName, s.fileWatcher))) // add watcher

	// service registry
	r, err := etcd.NewEtcdRegistry([]string{env.ServiceRegistryETCDHost()}) // r should not be reused.
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithRegistry(r))

	// log
	logger := kitexlogrus.NewLogger()
	logger.SetLevel(s.config.Log.Level.Value())
	logger.SetOutput(os.Stdout) // default stdout
	klog.SetLogger(logger)

	// open telemetry
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(s.serviceName),
		provider.WithExportEndpoint(s.config.OpenTelemetry.Endpoint),
		provider.WithInsecure(),
	)
	return opts
}
