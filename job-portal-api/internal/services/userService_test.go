package services

import (
	"context"
	"errors"
	"job-portal-api/internal/models"
	"job-portal-api/internal/repository"
	"reflect"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
)

func TestService_CreateUser(t *testing.T) {
	type args struct {
		ctx      context.Context
		userData models.NewUser
	}
	tests := []struct {
		name         string
		args         args
		want         models.User
		wantErr      bool
		mockResponse func() (models.User, error)
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Name:     "doejrnifhnsj",
					Email:    "fjsdijs@gmail.com",
					Password: "12648394",
				},
			},
			want: models.User{
				Name:         "doejrnifhnsj",
				Email:        "fjsdijs@gmail.com",
				PasswordHash: "hashed password",
			},
			wantErr: false,
			mockResponse: func() (models.User, error) {
				return models.User{
					Name:         "doejrnifhnsj",
					Email:        "fjsdijs@gmail.com",
					PasswordHash: "hashed password",
				}, nil
			},
		},
		// {
		// 	name: "failure in passwordhash",
		// 	args: args{
		// 		ctx: context.Background(),
		// 		userData: models.NewUser{
		// 			Name:     "surya",
		// 			Email:    "surya@gmail.com",
		// 			Password: "abcdefghi",
		// 		},
		// 	},
		// 	want:    models.User{},
		// 	wantErr: true,
		// 	mockResponse: func() (models.User, error) {
		// 		return models.User{}, nil
		// 	},
		// },
		{
			name: "failure in creating user",
			args: args{
				ctx: context.Background(),
				userData: models.NewUser{
					Name:     "chandra",
					Email:    "chandragmail.com",
					Password: "19083073987",
				},
			},
			want:    models.User{},
			wantErr: true,
			mockResponse: func() (models.User, error) {
				return models.User{}, errors.New("error while hashing the password")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockResponse != nil {
				mockRepo.EXPECT().CreateUser(gomock.Any()).Return(tt.mockResponse()).AnyTimes()
			}
			s, err := NewService(mockRepo)
			if err != nil {
				t.Errorf("error in initializing the repo layer")
				return
			}
			got, err := s.CreateUser(tt.args.ctx, tt.args.userData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_UserLogin(t *testing.T) {
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name         string
		args         args
		want         jwt.RegisteredClaims
		wantErr      bool
		mockResponse func() (models.User, error)
	}{
		{
			name: "failure in email",
			args: args{
				ctx:      context.Background(),
				email:    "teju@gmail.com",
				password: "12345",
			},
			want:    jwt.RegisteredClaims{},
			wantErr: true,
			mockResponse: func() (models.User, error) {
				return models.User{}, errors.New("error in email")
			},
		},
		{
			name: "failure in password",
			args: args{
				ctx:      context.Background(),
				email:    "teju@gmail.com",
				password: "12345",
			},
			want:    jwt.RegisteredClaims{},
			wantErr: true,
			mockResponse: func() (models.User, error) {
				return models.User{
					Name:         "teju",
					Email:        "teju@gmail.com",
					PasswordHash: "$2a$10$rlvsZ/vfboEmWRH2yhOEy.WpvOVdAk7Xlt47.X4EPWStRtrJJt9V",
				}, nil
			},
		},
		// {
		// 	name: "success in password",
		// 	args: args{
		// 		ctx:      context.Background(),
		// 		email:    "teju@gmail.com",
		// 		password: "12345",
		// 	},
		// 	want: jwt.RegisteredClaims{
		// 		Issuer:    "service project",
		// 		Subject:   "1",
		// 		Audience:  jwt.ClaimStrings{"users"},
		// 		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		// 		IssuedAt:  jwt.NewNumericDate(time.Now()),
		// 	},
		// 	wantErr: false,
		// 	mockResponse: func() (models.User, error) {
		// 		return models.User{
		// 			Name:         "teju",
		// 			Email:        "teju@gmail.com",
		// 			PasswordHash: "$2a$10$rlvsZ/vfboEmWRH2yhOEy.WpvOVdAk7Xlt47.X4EPWStRtrJJt9Vm",
		// 		}, nil
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			mockRepo := repository.NewMockUserRepo(mc)
			if tt.mockResponse != nil {
				mockRepo.EXPECT().UserLogin(tt.args.email).Return(tt.mockResponse()).AnyTimes()
			}

			s, err := NewService(mockRepo)
			if err != nil {
				t.Errorf("error in initializing the repo layer")
				return
			}
			got, err := s.UserLogin(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UserLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.UserLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}
