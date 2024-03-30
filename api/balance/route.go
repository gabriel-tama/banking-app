package balance

import (
	"github.com/gabriel-tama/banking-app/common/jwt"
	"github.com/gabriel-tama/banking-app/common/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.RouterGroup, controller *Controller, jwtService *jwt.JWTService) {
	router := r.Group("/balance")
	router.Use(middleware.AuthorizeJWT(*jwtService))

	{
		router.POST("/", controller.AddBalance)
		router.GET("/", controller.GetBalance)
		router.GET("/history", controller.GetHistory)
	}
}
