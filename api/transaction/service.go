package transaction

import (
	"fmt"

	"github.com/gabriel-tama/banking-app/common/jwt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	Add(ctx *gin.Context, req AddBalancePayload) error
	Reduce(ctx *gin.Context, req ReduceBalancePayload) error
	GetBalance(ctx *gin.Context) (*GetBalanceResponse, error)
	GetHistory(ctx *gin.Context, req GetHistoryPayload) (*GetHistoryResponse, int, error)
}

type transactionService struct {
	repository Repository
	jwtService jwt.JWTService
}

func NewService(repository Repository, jwtService jwt.JWTService) Service {
	return &transactionService{repository: repository, jwtService: jwtService}
}

func (s *transactionService) Add(ctx *gin.Context, req AddBalancePayload) error {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return ErrValidationFailed
	}
	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return ErrInvalidToken
	}

	return s.repository.Add(ctx, &req, token.UserID)

}

func (s *transactionService) Reduce(ctx *gin.Context, req ReduceBalancePayload) error {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		fmt.Println(err)
		return ErrValidationFailed
	}

	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return nil
	}

	return s.repository.Reduce(ctx, &req, token.UserID)

}

func (s *transactionService) GetBalance(ctx *gin.Context) (*GetBalanceResponse, error) {

	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return s.repository.Get(ctx, token.UserID)
}

func (s *transactionService) GetHistory(ctx *gin.Context, req GetHistoryPayload) (*GetHistoryResponse, int, error) {

	headerToken := ctx.GetHeader("Authorization")
	token, err := s.jwtService.GetPayload(headerToken)
	if err != nil {
		return nil, 0, ErrInvalidToken
	}

	return s.repository.GetHistory(ctx, &req, token.UserID)
}
