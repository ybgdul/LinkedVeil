package middlewares

import (
	"linkedveil/models"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimitingMiddleware(ipRateLimiter *models.IPRateLimiter) gin.HandlerFunc { 
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil { 
			c.JSON(http.StatusInternalServerError, gin.H{"error" : "Invalid IP"})
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		limiter := ipRateLimiter.GetLimiter(ip)
		if limiter.Allow() {
			c.Next()
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{"error" : "too many requests"})
		}
	}
}