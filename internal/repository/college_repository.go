package repository

import (
	"context"
	"fmt"

	"scholaris-v2/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// repository

type CollegeRepository struct {
	pool *pgxpool.Pool
}

func NewCollegeRepository(pool *pgxpool.Pool) *CollegeRepository {
	return &CollegeRepository{pool: pool}
}

// helpers

func (r *CollegeRepository) ctx() context.Context {
	return context.Background()
}

func searchPattern(s string) string {
	return "%" + s + "%"
}

func normalizeCollegeSort(sortBy string) string {
	switch sortBy {
	case "code":
		return "code"
	case "name":
		return "name"
	default:
		return "code"
	}
}

// queries

func (r *CollegeRepository) GetAll(search, sortBy, order string, page, pageSize int) ([]models.College, int, error) {
	offset := (page - 1) * pageSize
	pattern := searchPattern(search)
	sortColumn := normalizeCollegeSort(sortBy)
	sortOrder := normalizeSortOrder(order)

	query := fmt.Sprintf(`
		SELECT code, name
		FROM   college
		WHERE  name ILIKE $1
		OR     code ILIKE $1
		ORDER  BY %s %s
		LIMIT  $2
		OFFSET $3
	`, sortColumn, sortOrder)

	rows, err := r.pool.Query(r.ctx(), query, pattern, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("GetAll colleges: %w", err)
	}
	defer rows.Close()

	var colleges []models.College
	for rows.Next() {
		var c models.College
		if err := rows.Scan(&c.Code, &c.Name); err != nil {
			return nil, 0, fmt.Errorf("scan college: %w", err)
		}
		colleges = append(colleges, c)
	}

	var total int
	if err = r.pool.QueryRow(r.ctx(),
		`SELECT COUNT(*) FROM college WHERE name ILIKE $1 OR code ILIKE $1`,
		pattern,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count colleges: %w", err)
	}

	return colleges, total, nil
}

func (r *CollegeRepository) GetByCode(code string) (models.College, error) {
	var c models.College
	if err := r.pool.QueryRow(r.ctx(), `
		SELECT code, name
		FROM   college
		WHERE  code = $1
	`, code).Scan(&c.Code, &c.Name); err != nil {
		return models.College{}, fmt.Errorf("GetByCode %s: %w", code, err)
	}
	return c, nil
}

// mutations

func (r *CollegeRepository) Create(c models.College) error {
	if _, err := r.pool.Exec(r.ctx(), `
		INSERT INTO college (code, name)
		VALUES ($1, $2)
	`, c.Code, c.Name); err != nil {
		return fmt.Errorf("create college: %w", err)
	}
	return nil
}

func (r *CollegeRepository) Update(c models.College) error {
	if _, err := r.pool.Exec(r.ctx(), `
		UPDATE college
		SET    name = $1
		WHERE  code = $2
	`, c.Name, c.Code); err != nil {
		return fmt.Errorf("update college %s: %w", c.Code, err)
	}
	return nil
}

func (r *CollegeRepository) Delete(code string) error {
	if _, err := r.pool.Exec(r.ctx(), `
		DELETE FROM college
		WHERE  code = $1
	`, code); err != nil {
		return fmt.Errorf("delete college %s: %w", code, err)
	}
	return nil
}
