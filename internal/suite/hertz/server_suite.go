// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package hertz

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzlogrus "github.com/hertz-contrib/obs-opentelemetry/logging/logrus"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	hertztracing "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/registry/etcd"
	"gopkg.in/yaml.v3"

	"github.com/Trae-AI/stream-to-river/internal/suite/env"
	"github.com/Trae-AI/stream-to-river/internal/suite/log"
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

type ServerSuite struct {
	serviceName string
	config      *ServerConfig
}

func NewServerSuite() *ServerSuite {
	yamlConfigFile := GetConfFile()
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
	return &ServerSuite{
		serviceName: env.ServiceName(),
		config:      &cfg,
	}
}

func (s *ServerSuite) Options() (opts []config.Option, mws []app.HandlerFunc) {
	opts = append(opts, server.WithHostPorts(s.config.Address))

	// service registry
	r, err := etcd.NewEtcdRegistry([]string{env.ServiceRegistryETCDHost()}) // r should not be reused.
	if err != nil {
		panic(err)
	}
	addr, _ := net.ResolveTCPAddr("tcp", s.config.Address)
	opts = append(opts, server.WithRegistry(r, &registry.Info{
		ServiceName: s.serviceName,
		Addr:        addr,
		Weight:      10,
		Tags:        nil,
	}))

	// log
	logger := hertzlogrus.NewLogger()
	logger.SetLevel(hlog.Level(s.config.Log.Level.Value()))
	logger.SetOutput(os.Stdout) // default stdout
	hlog.SetLogger(logger)

	// open telemetry
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(s.serviceName),
		// Support setting ExportEndpoint via environment variables: OTEL_EXPORTER_OTLP_ENDPOINT
		provider.WithExportEndpoint("localhost:4317"),
		provider.WithInsecure(),
	)

	tracer, cfg := hertztracing.NewServerTracer()
	opts = append(opts, tracer)
	mws = append(mws, hertztracing.ServerMiddleware(cfg))
	return
}
