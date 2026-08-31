package controller

import (
	"backendTemp/internal/dto"
	"backendTemp/internal/service"
	"errors"
	"time"

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
func (a *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := a.authService.Register(ctx.Request.Context(), req)
	if err != nil {
		slog.Error("registration failed", "error", err)

		if err == service.ErrEmailAlreadyTaken {
			ctx.JSON(http.StatusConflict, gin.H{"error": "email already taken"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": res,
	})
}

// Login godoc
// @Summary      Login a user
// @Description  Login with email and password, returns access token and sets refresh token cookie
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "User Login Credentials"
// @Success      200 {object} map[string]string "Login successful"
// @Failure      400 {object} map[string]string "Invalid request body"
// @Failure      401 {object} map[string]string "Invalid credentials"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /auth/login [post]
func (a *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	res, err := a.authService.Login(ctx.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		slog.Error("login failed", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.SetCookie(
		"refresh_token",
		res.RefreshToken,
		int(7*24*time.Hour.Seconds()),
		"/auth/refresh",
		"",
		true,
		true,
	)

	ctx.SetSameSite(http.SameSiteStrictMode)

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "login successful",
		"access_token": res.AccessToken,
	})
}
