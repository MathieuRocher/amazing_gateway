package handler

import (
	"amazing_gateway/internal/adapter/application"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	useCase application.UserUsecaseInterface
}

func NewAuthHandler(uh application.UserUsecaseInterface) *AuthHandler {
	return &AuthHandler{useCase: uh}
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", h.Login)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	token, err := h.useCase.Authenticate(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
