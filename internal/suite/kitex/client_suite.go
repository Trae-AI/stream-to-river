// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package kitex

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/connpool"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/utils"
	fileclient "github.com/kitex-contrib/config-file/client"
	"github.com/kitex-contrib/config-file/filewatcher"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/Trae-AI/stream-to-river/internal/suite/env"
)

const clientDynamicConfigFile = "kitex_client.json"

type clientSuite struct {
	dest string
}

var defaultConnPoolConfig = connpool.IdleConfig{
	MaxIdlePerAddress: 10,
	MaxIdleGlobal:     1000,
	MaxIdleTimeout:    60 * time.Second,
	MinIdlePerAddress: 2,
}

var clientFileWatcher filewatcher.FileWatcher
var clientFileWatcherOnce sync.Once

func NewClientSuite(dest string) client.Suite {
	clientFileWatcherOnce.Do(func() {
		fw, err := filewatcher.NewFileWatcher(filepath.Join(utils.GetConfDir(), clientDynamicConfigFile))
		if err != nil {
			panic(err)
		}
		// start watching file changes
		if err = fw.StartWatching(); err != nil {
			panic(err)
		}
		clientFileWatcher = fw
	})
	return &clientSuite{
		dest: dest,
	}
}

func (s *clientSuite) Options() []client.Option {
	var opts []client.Option
	opts = append(opts, client.WithLongConnection(defaultConnPoolConfig))
	opts = append(opts, client.WithSuite(fileclient.NewSuite(s.dest,
		fmt.Sprintf("%s/%s", env.ServiceName(), s.dest), clientFileWatcher))) // timeout, retry, circuit break
	opts = append(opts, client.WithSuite(tracing.NewClientSuite())) // tracing
	opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: s.dest}))

	// service discovery
	r, err := etcd.NewEtcdResolver([]string{env.ServiceRegistryETCDHost()})
	if err != nil {
		panic(err)
	}
	opts = append(opts, client.WithResolver(r))
	return opts
}
