package dto

import "backendTemp/internal/model"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (r *RegisterRequest) ToModel(hashedPassword string) model.User {
	u := model.User{
		Email:        r.Email,
		PasswordHash: hashedPassword,
	}
	u.SetTimestamps()
	return u
}

type RegisterResponse struct {
	Email string `json:"email"`
}

func NewRegisterResponse(u *model.User) RegisterResponse {
	return RegisterResponse{
		Email: u.Email,
	}
}
