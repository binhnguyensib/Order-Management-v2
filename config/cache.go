package config

import (
	"time"

	"github.com/gin-contrib/cache/persistence"
)

var CacheStore persistence.CacheStore

func InitCache() {
	CacheStore = persistence.NewInMemoryStore(time.Minute * 5)
}
