package experience

import "github.com/google/uuid"

type Experience struct {
	Id              uuid.UUID
	CompanyName     string
	Position        string
	PeriodStart     string
	PeriodEnd       string
	RoleDescription string
}
