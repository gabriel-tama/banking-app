package router

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

var (
	limit ratelimit.Limiter
)

type RouterParam struct {
}

func leakBucket() gin.HandlerFunc {
	prev := time.Now()
	return func(ctx *gin.Context) {
		now := limit.Take()
		log.Printf("%v", now.Sub(prev))
		prev = now
	}
}

func SetupRouter(param RouterParam) *gin.Engine {
	limit = ratelimit.New(1000)
	router := gin.Default()

	router.SetTrustedProxies([]string{"::1"}) // This is for reverse proxy

	router.Use(leakBucket())
	router.Use(gin.Recovery())

	// Setup API version 1 routes
	// v1 := router.Group("/v1")
	// {
	// }

	router.GET("/rate", func(c *gin.Context) {
		c.JSON(200, "rate limiting test")
	})

	return router
}
