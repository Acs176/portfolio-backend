package pg

type SQLFunc func() (string, []any)

func insertExperienceEntry(id, companyName, position, start, end, description string) SQLFunc {
	query := `
	INSERT INTO testdb_schema.experience
	(id, company_name, position, period_start, period_end, role_description)
	VALUES ($1, $2, $3, $4, $5)
	`
	args := []any{id, companyName, position, start, end, description}
	return func() (string, []any) { return query, args }
}

func selectExperienceEntry(id string) SQLFunc {
	query := `
		SELECT (id, company_name, position, period_start, period_end, role_description)
		FROM testdb_schema.experience 
		WHERE id = $1 
	`
	args := []any{id}
	return func() (string, []any) { return query, args }
}
