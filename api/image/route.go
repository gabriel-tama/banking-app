package image

import (
	"github.com/gabriel-tama/banking-app/common/jwt"
	"github.com/gabriel-tama/banking-app/common/middleware"
	"github.com/gin-gonic/gin"
)

func NewImageRouter(r *gin.RouterGroup, controller *ImageController, jwtService *jwt.JWTService) {
	router := r.Group("/image")
	router.Use(middleware.AuthorizeJWT(*jwtService))

	{
		router.POST("/", controller.UploadImage)
	}
}
