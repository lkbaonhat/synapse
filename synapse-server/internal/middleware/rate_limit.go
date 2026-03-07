package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func init() {
	go cleanupVisitors()
}

func getVisitor(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	v, exists := visitors[ip]
	if !exists {
		v = &visitor{limiter: rate.NewLimiter(r, b)}
		visitors[ip] = v
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

// RateLimit allows `rps` requests per second per IP with a burst of `burst`.
func RateLimit(rps float64, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := getVisitor(c.ClientIP(), rate.Limit(rps), burst)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}
		c.Next()
	}
}
