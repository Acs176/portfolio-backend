package pg

import (
	"context"
	"localbe/experience"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExperienceRepository struct {
	pgpool *pgxpool.Pool
}

func NewExperienceRepository(pgpool *pgxpool.Pool) *ExperienceRepository {
	return &ExperienceRepository{
		pgpool: pgpool,
	}
}

func (er *ExperienceRepository) CreateExperienceEntry(
	ctx context.Context,
	companyName,
	position,
	start,
	end,
	description string,
) (*experience.Experience, error) {
	jobId := uuid.New().String()
	sqlFunc := insertExperienceEntry(jobId, companyName, position, start, end, description)
	queryInsert, argsInsert := sqlFunc()

	_, err := er.pgpool.Exec(ctx, queryInsert, argsInsert...)
	if err != nil {
		return nil, err
	}

	e, err := er.GetExperienceEntry(ctx, jobId)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (er *ExperienceRepository) GetExperienceEntry(ctx context.Context, jobId string) (*experience.Experience, error) {
	sqlFunc := selectExperienceEntry(jobId)
	querySelect, argsSelect := sqlFunc()

	var experienceEntry experience.Experience
	row := er.pgpool.QueryRow(ctx, querySelect, argsSelect)
	err := row.Scan(&experienceEntry)
	if err != nil {
		return nil, err
	}
	return &experienceEntry, nil
}
