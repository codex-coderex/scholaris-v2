package utils

import (
	"context"
	"time"
)

const DBQueryTimeout = 5 * time.Second

func NewDBContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DBQueryTimeout)
}

func NormalizeSortOrder(order string) string {
	if order == "DESC" {
		return "DESC"
	}
	return "ASC"
}

func SearchPattern(s string) string {
	return "%" + s + "%"
}
