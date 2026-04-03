package db

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/jackc/pgx/v5/pgxpool"
)

// data

var firstNames = []string{
	"James", "Maria", "John", "Ana", "Robert", "Karen", "Michael", "Lisa",
	"William", "Sandra", "David", "Ashley", "Richard", "Dorothy", "Joseph",
	"Jessica", "Thomas", "Emily", "Charles", "Sarah",
}

var lastNames = []string{
	"Santos", "Reyes", "Cruz", "Garcia", "Mendoza", "Torres", "Flores",
	"Rivera", "Gomez", "Diaz", "Martinez", "Hernandez", "Lopez", "Gonzalez",
	"Perez", "Rodriguez", "Ramos", "Ramirez", "Castillo", "Morales",
}

var genders = []string{"Male", "Female"}

var colleges = []struct {
	code string
	name string
}{
	{"CCS", "College of Computer Studies"},
	{"CED", "College of Education"},
	{"CASS", "College of Arts and Social Sciences"},
	{"CSM", "College of Science and Mathematics"},
	{"CHS", "College of Health Sciences"},
	{"CEBA", "College of Economics, Business and Accountancy"},
}

var programs = []struct {
	code        string
	name        string
	collegeCode string
}{
	// CCS
	{"BSCS", "Bachelor of Science in Computer Science", "CCS"},
	{"BSIT", "Bachelor of Science in Information Technology", "CCS"},
	{"BSIS", "Bachelor of Science in Information System", "CCS"},
	{"BSCA", "Bachelor of Science in Computer Application", "CCS"},
	// CSM
	{"BSBio-PB", "Bachelor of Science in Biology Major in Plant Biology", "CSM"},
	{"BSChem", "Bachelor of Science in Chemistry", "CSM"},
	{"BSBio-Micro", "Bachelor of Science in Biology Major in Microbiology", "CSM"},
	{"BSMath", "Bachelor of Science in Mathematics", "CSM"},
	{"BSBio-Bdv", "Bachelor of Science in Biology Major in Biodiversity", "CSM"},
	{"BSPhys", "Bachelor of Science in Physics", "CSM"},
	// CED
	{"BTLEd-IA", "Bachelor of Technology and Livelihood Education Major in Industrial Arts", "CED"},
	{"BEEd-SM", "Bachelor of Elementary Education Major in Science and Mathematics", "CED"},
	{"BSEd-Phys", "Bachelor of Secondary Education Major in Physics", "CED"},
	{"BSEd-Math", "Bachelor of Secondary Education Major in Mathematics", "CED"},
	{"BSEd-Chem", "Bachelor of Secondary Education Major in Chemistry", "CED"},
	{"BTLEd-HE", "Bachelor of Technology and Livelihood Education Major in Home Economics", "CED"},
	{"BEEd-LE", "Bachelor of Elementary Education Major in Language Education", "CED"},
	{"BPEd", "Bachelor of Physical Education", "CED"},
	{"BSEd-Fil", "Bachelor of Secondary Education Major in Filipino", "CED"},
	{"BSEd-Bio", "Bachelor of Secondary Education Major in Biology", "CED"},
	// CASS
	{"BAELS", "Bachelor of Arts in English Language Studies", "CASS"},
	{"BAPos", "Bachelor of Arts in Political Science", "CASS"},
	{"BSPsych", "Bachelor of Science in Psychology", "CASS"},
	{"BAPsych", "Bachelor of Arts in Psychology", "CASS"},
	{"BASoc", "Bachelor of Arts in Sociology", "CASS"},
	{"BAHis", "Bachelor of Arts in History", "CASS"},
	{"BALCS", "Bachelor of Arts in Literary and Cultural Studies", "CASS"},
	// CHS
	{"BSN", "Bachelor of Science in Nursing", "CHS"},
	// CEBA
	{"BSA", "Bachelor of Science in Accountancy", "CEBA"},
	{"BSBA-MM", "Bachelor of Science in Business Administration Major in Marketing Management", "CEBA"},
	{"BSHM", "Bachelor of Science in Hospitality Management", "CEBA"},
	{"BSEcon", "Bachelor of Science in Economics", "CEBA"},
	{"BSEntrep", "Bachelor of Science in Entrepreneurship", "CEBA"},
}

// seed

func Seed(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// colleges
	for _, c := range colleges {
		if _, err := pool.Exec(ctx, `
			INSERT INTO college (code, name)
			VALUES ($1, $2)
			ON CONFLICT (code) DO NOTHING
		`, c.code, c.name); err != nil {
			return fmt.Errorf("seed college %s: %w", c.code, err)
		}
	}

	// programs
	for _, p := range programs {
		if _, err := pool.Exec(ctx, `
			INSERT INTO program (code, name, college_code)
			VALUES ($1, $2, $3)
			ON CONFLICT (code) DO NOTHING
		`, p.code, p.name, p.collegeCode); err != nil {
			return fmt.Errorf("seed program %s: %w", p.code, err)
		}
	}

	// students — only seed if empty
	var count int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM student`).Scan(&count); err != nil {
		return fmt.Errorf("seed count check: %w", err)
	}

	if count < 5000 {
		programCodes := make([]string, len(programs))
		for i, p := range programs {
			programCodes[i] = p.code
		}

		for i := range 5000 {
			id := fmt.Sprintf("%04d-%04d", rand.Intn(9000)+1000, rand.Intn(9000)+1000)
			firstName := firstNames[rand.Intn(len(firstNames))]
			lastName := lastNames[rand.Intn(len(lastNames))]
			year := rand.Intn(4) + 1
			gender := genders[rand.Intn(len(genders))]
			programCode := programCodes[rand.Intn(len(programCodes))]

			if _, err := pool.Exec(ctx, `
				INSERT INTO student (id, first_name, last_name, year, gender, program_code)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT (id) DO NOTHING
			`, id, firstName, lastName, year, gender, programCode); err != nil {
				return fmt.Errorf("seed student %d: %w", i, err)
			}
		}
	}

	return nil
}
