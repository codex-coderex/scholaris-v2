package models

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
