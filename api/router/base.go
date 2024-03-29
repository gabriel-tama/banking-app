package router

import (
	"log"
	"time"

	"github.com/gabriel-tama/banking-app/api/balance"
	"github.com/gabriel-tama/banking-app/api/image"
	"github.com/gabriel-tama/banking-app/api/transaction"
	"github.com/gabriel-tama/banking-app/api/user"
	"github.com/gabriel-tama/banking-app/common/jwt"
	"github.com/gabriel-tama/banking-app/common/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/ratelimit"
)

var (
	limit ratelimit.Limiter
)

type RouterParam struct {
	JwtService            *jwt.JWTService
	UserController        *user.Controller
	ImageController       *image.ImageController
	TransactionController *transaction.Controller
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
	router.Use(middleware.GinPrometheusMiddleware())
	router.Use(leakBucket())
	router.Use(gin.Recovery())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/healthz", param.UserController.PingDB)

	// Setup API version 1 routes
	v1 := router.Group("/v1")
	{
		user.NewRouter(v1, param.UserController, param.JwtService)
		transaction.NewRouter(v1, param.TransactionController, param.JwtService)
		balance.NewRouter(v1, param.TransactionController, param.JwtService)

	}

	router.GET("/rate", func(c *gin.Context) {
		c.JSON(200, "rate limiting test")
	})

	return router
}
