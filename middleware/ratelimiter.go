package middleware

import (
	"intern-project-v2/logger"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Client struct {
	Count     int
	ResetTime time.Time
	LastSeen  time.Time
}

var (
	clients = make(map[string]*Client)
	mutex   = sync.RWMutex{}
)

func startCleanup() {
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			mutex.Lock()
			now := time.Now()
			for ip, client := range clients {
				if now.Sub(client.LastSeen) > 10*time.Minute {
					delete(clients, ip)
				}
			}
			mutex.Unlock()
			logger.Info("Cleanup completed: old clients removed")
		}
	}()
}

func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	static := false
	if !static {
		startCleanup()
		static = true
	}
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()

		mutex.Lock()
		defer mutex.Unlock()

		client, exists := clients[clientIP]

		if !exists {
			clients[clientIP] = &Client{
				Count:     1,
				ResetTime: now.Add(1 * time.Minute),
				LastSeen:  now,
			}
			c.Next()
			return
		}

		client.LastSeen = now
		if now.After(client.ResetTime) {
			client.Count = 1
			client.ResetTime = now.Add(1 * time.Minute)
			c.Next()
			return
		}
		if client.Count > requestsPerMinute {
			secondsLeft := int(time.Until(client.ResetTime).Seconds())

			c.JSON(429, gin.H{
				"error":       "Too many requests",
				"retry_after": secondsLeft})
			logger.Warn("Rate limit exceeded", "IP", clientIP, "Count", client.Count, "ResetTime", client.ResetTime)
			c.Abort()
			return
		}
		client.Count++
		c.Next()
	}

}
