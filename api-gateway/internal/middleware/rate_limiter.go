package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiter *rate.Limiter
}

var (
	clients = make(map[string]*clientLimiter)
	mu      sync.Mutex
)

func getClientLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if cl, exists := clients[ip]; exists {
		return cl.limiter
	}

	// 10 requests per second, burst of 20
	limiter := rate.NewLimiter(10, 20)
	clients[ip] = &clientLimiter{limiter: limiter}
	return limiter
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getClientLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, please slow down",
			})
			return
		}

		c.Next()
	}
}