package domain

import "time"

type Task struct {
	ID              string
	Title           string
	Description     string
	DueDate         time.Time
	Status          TaskStatus
	AssignedUserID  string
	CreatedByUserID string
	CreatedAt       time.Time
}

func (t *Task) IsExpired(now time.Time) bool {
	return now.After(t.DueDate)
}

func (t *Task) TransitionTo(next TaskStatus, now time.Time) error {
	// Verificamos vencimiento antes que la transición: una tarea vencida
	// no puede cambiar de estado independientemente de la transición solicitada.
	if t.IsExpired(now) {
		return ErrTaskExpired
	}
	if !t.Status.CanTransitionTo(next) {
		return ErrInvalidTaskTransition
	}
	t.Status = next
	return nil
}

// Solo ASIGNADO es modificable por admin — una vez iniciada, la tarea
// pertenece al ejecutor y el admin no puede interferir.
func (t Task) CanBeAdminModified() bool {
	return t.Status == StatusAssigned
}
