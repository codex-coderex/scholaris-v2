package programs

import (
	"fmt"
	"strings"

	"scholaris-v2/internal/shared/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// repository

type ProgramRepository struct {
	pool *pgxpool.Pool
}

const noCollegeFilter = "__NO_COLLEGE__"

func NewProgramRepository(pool *pgxpool.Pool) *ProgramRepository {
	return &ProgramRepository{pool: pool}
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

func nullableCollegeCode(collegeCode string) any {
	trimmed := strings.TrimSpace(collegeCode)
	if trimmed == "" {
		return nil
	}

	return trimmed
}

// queries

func (r *ProgramRepository) GetAll(search, sortBy, order string, page, pageSize int, collegeCode string) ([]Program, int, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

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
		AND    (
			$2 = ''
			OR ($2 = $5 AND p.college_code IS NULL)
			OR ($2 <> $5 AND p.college_code = $2)
		)
		ORDER  BY %s %s
		LIMIT  $3
		OFFSET $4
	`, sortColumn, sortOrder)

	rows, err := r.pool.Query(ctx, query, pattern, collegeCode, pageSize, offset, noCollegeFilter)
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

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate programs: %w", err)
	}

	var total int
	if err = r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM   program p
		LEFT JOIN college c ON p.college_code = c.code
		WHERE  (p.name ILIKE $1
		OR     p.code ILIKE $1
		OR     c.name ILIKE $1)
		AND    (
			$2 = ''
			OR ($2 = $3 AND p.college_code IS NULL)
			OR ($2 <> $3 AND p.college_code = $2)
		)
	`, pattern, collegeCode, noCollegeFilter).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count programs: %w", err)
	}

	return programs, total, nil
}

// mutations

func (r *ProgramRepository) Create(p Program) error {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collegeCode := nullableCollegeCode(p.CollegeCode)

	if _, err := r.pool.Exec(ctx, `
		INSERT INTO program (code, name, college_code)
		VALUES ($1, $2, $3)
	`, p.Code, p.Name, collegeCode); err != nil {
		return fmt.Errorf("create program: %w", err)
	}
	return nil
}

func (r *ProgramRepository) Update(p Program) error {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	collegeCode := nullableCollegeCode(p.CollegeCode)

	tag, err := r.pool.Exec(ctx, `
		UPDATE program
		SET    name         = $1,
		       college_code = $2
		WHERE  code         = $3
	`, p.Name, collegeCode, p.Code)
	if err != nil {
		return fmt.Errorf("update program %s: %w", p.Code, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("update program %s: no matching record", p.Code)
	}

	return nil
}

func (r *ProgramRepository) Delete(code string) error {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	tag, err := r.pool.Exec(ctx, `
		DELETE FROM program
		WHERE  code = $1
	`, code)
	if err != nil {
		return fmt.Errorf("delete program %s: %w", code, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete program %s: no matching record", code)
	}

	return nil
}
