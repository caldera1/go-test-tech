package models

import "task-api/internal/domain"

// Los mappers traducen entre modelos de GORM y entidades de dominio.
// El dominio no conoce ni tags ni estructuras de persistencia.

func UserToModel(u domain.User) User {
	return User{
		ID:                 u.ID,
		Email:              u.Email,
		PasswordHash:       u.PasswordHash,
		Role:               string(u.Role),
		MustChangePassword: u.MustChangePassword,
		CreatedAt:          u.CreatedAt,
	}
}

func UserToDomain(m User) domain.User {
	return domain.User{
		ID:                 m.ID,
		Email:              m.Email,
		PasswordHash:       m.PasswordHash,
		Role:               domain.Role(m.Role),
		MustChangePassword: m.MustChangePassword,
		CreatedAt:          m.CreatedAt,
	}
}

func TaskToModel(t domain.Task) Task {
	return Task{
		ID:              t.ID,
		Title:           t.Title,
		Description:     t.Description,
		DueDate:         t.DueDate,
		Status:          string(t.Status),
		AssignedUserID:  t.AssignedUserID,
		CreatedByUserID: t.CreatedByUserID,
		CreatedAt:       t.CreatedAt,
	}
}

func TaskToDomain(m Task) domain.Task {
	return domain.Task{
		ID:              m.ID,
		Title:           m.Title,
		Description:     m.Description,
		DueDate:         m.DueDate,
		Status:          domain.TaskStatus(m.Status),
		AssignedUserID:  m.AssignedUserID,
		CreatedByUserID: m.CreatedByUserID,
		CreatedAt:       m.CreatedAt,
	}
}

func CommentToModel(c domain.Comment) Comment {
	return Comment{
		ID:        c.ID,
		TaskID:    c.TaskID,
		AuthorID:  c.AuthorID,
		Body:      c.Body,
		CreatedAt: c.CreatedAt,
	}
}

func CommentToDomain(m Comment) domain.Comment {
	return domain.Comment{
		ID:        m.ID,
		TaskID:    m.TaskID,
		AuthorID:  m.AuthorID,
		Body:      m.Body,
		CreatedAt: m.CreatedAt,
	}
}
