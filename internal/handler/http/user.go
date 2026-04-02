package http

import (
	"net/http"
	"task-api/internal/domain"
	"task-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	users *usecase.UserUseCase
}

func NewUserHandler(users *usecase.UserUseCase) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
		Role  string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email y role son requeridos"})
		return
	}

	role := domain.Role(req.Role)
	if !role.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rol inválido, valores permitidos: executor, auditor"})
		return
	}

	result, err := h.users.Create(c.Request.Context(), req.Email, role)
	if err != nil {
		RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":                 result.User.ID,
		"email":              result.User.Email,
		"role":               result.User.Role,
		"temporary_password": result.TemporaryPassword,
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	user, err := h.users.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, userResponse(user))
}

func (h *UserHandler) Update(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
		Role  string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email y role son requeridos"})
		return
	}

	role := domain.Role(req.Role)
	if !role.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rol inválido, valores permitidos: executor, auditor"})
		return
	}

	user, err := h.users.Update(c.Request.Context(), c.Param("id"), req.Email, role)
	if err != nil {
		RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, userResponse(user))
}

func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.users.Delete(c.Request.Context(), c.Param("id")); err != nil {
		RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "usuario eliminado"})
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.users.List(c.Request.Context())
	if err != nil {
		RespondError(c, err)
		return
	}

	resp := make([]gin.H, len(users))
	for i, u := range users {
		resp[i] = userResponse(u)
	}

	c.JSON(http.StatusOK, resp)
}

func userResponse(u domain.User) gin.H {
	return gin.H{
		"id":    u.ID,
		"email": u.Email,
		"role":  u.Role,
	}
}
