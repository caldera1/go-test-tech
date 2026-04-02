package domain

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleExecutor Role = "executor"
	RoleAuditor  Role = "auditor"
)

func (r Role) IsValid() bool {
	return r == RoleAdmin || r == RoleExecutor || r == RoleAuditor
}
