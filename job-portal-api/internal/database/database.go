package database

import (
	"job-portal-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.Migrator().AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}
	err = db.Migrator().AutoMigrate(&models.Company{})
	if err != nil {
		return nil, err
	}
	err = db.Migrator().AutoMigrate(&models.Job{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
