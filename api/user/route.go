package user

import (
	"github.com/gabriel-tama/banking-app/common/jwt"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.RouterGroup, controller *Controller, jwtService *jwt.JWTService) {
	router := r.Group("/user")

	{
		router.POST("/register", controller.CreateUser)
		router.POST("/login", controller.LoginUser)
		router.GET("/health", controller.PingDB)
	}
}
