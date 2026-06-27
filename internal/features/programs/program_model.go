package programs

// program

type Program struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	CollegeCode  string  `json:"college_code"`
	College      College `json:"college"`
	OriginalCode string  `json:"original_code"`
}

type College struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
