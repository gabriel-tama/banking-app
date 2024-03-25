package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gabriel-tama/banking-app/api/image"
	"github.com/gabriel-tama/banking-app/api/router"
	"github.com/gabriel-tama/banking-app/api/user"
	C "github.com/gabriel-tama/banking-app/common/config"
	psql "github.com/gabriel-tama/banking-app/common/db"
	"github.com/gabriel-tama/banking-app/common/jwt"
	"github.com/gin-gonic/gin"
)

func main() {

	env, err := C.Get()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, dbErr := psql.Init(context.Background())
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close(context.Background())

	// Repository
	userRepository := user.NewRepository(db, env.BCRYPT_Salt)

	// Service
	jwtService := jwt.NewJWTService(env.JWTSecret, env.JWTExp)
	userService := user.NewService(userRepository, jwtService)
	s3Service := image.NewS3Service(env.S3ID, env.S3Secret, env.S3Bucket, env.S3Url, env.S3Region)

	// Controller
	userController := user.NewController(userService)
	imgController := image.NewImageController(s3Service)

	router := router.SetupRouter(router.RouterParam{
		JwtService:      &jwtService,
		ImageController: imgController,
		UserController:  userController,
	})

	router.GET("/v1/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := router.Run(fmt.Sprintf("%s:%s", env.Host, env.Port)); err != nil {
		log.Fatal("Server error:", err)
	}
}
