package http

import (
	"task-api/internal/domain"
	"task-api/internal/middleware"
	"task-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authUC *usecase.AuthUseCase, userUC *usecase.UserUseCase, taskUC *usecase.TaskUseCase, tokens usecase.TokenService) *gin.Engine {
	r := gin.Default()

	authHandler := NewAuthHandler(authUC)
	userHandler := NewUserHandler(userUC)
	taskHandler := NewTaskHandler(taskUC)

	r.POST("/auth/login", authHandler.Login)

	authenticated := r.Group("/")
	authenticated.Use(middleware.Auth(tokens))
	{
		authenticated.POST("/auth/logout", authHandler.Logout)
		authenticated.PUT("/auth/password", authHandler.ChangePassword)
	}

	protected := r.Group("/")
	protected.Use(middleware.Auth(tokens), middleware.RequirePasswordChanged())
	{
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole(domain.RoleAdmin))
		{
			admin.POST("/users", userHandler.Create)
			admin.GET("/users", userHandler.List)
			admin.GET("/users/:id", userHandler.GetByID)
			admin.PUT("/users/:id", userHandler.Update)
			admin.DELETE("/users/:id", userHandler.Delete)

			admin.POST("/tasks", taskHandler.AdminCreate)
			admin.GET("/tasks", taskHandler.AdminList)
			admin.GET("/tasks/:id", taskHandler.AdminGetDetail)
			admin.PUT("/tasks/:id", taskHandler.AdminUpdate)
			admin.DELETE("/tasks/:id", taskHandler.AdminDelete)
		}

		executor := protected.Group("/")
		executor.Use(middleware.RequireRole(domain.RoleExecutor))
		{
			executor.GET("/tasks", taskHandler.ListMine)
			executor.GET("/tasks/:id", taskHandler.GetMine)
			executor.PUT("/tasks/:id/status", taskHandler.UpdateStatus)
			executor.POST("/tasks/:id/comments", taskHandler.AddComment)
		}

		audit := protected.Group("/audit")
		audit.Use(middleware.RequireRole(domain.RoleAuditor))
		{
			audit.GET("/tasks", taskHandler.AuditList)
			audit.GET("/tasks/:id", taskHandler.AuditGetDetail)
		}
	}

	return r
}
