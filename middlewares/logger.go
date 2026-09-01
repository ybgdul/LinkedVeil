package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func StartLogging(c *gin.Context) {
	log.Println("About to start processing a request")
	start := time.Now()
	c.Set("request_start", start)
	c.Next()

	log.Println("Request ended after: ", time.Since(start))
}