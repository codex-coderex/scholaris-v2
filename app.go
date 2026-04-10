package main

import (
	"context"
	"fmt"

	"scholaris-v2/internal/config"
	"scholaris-v2/internal/db"
	"scholaris-v2/internal/features/colleges"
	"scholaris-v2/internal/features/programs"
	"scholaris-v2/internal/features/students"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// service

type AppService struct {
	collegeRepo *colleges.CollegeRepository
	programRepo *programs.ProgramRepository
	studentRepo *students.StudentRepository
}

func NewAppService() *AppService {
	return &AppService{}
}

func (a *AppService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	pool, err := db.Init(cfg)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}

	if err := db.CreateTables(pool); err != nil {
		return fmt.Errorf("create tables: %w", err)
	}

	if err := db.Seed(pool); err != nil {
		return fmt.Errorf("seed: %w", err)
	}

	a.collegeRepo = colleges.NewCollegeRepository(pool)
	a.programRepo = programs.NewProgramRepository(pool)
	a.studentRepo = students.NewStudentRepository(pool)

	return nil
}

// colleges

func (a *AppService) GetColleges(search, sortBy, order string, page, pageSize int) ([]colleges.College, int, error) {
	return a.collegeRepo.GetAll(search, sortBy, order, page, pageSize)
}

func (a *AppService) GetCollege(code string) (colleges.College, error) {
	return a.collegeRepo.GetByCode(code)
}

func (a *AppService) CreateCollege(c colleges.College) error {
	return a.collegeRepo.Create(c)
}

func (a *AppService) UpdateCollege(c colleges.College) error {
	return a.collegeRepo.Update(c)
}

func (a *AppService) DeleteCollege(code string) error {
	return a.collegeRepo.Delete(code)
}

// programs

func (a *AppService) GetPrograms(search, sortBy, order string, page, pageSize int, collegeCode string) ([]programs.Program, int, error) {
	return a.programRepo.GetAll(search, sortBy, order, page, pageSize, collegeCode)
}

func (a *AppService) GetProgram(code string) (programs.Program, error) {
	return a.programRepo.GetByCode(code)
}

func (a *AppService) CreateProgram(p programs.Program) error {
	return a.programRepo.Create(p)
}

func (a *AppService) UpdateProgram(p programs.Program) error {
	return a.programRepo.Update(p)
}

func (a *AppService) DeleteProgram(code string) error {
	return a.programRepo.Delete(code)
}

// students

func (a *AppService) GetStudents(search, sortBy, order string, page, pageSize int) ([]students.Student, int, error) {
	return a.studentRepo.GetAll(search, sortBy, order, page, pageSize)
}

func (a *AppService) GetStudent(id string) (students.Student, error) {
	return a.studentRepo.GetById(id)
}

func (a *AppService) CreateStudent(s students.Student) error {
	return a.studentRepo.Create(s)
}

func (a *AppService) UpdateStudent(s students.Student) error {
	return a.studentRepo.Update(s)
}

func (a *AppService) DeleteStudent(id string) error {
	return a.studentRepo.Delete(id)
}
