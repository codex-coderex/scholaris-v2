package main

import (
	"context"
	"errors"
	"fmt"

	"scholaris-v2/internal/config"
	"scholaris-v2/internal/db"
	"scholaris-v2/internal/features/colleges"
	"scholaris-v2/internal/features/programs"
	"scholaris-v2/internal/features/students"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// service

type AppService struct {
	collegeRepo *colleges.CollegeRepository
	programRepo *programs.ProgramRepository
	studentRepo *students.StudentRepository
	pool        *pgxpool.Pool
	startupErr  error
}

var errNoDatabase = errors.New("no database")

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
		a.startupErr = fmt.Errorf("%w: %v", errNoDatabase, err)
		return nil
	}

	if err := db.CreateTables(pool); err != nil {
		a.startupErr = fmt.Errorf("%w: %v", errNoDatabase, err)
		db.Close(pool)
		return nil
	}

	if err := db.Seed(pool); err != nil {
		a.startupErr = fmt.Errorf("%w: %v", errNoDatabase, err)
		db.Close(pool)
		return nil
	}

	a.collegeRepo = colleges.NewCollegeRepository(pool)
	a.programRepo = programs.NewProgramRepository(pool)
	a.studentRepo = students.NewStudentRepository(pool)
	a.pool = pool
	a.startupErr = nil

	return nil
}

func (a *AppService) ServiceShutdown() error {
	db.Close(a.pool)
	a.pool = nil
	a.collegeRepo = nil
	a.programRepo = nil
	a.studentRepo = nil
	return nil
}

func (a *AppService) databaseError() error {
	if a.startupErr != nil {
		return a.startupErr
	}

	return errNoDatabase
}

// colleges

func (a *AppService) GetColleges(search, sortBy, order string, page, pageSize int) ([]colleges.College, int, error) {
	if a.collegeRepo == nil {
		return []colleges.College{}, 0, a.databaseError()
	}

	return a.collegeRepo.GetAll(search, sortBy, order, page, pageSize)
}

func (a *AppService) CreateCollege(c colleges.College) error {
	if a.collegeRepo == nil {
		return a.databaseError()
	}

	return a.collegeRepo.Create(c)
}

func (a *AppService) UpdateCollege(c colleges.College) error {
	if a.collegeRepo == nil {
		return a.databaseError()
	}

	return a.collegeRepo.Update(c)
}

func (a *AppService) DeleteCollege(code string) error {
	if a.collegeRepo == nil {
		return a.databaseError()
	}

	return a.collegeRepo.Delete(code)
}

// programs

func (a *AppService) GetPrograms(search, sortBy, order string, page, pageSize int, collegeCode string) ([]programs.Program, int, error) {
	if a.programRepo == nil {
		return []programs.Program{}, 0, a.databaseError()
	}

	return a.programRepo.GetAll(search, sortBy, order, page, pageSize, collegeCode)
}

func (a *AppService) CreateProgram(p programs.Program) error {
	if a.programRepo == nil {
		return a.databaseError()
	}

	return a.programRepo.Create(p)
}

func (a *AppService) UpdateProgram(p programs.Program) error {
	if a.programRepo == nil {
		return a.databaseError()
	}

	return a.programRepo.Update(p)
}

func (a *AppService) DeleteProgram(code string) error {
	if a.programRepo == nil {
		return a.databaseError()
	}

	return a.programRepo.Delete(code)
}

// students

func (a *AppService) GetStudents(search, sortBy, order string, page, pageSize int) ([]students.Student, int, error) {
	if a.studentRepo == nil {
		return []students.Student{}, 0, a.databaseError()
	}

	return a.studentRepo.GetAll(search, sortBy, order, page, pageSize)
}

func (a *AppService) CreateStudent(s students.Student) error {
	if a.studentRepo == nil {
		return a.databaseError()
	}

	return a.studentRepo.Create(s)
}

func (a *AppService) UpdateStudent(s students.Student) error {
	if a.studentRepo == nil {
		return a.databaseError()
	}

	return a.studentRepo.Update(s)
}

func (a *AppService) DeleteStudent(id string) error {
	if a.studentRepo == nil {
		return a.databaseError()
	}

	return a.studentRepo.Delete(id)
}
