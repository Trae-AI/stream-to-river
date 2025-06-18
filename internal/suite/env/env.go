// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package env

import "os"

const (
	envServiceRegistryETCDHost = "SERVICE_REGISTRY_ETCD_HOST"
	envServiceName             = "MICRO_SERVICE_NAME"
)

var (
	serviceRegistryETCDHost string
	serviceName             string
)

func init() {
	serviceRegistryETCDHost = "127.0.0.1:2379"
	if h := os.Getenv(envServiceRegistryETCDHost); h != "" {
		serviceRegistryETCDHost = h
	}

	serviceName = os.Getenv(envServiceName)
}

func ServiceRegistryETCDHost() string {
	return serviceRegistryETCDHost
}

func ServiceName() string {
	return serviceName
}
