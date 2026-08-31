package controller

import (
	"backendTemp/internal/dto"
	"backendTemp/internal/service"

	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "User Registration Credentials"
// @Success      201 {object} dto.RegisterResponse
// @Failure      400 {object} map[string]string "Invalid request body"
// @Failure      409 {object} map[string]string "Email already taken"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /auth/register [post]
func (a *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := a.authService.Register(c.Request.Context(), req)
	if err != nil {
		slog.Error("registration failed", "error", err)

		if err == service.ErrEmailAlreadyTaken {
			c.JSON(http.StatusConflict, gin.H{"error": "email already taken"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": res,
	})
}
