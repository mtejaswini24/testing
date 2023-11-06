package services

import (
	"context"
	"job-portal-api/internal/models"
)

func (s *Service) CreateCompany(ctx context.Context, companyData models.NewCompany) (models.Company, error) {
	//prepare user record
	companyDetails := models.Company{
		CompanyName: companyData.CompanyName,
		Location:    companyData.Location,
	}
	companyDetails, err := s.userRepo.CreateCompany(companyDetails)
	if err != nil {
		return models.Company{}, err
	}
	return companyDetails, nil
}
func (s *Service) FetchCompanies() ([]models.Company, error) {
	companyDetails, err := s.userRepo.ViewCompanies()
	if err != nil {
		return nil, err
	}
	return companyDetails, nil
}
func (s *Service) FetchCompanyById(cid uint64) (models.Company, error) {
	companyDetails, err := s.userRepo.ViewCompanyById(cid)
	if err != nil {
		return models.Company{}, err
	}
	return companyDetails, nil
}
