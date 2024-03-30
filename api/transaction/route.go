package transaction

import (
	"github.com/gabriel-tama/banking-app/api/balance"
	"github.com/gabriel-tama/banking-app/common/jwt"
	"github.com/gabriel-tama/banking-app/common/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(r *gin.RouterGroup, controller *balance.Controller, jwtService *jwt.JWTService) {
	router := r.Group("/transaction")
	router.Use(middleware.AuthorizeJWT(*jwtService))

	{
		router.POST("/", controller.ReduceBalance)
	}
}
