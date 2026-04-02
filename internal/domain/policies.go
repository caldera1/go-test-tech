package domain

import "time"

// CanExecutorUpdateTask valida ownership y que la tarea no esté vencida.
// No valida la transición en sí — eso es responsabilidad de task.TransitionTo.
func CanExecutorUpdateTask(task Task, userID string, now time.Time) bool {
	return task.AssignedUserID == userID && !task.IsExpired(now)
}

// Los comentarios solo están permitidos en tareas vencidas.
// Esto da al ejecutor un canal para documentar por qué no pudo completarla a tiempo.
func CanAddComment(task Task, userID string, now time.Time) bool {
	return task.AssignedUserID == userID && task.IsExpired(now)
}

func CanViewTask(task Task, userID string, role Role) bool {
	return role == RoleAuditor || role == RoleAdmin || task.AssignedUserID == userID
}
