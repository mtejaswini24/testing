package services

import (
	"context"
	"errors"
	"job-portal-api/internal/models"
	"job-portal-api/internal/pkg"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func (s *Service) CreateUser(ctx context.Context, userData models.NewUser) (models.User, error) {
	//method that creates a new record in  db
	hashedPass, err := pkg.HashPassword(userData.Password)
	if err != nil {
		return models.User{}, err
	}
	//prepare user record
	userDetails := models.User{
		Name:         userData.Name,
		Email:        userData.Email,
		PasswordHash: string(hashedPass),
	}
	userDetails, err = s.userRepo.CreateUser(userDetails)
	if err != nil {
		return models.User{}, err
	}
	return userDetails, nil
}
func (s *Service) UserLogin(ctx context.Context, email, password string) (jwt.RegisteredClaims, error) {
	//checking the email in database
	userDetails, err := s.userRepo.UserLogin(email)
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}
	// We check if the provided password matches the hashed password in the database.
	err = pkg.CheckPassword(password, userDetails.PasswordHash)
	if err != nil {
		log.Info().Err(err).Send()
		return jwt.RegisteredClaims{}, errors.New("entered password is wrong")
	}
	// Successful authentication! Generate JWT claims.
	claims := jwt.RegisteredClaims{
		Issuer:    "service project",
		Subject:   strconv.FormatUint(uint64(userDetails.ID), 10),
		Audience:  jwt.ClaimStrings{"users"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	// And return those claims.
	return claims, nil

	// token, err := s.auth.GenerateAuthToken(claims)
	// if err != nil {
	// 	return "", err
	// }

	// return token, nil
}
