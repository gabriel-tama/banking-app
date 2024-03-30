package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) CreateUser(ctx *gin.Context) {
	var req RegisterPayload

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	data, err := c.service.Create(ctx, req)

	if errors.Is(err, ErrEmailAlreadyExists) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user create succesfully", "data": data})

}

func (c *Controller) LoginUser(ctx *gin.Context) {
	var req LoginPayload

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	data, err := c.service.FindByCredential(ctx, req)

	if errors.Is(err, ErrUserNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if errors.Is(err, ErrWrongPassword) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user logged in", "data": data})

}

func (c *Controller) PingDB(ctx *gin.Context) {
	err := c.service.PingDB(ctx)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"service": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"service": "ok"})
}
