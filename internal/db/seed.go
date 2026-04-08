package db

import (
	"context"
	"fmt"
	"math/rand"
	"time"

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
	{"COE", "College of Engineering"},
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
	{"BSBio-AB", "Bachelor of Science in Biology Major in Animal Biology", "CSM"},
	{"BSBio-Micro", "Bachelor of Science in Biology Major in Microbiology", "CSM"},
	{"BSBio-Bdv", "Bachelor of Science in Biology Major in Biodiversity", "CSM"},
	{"BSChem", "Bachelor of Science in Chemistry", "CSM"},
	{"BSMath", "Bachelor of Science in Mathematics", "CSM"},
	{"BSPhys", "Bachelor of Science in Physics", "CSM"},
	{"BSStat", "Bachelor of Science in Statistics", "CSM"},
	{"BS-MB", "Bachelor of Science in Marine Biology", "CSM"},
	// CED
	{"BTLEd-IA", "Bachelor of Technology and Livelihood Education Major in Industrial Arts", "CED"},
	{"BTLEd-HE", "Bachelor of Technology and Livelihood Education Major in Home Economics", "CED"},
	{"BTVTEd-DT", "Bachelor of Technical-Vocational Teacher Education Major in Drafting Technology", "CED"},
	{"BEEd-SM", "Bachelor of Elementary Education Major in Science and Mathematics", "CED"},
	{"BEEd-LE", "Bachelor of Elementary Education Major in Language Education", "CED"},
	{"BSEd-Phys", "Bachelor of Secondary Education Major in Physics", "CED"},
	{"BSEd-Math", "Bachelor of Secondary Education Major in Mathematics", "CED"},
	{"BSEd-Chem", "Bachelor of Secondary Education Major in Chemistry", "CED"},
	{"BSEd-Fil", "Bachelor of Secondary Education Major in Filipino", "CED"},
	{"BSEd-Bio", "Bachelor of Secondary Education Major in Biology", "CED"},
	{"BPEd", "Bachelor of Physical Education", "CED"},
	// CASS
	{"BAELS", "Bachelor of Arts in English Language Studies", "CASS"},
	{"BAPos", "Bachelor of Arts in Political Science", "CASS"},
	{"BSPsych", "Bachelor of Science in Psychology", "CASS"},
	{"BAPsych", "Bachelor of Arts in Psychology", "CASS"},
	{"BASoc", "Bachelor of Arts in Sociology", "CASS"},
	{"BAHis", "Bachelor of Arts in History", "CASS"},
	{"BALCS", "Bachelor of Arts in Literary and Cultural Studies", "CASS"},
	{"BSPhil", "Bachelor of Science in Philosophy Applied Ethics", "CASS"},
	{"BAPan", "Batsilyer ng Sining sa Panitikan", "CASS"},
	{"BAFil", "Batsilyer ng Sining sa Filipino", "CASS"},
	// CHS
	{"BSN", "Bachelor of Science in Nursing", "CHS"},
	// CEBA
	{"BSA", "Bachelor of Science in Accountancy", "CEBA"},
	{"BSBA-MM", "Bachelor of Science in Business Administration Major in Marketing Management", "CEBA"},
	{"BSBA-BE", "Bachelor of Science in Business Administration Major in Business Economics", "CEBA"},
	{"BSHM", "Bachelor of Science in Hospitality Management", "CEBA"},
	{"BSEcon", "Bachelor of Science in Economics", "CEBA"},
	{"BSEntrep", "Bachelor of Science in Entrepreneurship", "CEBA"},
	// COE
	{"BSCE", "Bachelor of Science in Civil Engineering", "COE"},
	{"BSEE", "Bachelor of Science in Electrical Engineering", "COE"},
	{"BSME", "Bachelor of Science in Mechanical Engineering", "COE"},
	{"BSCHE", "Bachelor of Science in Chemical Engineering", "COE"},
	{"BSCPE", "Bachelor of Science in Computer Engineering", "COE"},
	{"BSEsE", "Bachelor of Science in Electronics Engineering", "COE"},
	{"BSEnE", "Bachelor of Science in Environmental Engineering", "COE"},
	{"BSMETE", "Bachelor of Science in Metallurgical Engineering", "COE"},
	{"BSEM", "Bachelor of Science in Mining Engineering", "COE"},
	{"BSCERE", "Bachelor of Science in Ceramic Engineering", "COE"},
	{"BSIAM", "Bachelor of Science in Industrial Automation and Mechatronics", "COE"},
	{"BET-CHET", "Bachelor of Engineering Technology Major in Chemical Engineering Technology", "COE"},
	{"BET-ELET", "Bachelor of Engineering Technology Major in Electrical Engineering Technology", "COE"},
	{"BET-MET", "Bachelor of Engineering Technology Major in Mechanical Engineering Technology", "COE"},
	{"BET-ESET", "Bachelor of Engineering Technology Major in Electronics Engineering Technology", "COE"},
	{"BET-MMT", "Bachelor of Engineering Technology Major in Metallurgical and Materials Engineering Technology", "COE"},
	{"BET-CET", "Bachelor of Engineering Technology Major in Civil Engineering Technology", "COE"},
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

	// students — only seed if table is empty
	var count int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM student`).Scan(&count); err != nil {
		return fmt.Errorf("seed count check: %w", err)
	}

	if count == 0 {
		programCodes := make([]string, len(programs))
		for i, p := range programs {
			programCodes[i] = p.code
		}

		for i := range 7500 {
			currentYear := time.Now().Year()
			id := fmt.Sprintf("%04d-%04d", rand.Intn(currentYear-1968+1)+1968, rand.Intn(currentYear-1968+1)+1968)
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
