package models

import (
	"gorm.io/gorm"
)

type NewCompany struct {
	CompanyName string `json:"companyname" validate:"required" `
	Location    string `json:"location" validate:"required"`
}
type Company struct {
	gorm.Model
	CompanyName string `validate:"required,unique" gorm:"unique;not null"`
	Location    string `json:"location"`
}

type NewJob struct {
	JobRole        string `json:"Role" validate:"required"`
	JobDescription string `json:"Description" validate:"required"`
}
type Job struct {
	gorm.Model
	Company        Company `json:"-" gorm:"ForeignKey:cid"`
	Cid            uint    `json:"cid"`
	JobRole        string  `json:"Role"`
	JobDescription string  `json:"Description"`
}
