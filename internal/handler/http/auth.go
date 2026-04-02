package http

import (
	"net/http"
	"task-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *usecase.AuthUseCase
}

func NewAuthHandler(auth *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email y password son requeridos"})
		return
	}

	result, err := h.auth.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":         result.Tokens.AccessToken,
		"refresh_token":        result.Tokens.RefreshToken,
		"must_change_password": result.MustChangePassword,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token es requerido"})
		return
	}

	if err := h.auth.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sesión cerrada"})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current_password y new_password son requeridos"})
		return
	}

	userID := c.GetString("userID")
	if err := h.auth.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
		RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "contraseña actualizada"})
}
