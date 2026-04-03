package repository

func normalizeSortOrder(order string) string {
	if order == "DESC" {
		return "DESC"
	}
	return "ASC"
}