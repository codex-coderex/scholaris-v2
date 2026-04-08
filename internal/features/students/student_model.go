package students

// student

type Student struct {
	Id          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Year        int     `json:"year"`
	Gender      string  `json:"gender"`
	ProgramCode string  `json:"program_code"`
	Program     Program `json:"program"`
}

type Program struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	CollegeCode string  `json:"college_code"`
	College     College `json:"college"`
}

type College struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
