package programs

import (
	"context"
	"fmt"

	"scholaris-v2/internal/shared/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// repository

type ProgramRepository struct {
	pool *pgxpool.Pool
}

func NewProgramRepository(pool *pgxpool.Pool) *ProgramRepository {
	return &ProgramRepository{pool: pool}
}

// helpers

func (r *ProgramRepository) ctx() context.Context {
	return context.Background()
}

func normalizeProgramSort(sortBy string) string {
	switch sortBy {
	case "code", "p.code":
		return "p.code"
	case "name", "p.name":
		return "p.name"
	case "college", "c.name":
		return "c.name"
	default:
		return "p.code"
	}
}

// queries

func (r *ProgramRepository) GetAll(search, sortBy, order string, page, pageSize int, collegeCode string) ([]Program, int, error) {
	offset := (page - 1) * pageSize
	pattern := utils.SearchPattern(search)
	sortColumn := normalizeProgramSort(sortBy)
	sortOrder := utils.NormalizeSortOrder(order)

	query := fmt.Sprintf(`
		SELECT p.code, p.name, COALESCE(p.college_code, ''),
		       COALESCE(c.code, ''), COALESCE(c.name, '')
		FROM   program p
		LEFT JOIN college c ON p.college_code = c.code
		WHERE  (p.name ILIKE $1
		OR     p.code ILIKE $1
		OR     c.name ILIKE $1)
		AND    ($2 = '' OR p.college_code = $2)
		ORDER  BY %s %s
		LIMIT  $3
		OFFSET $4
	`, sortColumn, sortOrder)

	rows, err := r.pool.Query(r.ctx(), query, pattern, collegeCode, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("GetAll programs: %w", err)
	}
	defer rows.Close()

	var programs []Program
	for rows.Next() {
		var p Program
		if err := rows.Scan(
			&p.Code, &p.Name, &p.CollegeCode,
			&p.College.Code, &p.College.Name,
		); err != nil {
			return nil, 0, fmt.Errorf("scan program: %w", err)
		}
		programs = append(programs, p)
	}

	var total int
	if err = r.pool.QueryRow(r.ctx(), `
		SELECT COUNT(*)
		FROM   program p
		LEFT JOIN college c ON p.college_code = c.code
		WHERE  (p.name ILIKE $1
		OR     p.code ILIKE $1
		OR     c.name ILIKE $1)
		AND    ($2 = '' OR p.college_code = $2)
	`, pattern, collegeCode).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count programs: %w", err)
	}

	return programs, total, nil
}

func (r *ProgramRepository) GetByCode(code string) (Program, error) {
	var p Program
	if err := r.pool.QueryRow(r.ctx(), `
		SELECT p.code, p.name, COALESCE(p.college_code, ''),
		       COALESCE(c.code, ''), COALESCE(c.name, '')
		FROM   program p
		LEFT JOIN college c ON p.college_code = c.code
		WHERE  p.code = $1
	`, code).Scan(
		&p.Code, &p.Name, &p.CollegeCode,
		&p.College.Code, &p.College.Name,
	); err != nil {
		return Program{}, fmt.Errorf("GetByCode %s: %w", code, err)
	}
	return p, nil
}

// mutations

func (r *ProgramRepository) Create(p Program) error {
	if _, err := r.pool.Exec(r.ctx(), `
		INSERT INTO program (code, name, college_code)
		VALUES ($1, $2, $3)
	`, p.Code, p.Name, p.CollegeCode); err != nil {
		return fmt.Errorf("create program: %w", err)
	}
	return nil
}

func (r *ProgramRepository) Update(p Program) error {
	if _, err := r.pool.Exec(r.ctx(), `
		UPDATE program
		SET    name         = $1,
		       college_code = $2
		WHERE  code         = $3
	`, p.Name, p.CollegeCode, p.Code); err != nil {
		return fmt.Errorf("update program %s: %w", p.Code, err)
	}
	return nil
}

func (r *ProgramRepository) Delete(code string) error {
	if _, err := r.pool.Exec(r.ctx(), `
		DELETE FROM program
		WHERE  code = $1
	`, code); err != nil {
		return fmt.Errorf("delete program %s: %w", code, err)
	}
	return nil
}
