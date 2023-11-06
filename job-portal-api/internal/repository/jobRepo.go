package repository

import (
	"errors"
	"job-portal-api/internal/models"

	"github.com/rs/zerolog/log"
)

func (r *Repo) CreateJob(jobData models.Job) (models.Job, error) {
	err := r.db.Create(&jobData).Error
	//calling default create method
	if err != nil {
		log.Info().Err(err).Send()
		return models.Job{}, errors.New("unable to create job")
	}
	return jobData, nil
}
func (r *Repo) ViewJobByCompanyId(cid uint64) ([]models.Job, error) {
	var jobData []models.Job
	result := r.db.Where("cid=?", cid).Find(&jobData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, errors.New("jobs are not there with that company id")
	}
	return jobData, nil
}
func (r *Repo) ViewAllJobs() ([]models.Job, error) {
	var jobDetails []models.Job
	result := r.db.Find(&jobDetails)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return nil, errors.New("jobs data is empty")
	}
	return jobDetails, nil
}
func (r *Repo) ViewJobDetailsById(jid uint64) (models.Job, error) {
	var jobData models.Job
	result := r.db.Where("id=?", jid).First(&jobData)
	if result.Error != nil {
		log.Info().Err(result.Error).Send()
		return models.Job{}, errors.New("job is not there with that id")
	}
	return jobData, nil
}
