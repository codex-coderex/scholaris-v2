package utils

func NormalizeSortOrder(order string) string {
	if order == "DESC" {
		return "DESC"
	}
	return "ASC"
}

func SearchPattern(s string) string {
	return "%" + s + "%"
}
