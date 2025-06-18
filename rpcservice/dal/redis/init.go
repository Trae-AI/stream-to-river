// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package redis

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Cache is a global variable that holds the in - memory cache instance.
// It's used to store frequently accessed data to reduce database queries.
var Cache *cache.Cache

// CacheExpireTime defines the default expiration time for cache items.
// Items in the cache will be automatically deleted after this duration.
var CacheExpireTime = 24 * time.Hour

// CacheCleanupTime determines how often the cache performs cleanup operations.
// The cache will remove expired items at this interval.
var CacheCleanupTime = 72 * time.Hour

// InitCache initializes the in - memory cache with the specified expiration and cleanup times.
func InitCache() {
	Cache = cache.New(CacheExpireTime, CacheCleanupTime)
}
