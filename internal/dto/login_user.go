package dto

import "backendTemp/internal/model"

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (l *LoginRequest) ToModel(passwordHash string) model.User {
	return model.User{
		Email:       l.Email,
		PasswordHash: passwordHash,
	}
}

type LoginResponse struct {
	AccessToken string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
