package user

import (
	"context"
	"fmt"

	"github.com/gabriel-tama/banking-app/common/jwt"
	"github.com/gabriel-tama/banking-app/common/password"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Create(ctx context.Context, req RegisterPayload) (*AuthResponse, error)
	FindByCredential(ctx context.Context, req LoginPayload) (*AuthResponse, error)
	PingDB(ctx *gin.Context) error
}

type userService struct {
	repository Repository
	jwtService jwt.JWTService
}

func NewService(repository Repository, jwtService jwt.JWTService) Service {
	return &userService{repository: repository, jwtService: jwtService}
}

func (s *userService) Create(ctx context.Context, req RegisterPayload) (*AuthResponse, error) {

	hashedPassword, err := password.Hash(req.Password, s.repository.GetSalt())
	if err != nil {
		return nil, err
	}
	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	err = s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	accessToken, err := s.jwtService.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Name:        req.Name,
		Email:       req.Email,
		AccessToken: accessToken,
	}, nil
}

func (s *userService) FindByCredential(ctx context.Context, req LoginPayload) (*AuthResponse, error) {
	user := &User{
		Email: req.Email,
	}

	err := s.repository.FindByCredential(ctx, user)
	if err != nil {
		return nil, err
	}
	match, err := password.Matches(req.Password, user.Password)

	if err != nil {
		return nil, err
	}

	if !match {
		return nil, fmt.Errorf("%w:%w", ErrWrongPassword, err)
	}
	accessToken, err := s.jwtService.CreateToken(int(user.ID))
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: accessToken,
	}, nil

}

func (s *userService) PingDB(ctx *gin.Context) error {
	return s.repository.PingDB(ctx)
}
