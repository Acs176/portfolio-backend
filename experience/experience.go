package experience

import (
	"context"

	"github.com/google/uuid"
)

type Experience struct {
	Id              uuid.UUID
	CompanyName     string
	Position        string
	PeriodStart     string
	PeriodEnd       string
	RoleDescription string
}

type CreateExperienceEntryFunc func(ctx context.Context, companyName, position, start, end, description string) (*Experience, error)
