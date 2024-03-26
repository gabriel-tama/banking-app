package transaction

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

func (c *Controller) AddBalance(ctx *gin.Context) {
	var req AddBalancePayload

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := c.service.Add(ctx, req)

	if errors.Is(err, ErrValidationFailed) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if errors.Is(err, ErrInvalidToken) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "succesfully add balance"})

}

func (c *Controller) GetBalance(ctx *gin.Context) {
	data, err := c.service.GetBalance(ctx)
	if errors.Is(err, ErrInvalidToken) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ok gas", "data": data})

}

func (c *Controller) ReduceBalance(ctx *gin.Context) {
	var req ReduceBalancePayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := c.service.Reduce(ctx, req)

	if errors.Is(err, ErrValidationFailed) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if errors.Is(err, ErrInvalidToken) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ok gas"})
}

func (c *Controller) GetHistory(ctx *gin.Context) {
	var req GetHistoryPayload
	var pagination Pagination
	if err := ctx.ShouldBind(&req); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	paramPairs := ctx.Request.URL.Query()
	for _, values := range paramPairs {
		if values[0] == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}
	}

	data, total, err := c.service.GetHistory(ctx, req)

	if errors.Is(err, ErrInvalidToken) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
		return
	}
	pagination.Limit = req.Limit
	pagination.Offset = req.Offset
	pagination.Total = total
	ctx.JSON(http.StatusOK, gin.H{"message": "ok", "data": data, "meta": pagination})
}
