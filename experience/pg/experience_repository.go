package pg

import (
	"context"
	"fmt"
	"localbe/experience"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	jobId := uuid.New()
	sqlFunc := insertExperienceEntry(jobId, companyName, position, start, end, description)
	queryInsert, argsInsert := sqlFunc()

	_, err := er.pgpool.Exec(ctx, queryInsert, argsInsert...)
	if err != nil {
		return nil, err
	}
	fmt.Println("Created")
	e, err := er.GetExperienceEntry(ctx, jobId)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (er *ExperienceRepository) GetExperienceEntry(ctx context.Context, jobId uuid.UUID) (*experience.Experience, error) {
	sqlFunc := selectExperienceEntry(jobId)
	querySelect, argsSelect := sqlFunc()

	rows, err := er.pgpool.Query(ctx, querySelect, argsSelect...)
	if err != nil {
		return nil, err
	}
	experienceEntry, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[experience.Experience])
	if err != nil {
		return nil, err
	}
	return &experienceEntry, nil
}

func (er *ExperienceRepository) GetExperience(ctx context.Context) ([]experience.Experience, error) {
	sqlFunc := selectAllExperience()
	querySelect, argsSelect := sqlFunc()

	rows, err := er.pgpool.Query(ctx, querySelect, argsSelect...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	experienceRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[experience.Experience])
	if err != nil {
		return nil, err
	}
	return experienceRows, nil
}
