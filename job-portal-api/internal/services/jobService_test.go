package services

import (
	"context"
	"errors"
	"job-portal-api/internal/models"
	"job-portal-api/internal/repository"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestService_FetchJobById(t *testing.T) {
	type args struct {
		jid uint64
	}
	tests := []struct {
		name             string
		args             args
		want             models.Job
		wantErr          bool
		mockRepoResponse func() (models.Job, error)
	}{
		{
			name: "success",
			want: models.Job{
				Company: models.Company{
					CompanyName: "hp",
					Location:    "musuru",
				},
				Cid:            1,
				JobRole:        "front-end",
				JobDescription: "something",
			},
			args: args{
				jid: 15,
			},
			wantErr: false,
			mockRepoResponse: func() (models.Job, error) {
				return models.Job{
					Company: models.Company{
						CompanyName: "hp",
						Location:    "musuru",
					},
					Cid:            1,
					JobRole:        "front-end",
					JobDescription: "something",
				}, nil
			},
		},
		{
			name: "invalid job id",
			want: models.Job{},
			args: args{
				jid: 5,
			},
			mockRepoResponse: func() (models.Job, error) {
				return models.Job{}, errors.New("error test")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().ViewJobDetailsById(tt.args.jid).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, err := NewService(mockRepo)
			if err != nil {
				t.Errorf("error in initializing the repo layer")
				return
			}
			got, err := s.FetchJobById(tt.args.jid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.FetchJobById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.FetchJobById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_FetchJob(t *testing.T) {
	tests := []struct {
		name             string
		want             []models.Job
		wantErr          bool
		mockRepoResponse func() ([]models.Job, error)
	}{
		{
			name: "database success",
			want: []models.Job{
				{
					Cid:            1,
					JobRole:        "test",
					JobDescription: "something",
				},
				{
					Cid:            2,
					JobRole:        "putvi",
					JobDescription: "having",
				},
			},
			wantErr: false,
			mockRepoResponse: func() ([]models.Job, error) {
				return []models.Job{
					{
						Cid:            1,
						JobRole:        "test",
						JobDescription: "something",
					},
					{
						Cid:            2,
						JobRole:        "putvi",
						JobDescription: "having",
					},
				}, nil
			},
		},
		{
			name:    "database failure",
			want:    nil,
			wantErr: true,
			mockRepoResponse: func() ([]models.Job, error) {
				return nil, errors.New("error test")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().ViewAllJobs().Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, err := NewService(mockRepo)
			if err != nil {
				t.Errorf("error in initializing the repo layer")
				return
			}
			got, err := s.FetchJob()
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.FetchJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.FetchJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_FetchJobByCompanyId(t *testing.T) {
	type args struct {
		cid uint64
	}
	tests := []struct {
		name             string
		args             args
		want             []models.Job
		wantErr          bool
		mockRepoResponse func() ([]models.Job, error)
	}{
		{
			name: "success",
			args: args{
				cid: 2,
			},
			want: []models.Job{
				{
					Cid:            2,
					JobRole:        "xyz",
					JobDescription: "sleeping",
				},
				{
					Cid:            2,
					JobRole:        "abc",
					JobDescription: "walking",
				},
			},
			wantErr: false,
			mockRepoResponse: func() ([]models.Job, error) {
				return []models.Job{
					{
						Cid:            2,
						JobRole:        "xyz",
						JobDescription: "sleeping",
					},
					{
						Cid:            2,
						JobRole:        "abc",
						JobDescription: "walking",
					},
				}, nil
			},
		},
		{
			name: "failure",
			args: args{
				cid: 10,
			},
			want:    nil,
			wantErr: true,
			mockRepoResponse: func() ([]models.Job, error) {
				return nil, errors.New("data is not there")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().ViewJobByCompanyId(tt.args.cid).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, err := NewService(mockRepo)
			if err != nil {
				t.Errorf("error in initializing the repo layer")
				return
			}
			got, err := s.FetchJobByCompanyId(tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.FetchJobByCompanyId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.FetchJobByCompanyId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_CreateJob(t *testing.T) {
	type args struct {
		ctx     context.Context
		jobData models.NewJob
		cid     uint64
	}
	tests := []struct {
		name             string
		args             args
		want             models.Job
		wantErr          bool
		mockRepoResponse func() (models.Job, error)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				jobData: models.NewJob{
					JobRole:        "abc",
					JobDescription: "do something",
				},
				cid: 1,
			},
			want: models.Job{
				Cid:            1,
				JobRole:        "abc",
				JobDescription: "do something",
			},
			wantErr: false,
			mockRepoResponse: func() (models.Job, error) {
				return models.Job{
					Cid:            1,
					JobRole:        "abc",
					JobDescription: "do something",
				}, nil
			},
		},
		{
			name: "failure",
			args: args{
				ctx: context.Background(),
				jobData: models.NewJob{
					JobRole:        "red",
					JobDescription: "frfytfguy",
				},
				cid: 2,
			},
			want:    models.Job{},
			wantErr: true,
			mockRepoResponse: func() (models.Job, error) {
				return models.Job{}, errors.New("job is not created")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockRepoResponse != nil {
				mockRepo.EXPECT().CreateJob(gomock.Any()).Return(tt.mockRepoResponse()).AnyTimes()
			}
			s, err := NewService(mockRepo)
			if err != nil {
				t.Errorf("error in initializing the repo layer")
				return
			}
			got, err := s.CreateJob(tt.args.ctx, tt.args.jobData, tt.args.cid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.CreateJob() = %v, want %v", got, tt.want)
			}
		})
	}
}
