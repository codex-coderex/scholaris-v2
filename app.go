package main

import (
	"context"
	"fmt"

	"scholaris-v2/internal/config"
	"scholaris-v2/internal/db"
	"scholaris-v2/internal/models"
	"scholaris-v2/internal/repository"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// service

type AppService struct {
	collegeRepo *repository.CollegeRepository
	programRepo *repository.ProgramRepository
	studentRepo *repository.StudentRepository
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

	a.collegeRepo = repository.NewCollegeRepository(pool)
	a.programRepo = repository.NewProgramRepository(pool)
	a.studentRepo = repository.NewStudentRepository(pool)

	return nil
}

// colleges

func (a *AppService) GetColleges(search, sortBy, order string, page, pageSize int) ([]models.College, int, error) {
	return a.collegeRepo.GetAll(search, sortBy, order, page, pageSize)
}

func (a *AppService) GetCollege(code string) (models.College, error) {
	return a.collegeRepo.GetByCode(code)
}

func (a *AppService) CreateCollege(c models.College) error {
	return a.collegeRepo.Create(c)
}

func (a *AppService) UpdateCollege(c models.College) error {
	return a.collegeRepo.Update(c)
}

func (a *AppService) DeleteCollege(code string) error {
	return a.collegeRepo.Delete(code)
}

// programs

func (a *AppService) GetPrograms(search, sortBy, order string, page, pageSize int) ([]models.Program, int, error) {
	return a.programRepo.GetAll(search, sortBy, order, page, pageSize)
}

func (a *AppService) GetProgram(code string) (models.Program, error) {
	return a.programRepo.GetByCode(code)
}

func (a *AppService) CreateProgram(p models.Program) error {
	return a.programRepo.Create(p)
}

func (a *AppService) UpdateProgram(p models.Program) error {
	return a.programRepo.Update(p)
}

func (a *AppService) DeleteProgram(code string) error {
	return a.programRepo.Delete(code)
}

// students

func (a *AppService) GetStudents(search, sortBy, order string, page, pageSize int) ([]models.Student, int, error) {
	return a.studentRepo.GetAll(search, sortBy, order, page, pageSize)
}

func (a *AppService) GetStudent(id string) (models.Student, error) {
	return a.studentRepo.GetById(id)
}

func (a *AppService) CreateStudent(s models.Student) error {
	return a.studentRepo.Create(s)
}

func (a *AppService) UpdateStudent(s models.Student) error {
	return a.studentRepo.Update(s)
}

func (a *AppService) DeleteStudent(id string) error {
	return a.studentRepo.Delete(id)
}
