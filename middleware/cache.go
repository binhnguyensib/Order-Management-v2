package middleware

import (
	"intern-project-v2/config"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-gonic/gin"
)

func CacheMiddleware(duration time.Duration, handler gin.HandlerFunc) gin.HandlerFunc {
	return cache.CachePage(config.CacheStore, duration, handler)
}
