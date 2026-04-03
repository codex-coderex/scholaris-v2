package repository

import (
	"context"
	"fmt"

	"scholaris-v2/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// repository

type StudentRepository struct {
	pool *pgxpool.Pool
}

func NewStudentRepository(pool *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{pool: pool}
}

// helpers

func (r *StudentRepository) ctx() context.Context {
	return context.Background()
}

func normalizeStudentSort(sortBy string) string {
	switch sortBy {
	case "id", "s.id":
		return "s.id"
	case "first_name", "s.first_name":
		return "s.first_name"
	case "last_name", "s.last_name":
		return "s.last_name"
	case "program", "p.code":
		return "p.code"
	case "college", "c.name":
		return "c.name"
	default:
		return "s.id"
	}
}

// queries

func (r *StudentRepository) GetAll(search, sortBy, order string, page, pageSize int) ([]models.Student, int, error) {
	offset := (page - 1) * pageSize
	pattern := searchPattern(search)
	sortColumn := normalizeStudentSort(sortBy)
	sortOrder := normalizeSortOrder(order)

	query := fmt.Sprintf(`
		SELECT s.id, s.first_name, s.last_name,
		       s.year, s.gender, s.program_code,
		       p.code, p.name,
		       c.code, c.name
		FROM   student s
		JOIN   program p ON s.program_code = p.code
		JOIN   college c ON p.college_code = c.code
		WHERE  s.id         ILIKE $1
		OR     s.first_name ILIKE $1
		OR     s.last_name  ILIKE $1
		OR     p.code       ILIKE $1
		OR     p.name       ILIKE $1
		ORDER  BY %s %s
		LIMIT  $2
		OFFSET $3
	`, sortColumn, sortOrder)

	rows, err := r.pool.Query(r.ctx(), query, pattern, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("GetAll students: %w", err)
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(
			&s.Id, &s.FirstName, &s.LastName,
			&s.Year, &s.Gender, &s.ProgramCode,
			&s.Program.Code, &s.Program.Name,
			&s.Program.College.Code, &s.Program.College.Name,
		); err != nil {
			return nil, 0, fmt.Errorf("scan student: %w", err)
		}
		students = append(students, s)
	}

	var total int
	if err = r.pool.QueryRow(r.ctx(), `
		SELECT COUNT(*)
		FROM   student s
		JOIN   program p ON s.program_code = p.code
		JOIN   college c ON p.college_code = c.code
		WHERE  s.id         ILIKE $1
		OR     s.first_name ILIKE $1
		OR     s.last_name  ILIKE $1
		OR     p.code       ILIKE $1
		OR     p.name       ILIKE $1
	`, pattern).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count students: %w", err)
	}

	return students, total, nil
}

func (r *StudentRepository) GetById(id string) (models.Student, error) {
	var s models.Student
	if err := r.pool.QueryRow(r.ctx(), `
		SELECT s.id, s.first_name, s.last_name,
		       s.year, s.gender, s.program_code,
		       p.code, p.name,
		       c.code, c.name
		FROM   student s
		JOIN   program p ON s.program_code = p.code
		JOIN   college c ON p.college_code = c.code
		WHERE  s.id = $1
	`, id).Scan(
		&s.Id, &s.FirstName, &s.LastName,
		&s.Year, &s.Gender, &s.ProgramCode,
		&s.Program.Code, &s.Program.Name,
		&s.Program.College.Code, &s.Program.College.Name,
	); err != nil {
		return models.Student{}, fmt.Errorf("GetById %s: %w", id, err)
	}
	return s, nil
}

// mutations

func (r *StudentRepository) Create(s models.Student) error {
	if _, err := r.pool.Exec(r.ctx(), `
		INSERT INTO student (id, first_name, last_name, year, gender, program_code)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, s.Id, s.FirstName, s.LastName, s.Year, s.Gender, s.ProgramCode); err != nil {
		return fmt.Errorf("create student: %w", err)
	}
	return nil
}

func (r *StudentRepository) Update(s models.Student) error {
	if _, err := r.pool.Exec(r.ctx(), `
		UPDATE student
		SET    first_name   = $1,
		       last_name    = $2,
		       year         = $3,
		       gender       = $4,
		       program_code = $5
		WHERE  id           = $6
	`, s.FirstName, s.LastName, s.Year, s.Gender, s.ProgramCode, s.Id); err != nil {
		return fmt.Errorf("update student %s: %w", s.Id, err)
	}
	return nil
}

func (r *StudentRepository) Delete(id string) error {
	if _, err := r.pool.Exec(r.ctx(), `
		DELETE FROM student
		WHERE  id = $1
	`, id); err != nil {
		return fmt.Errorf("delete student %s: %w", id, err)
	}
	return nil
}
