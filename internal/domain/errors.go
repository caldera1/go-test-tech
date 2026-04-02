package domain

import "errors"

var (
	ErrNotFound               = errors.New("recurso no encontrado")
	ErrInvalidCredentials     = errors.New("credenciales inválidas")
	ErrInvalidToken           = errors.New("token inválido o expirado")
	ErrForbidden              = errors.New("acción no permitida")
	ErrInvalidTaskTransition  = errors.New("transición de estado inválida")
	ErrTaskExpired            = errors.New("la tarea está vencida")
	ErrAdminCannotBeCreated   = errors.New("no se puede crear un usuario administrador")
	ErrAssigneeNotExecutor    = errors.New("el usuario asignado debe tener perfil ejecutor")
	ErrPasswordChangeRequired = errors.New("debe cambiar su contraseña antes de continuar")
	ErrUserHasAssignedTasks   = errors.New("no se puede eliminar un usuario con tareas asignadas")
)
