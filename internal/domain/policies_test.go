package domain

import (
	"testing"
	"time"
)

func TestCanExecutorUpdateTask(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)
	executorID := "exec-01"

	cases := []struct {
		name   string
		task   Task
		userID string
		want   bool
	}{
		{"ejecutor asignado, tarea vigente", Task{AssignedUserID: executorID, DueDate: future}, executorID, true},
		{"otro ejecutor no puede modificar", Task{AssignedUserID: "otro-exec", DueDate: future}, executorID, false},
		{"tarea vencida bloquea actualizacion", Task{AssignedUserID: executorID, DueDate: past}, executorID, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CanExecutorUpdateTask(tc.task, tc.userID, now)
			if got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestCanAddComment(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)
	userID := "exec-01"

	cases := []struct {
		name   string
		task   Task
		userID string
		want   bool
	}{
		{"ejecutor comenta tarea vencida propia", Task{AssignedUserID: userID, DueDate: past}, userID, true},
		{"no puede comentar tarea vigente", Task{AssignedUserID: userID, DueDate: future}, userID, false},
		{"no puede comentar tarea ajena aunque este vencida", Task{AssignedUserID: "otro-exec", DueDate: past}, userID, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CanAddComment(tc.task, tc.userID, now)
			if got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}
