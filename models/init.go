package models

import (
	"fmt"
	"simpleAuth/config"
	"simpleAuth/logger"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Create new connection to DB
func NewDBConnection(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.CustomGormLogger(),
	})
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	logrus.Info("Successfully connected to the database")

	var tables []string
	result := DB.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables)
	if result.Error != nil {
		logrus.WithError(result.Error).Fatal("Failed retrieving table list")
	}

	found := false
	for _, table := range tables {
		if table == "sessions" {
			found = true
			break
		}
	}

	if !found {
		logrus.Fatal("Session table in db not found")
	}

	return DB
}
