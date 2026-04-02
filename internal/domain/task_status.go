package domain

type TaskStatus string

const (
	StatusAssigned  TaskStatus = "ASIGNADO"
	StatusStarted   TaskStatus = "INICIADO"
	StatusDoneOk    TaskStatus = "FINALIZADO_EXITO"
	StatusDoneError TaskStatus = "FINALIZADO_ERROR"
	StatusOnHold    TaskStatus = "EN_ESPERA"
)

// EN_ESPERA puede volver a INICIADO — el ejecutor puede retomar una tarea en espera.
// Los estados finalizados no tienen transiciones salientes.
var validTransitions = map[TaskStatus][]TaskStatus{
	StatusAssigned: {StatusStarted},
	StatusStarted:  {StatusDoneOk, StatusDoneError, StatusOnHold},
	StatusOnHold:   {StatusStarted},
}

func (s TaskStatus) CanTransitionTo(next TaskStatus) bool {
	allowed, ok := validTransitions[s]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == next {
			return true
		}
	}
	return false
}
