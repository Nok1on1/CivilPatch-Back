package service

import (
	"backendTemp/internal/dto"
	"backendTemp/internal/repository"
	"context"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
}

type DefaultAuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &DefaultAuthService{userRepo: repo}
}

func (s *DefaultAuthService) Register(ctx context.Context, req dto.RegisterRequest) (dto.RegisterResponse, error) {
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

func (s *DefaultAuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(2 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		AccessToken:  tokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
