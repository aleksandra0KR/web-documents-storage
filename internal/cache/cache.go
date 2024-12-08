package cache

import (
	"github.com/dgraph-io/ristretto"
)

var Cache *ristretto.Cache

func InitializeCache() {
	Cache, _ = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
}
