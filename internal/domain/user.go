package domain

import "time"

type User struct {
	ID                 string
	Email              string
	PasswordHash       string
	Role               Role
	MustChangePassword bool
	CreatedAt          time.Time
}
