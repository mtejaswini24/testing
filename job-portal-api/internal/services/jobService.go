package services

import (
	"context"
	"job-portal-api/internal/models"
)

func (s *Service) CreateJob(ctx context.Context, jobData models.NewJob, cid uint64) (models.Job, error) {
	jobDetails := models.Job{
		JobRole:        jobData.JobRole,
		JobDescription: jobData.JobDescription,
		Cid:            uint(cid),
	}
	jobDetails, err := s.userRepo.CreateJob(jobDetails)
	if err != nil {
		return models.Job{}, err
	}
	return jobDetails, nil
}
func (s *Service) FetchJob() ([]models.Job, error) {
	jobDetails, err := s.userRepo.ViewAllJobs()
	if err != nil {
		return nil, err
	}
	return jobDetails, nil
}
func (s *Service) FetchJobById(jid uint64) (models.Job, error) {
	jobDetails, err := s.userRepo.ViewJobDetailsById(jid)
	if err != nil {
		return models.Job{}, err
	}
	return jobDetails, nil
}
func (s *Service) FetchJobByCompanyId(cid uint64) ([]models.Job, error) {
	jobDetails, err := s.userRepo.ViewJobByCompanyId(cid)
	if err != nil {
		return nil, err
	}
	return jobDetails, nil
}
