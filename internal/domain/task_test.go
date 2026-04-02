package domain

import (
	"errors"
	"testing"
	"time"
)

func TestTaskTransitionTo(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)

	cases := []struct {
		name    string
		from    TaskStatus
		to      TaskStatus
		dueDate time.Time
		wantErr error
	}{
		{"ejecutor inicia tarea asignada", StatusAssigned, StatusStarted, future, nil},
		{"ejecutor finaliza con exito", StatusStarted, StatusDoneOk, future, nil},
		{"ejecutor finaliza con error", StatusStarted, StatusDoneError, future, nil},
		{"ejecutor pone en espera", StatusStarted, StatusOnHold, future, nil},
		{"ejecutor retoma tarea en espera", StatusOnHold, StatusStarted, future, nil},
		{"no puede saltar de asignado a finalizado", StatusAssigned, StatusDoneOk, future, ErrInvalidTaskTransition},
		{"tarea vencida bloquea cualquier transicion", StatusAssigned, StatusStarted, past, ErrTaskExpired},
		{"estado final no tiene transiciones salientes", StatusDoneOk, StatusStarted, future, ErrInvalidTaskTransition},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			task := &Task{Status: tc.from, DueDate: tc.dueDate}
			err := task.TransitionTo(tc.to, now)

			if tc.wantErr == nil && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected %v, got %v", tc.wantErr, err)
			}
		})
	}
}

func TestIsExpired(t *testing.T) {
	now := time.Now()

	task := &Task{DueDate: now.Add(-1 * time.Hour)}
	if !task.IsExpired(now) {
		t.Fatal("expected task to be expired")
	}

	task = &Task{DueDate: now.Add(1 * time.Hour)}
	if task.IsExpired(now) {
		t.Fatal("expected task not to be expired")
	}
}
