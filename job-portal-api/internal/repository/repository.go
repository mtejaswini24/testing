package repository

import (
	"errors"

	"job-portal-api/internal/models"

	"gorm.io/gorm"
)

type Repo struct {

	// db is an instance of the SQLite database.
	db *gorm.DB
}

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=repository
type UserRepo interface {
	CreateUser(userData models.User) (models.User, error)
	UserLogin(email string) (models.User, error)

	CreateCompany(companyData models.Company) (models.Company, error)
	ViewCompanies() ([]models.Company, error)
	ViewCompanyById(cid uint64) (models.Company, error)

	CreateJob(jobData models.Job) (models.Job, error)
	ViewJobByCompanyId(cid uint64) ([]models.Job, error)
	ViewAllJobs() ([]models.Job, error)
	ViewJobDetailsById(jid uint64) (models.Job, error)
}

func NewRepo(db *gorm.DB) (UserRepo, error) {
	if db == nil {
		return nil, errors.New("database cannot be null")
	}
	return &Repo{
		db: db,
	}, nil
}
