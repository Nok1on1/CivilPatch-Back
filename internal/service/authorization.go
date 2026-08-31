package service

import (
	"backendTemp/internal/dto"
	"backendTemp/internal/repository"
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
}

type DefaultauthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &DefaultauthService{userRepo: repo}
}

func (s *DefaultauthService) Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.RegisterResponse{}, err
	}

	user := req.ToModel(string(hashedPassword))

	if err = s.userRepo.Create(ctx, &user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return dto.RegisterResponse{}, ErrEmailAlreadyTaken
		}
		return dto.RegisterResponse{}, err
	}

	return dto.NewRegisterResponse(&user), nil
}
