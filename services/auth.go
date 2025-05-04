package services

import (
	"simpleAuth/models"

	"gorm.io/gorm"
)

func SignOut(db *gorm.DB, sessionID string) error {
	return models.DeleteSession(db, sessionID)
}
