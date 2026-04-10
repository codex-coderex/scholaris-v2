package colleges

import (
	"fmt"

	"scholaris-v2/internal/shared/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

// repository

type CollegeRepository struct {
	pool *pgxpool.Pool
}

func NewCollegeRepository(pool *pgxpool.Pool) *CollegeRepository {
	return &CollegeRepository{pool: pool}
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

func (r *CollegeRepository) GetAll(search, sortBy, order string, page, pageSize int) ([]College, int, error) {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	offset := (page - 1) * pageSize
	pattern := utils.SearchPattern(search)
	sortColumn := normalizeCollegeSort(sortBy)
	sortOrder := utils.NormalizeSortOrder(order)

	query := fmt.Sprintf(`
		SELECT code, name
		FROM   college
		WHERE  name ILIKE $1
		OR     code ILIKE $1
		ORDER  BY %s %s
		LIMIT  $2
		OFFSET $3
	`, sortColumn, sortOrder)

	rows, err := r.pool.Query(ctx, query, pattern, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("GetAll colleges: %w", err)
	}
	defer rows.Close()

	var colleges []College
	for rows.Next() {
		var c College
		if err := rows.Scan(&c.Code, &c.Name); err != nil {
			return nil, 0, fmt.Errorf("scan college: %w", err)
		}
		colleges = append(colleges, c)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate colleges: %w", err)
	}

	var total int
	if err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM college WHERE name ILIKE $1 OR code ILIKE $1`,
		pattern,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count colleges: %w", err)
	}

	return colleges, total, nil
}

// mutations

func (r *CollegeRepository) Create(c College) error {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	if _, err := r.pool.Exec(ctx, `
		INSERT INTO college (code, name)
		VALUES ($1, $2)
	`, c.Code, c.Name); err != nil {
		return fmt.Errorf("create college: %w", err)
	}
	return nil
}

func (r *CollegeRepository) Update(c College) error {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	tag, err := r.pool.Exec(ctx, `
		UPDATE college
		SET    name = $1
		WHERE  code = $2
	`, c.Name, c.Code)
	if err != nil {
		return fmt.Errorf("update college %s: %w", c.Code, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("update college %s: no matching record", c.Code)
	}

	return nil
}

func (r *CollegeRepository) Delete(code string) error {
	ctx, cancel := utils.NewDBContext()
	defer cancel()

	tag, err := r.pool.Exec(ctx, `
		DELETE FROM college
		WHERE  code = $1
	`, code)
	if err != nil {
		return fmt.Errorf("delete college %s: %w", code, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete college %s: no matching record", code)
	}

	return nil
}
