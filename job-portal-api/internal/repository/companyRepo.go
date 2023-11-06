package repository

import (
	"errors"
	"job-portal-api/internal/models"

	"github.com/rs/zerolog/log"
)

func (r *Repo) CreateCompany(companyData models.Company) (models.Company, error) {
	err := r.db.Create(&companyData).Error
	//calling default create method
	if err != nil {
		log.Info().Err(err).Send()
		return models.Company{}, errors.New("unable to create company")
	}
	return companyData, nil
}
func (r *Repo) ViewCompanies() ([]models.Company, error) {
	var companyDetails []models.Company
	result := r.db.Find(&companyDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, errors.New("companies details is not there")
	}
	return companyDetails, nil
}
func (r *Repo) ViewCompanyById(cid uint64) (models.Company, error) {
	var companyData models.Company
	result := r.db.Where("id=?", cid).First(&companyData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.Company{}, errors.New("company data is not there with that id")
	}
	return companyData, nil
}
