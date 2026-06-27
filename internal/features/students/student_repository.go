package students

import (
	"fmt"
	"regexp"

	"scholaris-v2/internal/shared/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// repository

type StudentRepository struct {
	pool *pgxpool.Pool
}

func NewStudentRepository(pool *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{pool: pool}
}

func normalizeStudentSort(sortBy string) string {
	switch sortBy {
	case "id", "s.id":
		return "s.id"
	case "first_name", "s.first_name":
		return "s.first_name"
	case "last_name", "s.last_name":
		return "s.last_name"
	case "year", "s.year":
		return "s.year"
	case "gender", "s.gender":
		return "s.gender"
	case "program", "p.code":
		// Sort by what the UI displays:
		// - real program code when the program still exists
		// - N/A when the student's saved program_code points to a deleted program
		return "CASE WHEN p.code IS NULL THEN 'N/A' ELSE p.code END"
	case "college", "c.code", "c.name", "p.college_code":
		// Sort by what the UI displays:
		// - real college code when both program and college still exist
		// - N/A when the program is deleted or the college is deleted
		return "CASE WHEN p.code IS NULL OR c.code IS NULL THEN 'N/A' ELSE c.code END"
	default:
		return "s.id"
	}
}

// queries

func (r *StudentRepository) GetAll(search, sortBy, order string, page, pageSize int) ([]Student, int, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	offset := (page - 1) * pageSize
	pattern := utils.SearchPattern(search)
	sortColumn := normalizeStudentSort(sortBy)
	sortOrder := utils.NormalizeSortOrder(order)

	query := fmt.Sprintf(`
		SELECT s.id, s.first_name, s.last_name,
		       s.year, s.gender, COALESCE(s.program_code, ''),
		       COALESCE(p.code, ''), COALESCE(p.name, ''), COALESCE(p.college_code, ''),
		       COALESCE(c.code, ''), COALESCE(c.name, '')
		FROM   student s
		LEFT JOIN program p ON s.program_code = p.code
		LEFT JOIN college c ON p.college_code = c.code
		WHERE  s.id         ILIKE $1
		OR     s.first_name ILIKE $1
		OR     s.last_name  ILIKE $1
		OR     p.code       ILIKE $1
		OR     p.name       ILIKE $1
		OR     c.code       ILIKE $1
		OR     c.name       ILIKE $1
		ORDER  BY %s %s, s.id ASC
		LIMIT  $2
		OFFSET $3
	`, sortColumn, sortOrder)

	rows, err := r.pool.Query(ctx, query, pattern, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("GetAll students: %w", err)
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(
			&s.Id, &s.FirstName, &s.LastName,
			&s.Year, &s.Gender, &s.ProgramCode,
			&s.Program.Code, &s.Program.Name, &s.Program.CollegeCode,
			&s.Program.College.Code, &s.Program.College.Name,
		); err != nil {
			return nil, 0, fmt.Errorf("scan student: %w", err)
		}
		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate students: %w", err)
	}

	var total int
	if err = r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM   student s
		LEFT JOIN program p ON s.program_code = p.code
		LEFT JOIN college c ON p.college_code = c.code
		WHERE  s.id         ILIKE $1
		OR     s.first_name ILIKE $1
		OR     s.last_name  ILIKE $1
		OR     p.code       ILIKE $1
		OR     p.name       ILIKE $1
		OR     c.code       ILIKE $1
		OR     c.name       ILIKE $1
	`, pattern).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count students: %w", err)
	}

	return students, total, nil
}

// mutations
var studentIDPattern = regexp.MustCompile(`^\d{4}-\d{4}$`)

func (r *StudentRepository) Create(s Student) error {
	if !studentIDPattern.MatchString(s.Id) {
		return fmt.Errorf("invalid student ID format: must be YYYY-NNNN (e.g. 2024-0001)")
	}

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	if _, err := r.pool.Exec(ctx, `
		INSERT INTO student (id, first_name, last_name, year, gender, program_code)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, s.Id, s.FirstName, s.LastName, s.Year, s.Gender, s.ProgramCode); err != nil {
		return fmt.Errorf("create student: %w", err)
	}
	return nil
}

func (r *StudentRepository) Update(s Student) error {
	if !studentIDPattern.MatchString(s.Id) {
		return fmt.Errorf("invalid student ID format: must be YYYY-NNNN (e.g. 2024-0001)")
	}

	ctx, cancel := utils.NewDBContext()
	defer cancel()

	tag, err := r.pool.Exec(ctx, `
        UPDATE student
        SET    id           = $1,
               first_name   = $2,
               last_name    = $3,
               year         = $4,
               gender       = $5,
               program_code = $6
        WHERE  id           = $7
    `, s.Id, s.FirstName, s.LastName, s.Year, s.Gender, s.ProgramCode, s.OriginalId)
	if err != nil {
		return fmt.Errorf("update student %s: %w", s.OriginalId, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("update student %s: no matching record", s.OriginalId)
	}

	return nil
}

func (r *StudentRepository) Delete(id string) error {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	tag, err := r.pool.Exec(ctx, `
		DELETE FROM student
		WHERE  id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("delete student %s: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete student %s: no matching record", id)
	}

	return nil
}
