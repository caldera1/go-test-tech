package http

import (
	"net/http"
	"task-api/internal/domain"
	"task-api/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	tasks *usecase.TaskUseCase
}

func NewTaskHandler(tasks *usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{tasks: tasks}
}

func (h *TaskHandler) AdminCreate(c *gin.Context) {
	var req struct {
		Title          string `json:"title" binding:"required"`
		Description    string `json:"description"`
		DueDate        string `json:"due_date" binding:"required"`
		AssignedUserID string `json:"assigned_user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "campos requeridos: title, due_date, assigned_user_id"})
		return
	}

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "due_date debe estar en formato RFC3339"})
		return
	}

	createdBy := c.GetString("userID")
	task, err := h.tasks.Create(c.Request.Context(), req.Title, req.Description, dueDate, req.AssignedUserID, createdBy)
	if err != nil {
		RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, taskResponse(task))
}

func (h *TaskHandler) AdminUpdate(c *gin.Context) {
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		DueDate     string `json:"due_date" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "campos requeridos: title, due_date"})
		return
	}

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "due_date debe estar en formato RFC3339"})
		return
	}

	if err := h.tasks.AdminUpdate(c.Request.Context(), c.Param("id"), req.Title, req.Description, dueDate); err != nil {
		RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tarea actualizada"})
}

func (h *TaskHandler) AdminDelete(c *gin.Context) {
	if err := h.tasks.AdminDelete(c.Request.Context(), c.Param("id")); err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tarea eliminada"})
}

func (h *TaskHandler) AdminGetDetail(c *gin.Context) {
	userID := c.GetString("userID")
	role, ok := c.Get("role")
	if !ok {
		RespondError(c, domain.ErrInvalidToken)
		return
	}
	r, ok := role.(domain.Role)
	if !ok {
		RespondError(c, domain.ErrInvalidToken)
		return
	}
	detail, err := h.tasks.GetDetail(c.Request.Context(), c.Param("id"), userID, r)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, taskDetailResponse(detail))
}

func (h *TaskHandler) AdminList(c *gin.Context) {
	tasks, err := h.tasks.ListAll(c.Request.Context())
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasksResponse(tasks))
}

func (h *TaskHandler) ListMine(c *gin.Context) {
	userID := c.GetString("userID")
	tasks, err := h.tasks.ListMine(c.Request.Context(), userID)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasksResponse(tasks))
}

func (h *TaskHandler) GetMine(c *gin.Context) {
	userID := c.GetString("userID")
	detail, err := h.tasks.GetMine(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, taskDetailResponse(detail))
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status es requerido"})
		return
	}

	userID := c.GetString("userID")
	if err := h.tasks.UpdateStatus(c.Request.Context(), c.Param("id"), userID, domain.TaskStatus(req.Status)); err != nil {
		RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "estado actualizado"})
}

func (h *TaskHandler) AddComment(c *gin.Context) {
	var req struct {
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "body es requerido"})
		return
	}

	userID := c.GetString("userID")
	if err := h.tasks.AddComment(c.Request.Context(), c.Param("id"), userID, req.Body); err != nil {
		RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "comentario agregado"})
}

func (h *TaskHandler) AuditList(c *gin.Context) {
	tasks, err := h.tasks.ListAll(c.Request.Context())
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasksResponse(tasks))
}

func (h *TaskHandler) AuditGetDetail(c *gin.Context) {
	userID := c.GetString("userID")
	role, ok := c.Get("role")
	if !ok {
		RespondError(c, domain.ErrInvalidToken)
		return
	}
	r, ok := role.(domain.Role)
	if !ok {
		RespondError(c, domain.ErrInvalidToken)
		return
	}
	detail, err := h.tasks.GetDetail(c.Request.Context(), c.Param("id"), userID, r)
	if err != nil {
		RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, taskDetailResponse(detail))
}

func taskResponse(t domain.Task) gin.H {
	return gin.H{
		"id":               t.ID,
		"title":            t.Title,
		"description":      t.Description,
		"due_date":         t.DueDate.Format(time.RFC3339),
		"status":           t.Status,
		"assigned_user_id": t.AssignedUserID,
		"created_by":       t.CreatedByUserID,
		"created_at":       t.CreatedAt.Format(time.RFC3339),
	}
}

func taskDetailResponse(d usecase.TaskDetail) gin.H {
	comments := make([]gin.H, len(d.Comments))
	for i, c := range d.Comments {
		comments[i] = gin.H{
			"id":         c.ID,
			"author_id":  c.AuthorID,
			"body":       c.Body,
			"created_at": c.CreatedAt.Format(time.RFC3339),
		}
	}
	resp := taskResponse(d.Task)
	resp["comments"] = comments
	return resp
}

func tasksResponse(tasks []domain.Task) []gin.H {
	resp := make([]gin.H, len(tasks))
	for i, t := range tasks {
		resp[i] = taskResponse(t)
	}
	return resp
}
