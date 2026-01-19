package cache

import "errors"

// ErrCacheMiss is returned when a key is not found in the cache
var ErrCacheMiss = errors.New("cache miss")
