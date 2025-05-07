package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	SessionID    string    `json:"session_id"    gorm:"primaryKey; type:varchar(36)"`
	UserID       string    `json:"user_id"       gorm:"type:varchar(36); not null"`
	IP           string    `json:"ip"            gorm:"type:varchar(45)"`
	UserAgent    string    `json:"user_agent"    gorm:"type:varchar(512)"`
	RefreshToken string    `json:"refresh_token" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at"    gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at"    gorm:"autoUpdateTime"`
	ExpireAt     time.Time `json:"expire_at"     gorm:"not null"`
}

type UserResponse struct {
	UserID string `json:"user_id"`
}

type SignOutResponse struct {
	Message string `json:"message"`
}

// Adds a new session to the sessions table.
func CreateSession(db *gorm.DB, session *Session) (string, error) {
	if session.SessionID == "" {
		session.SessionID = uuid.New().String()
	}

	if err := db.Create(session).Error; err != nil {
		return "", err
	}

	return session.SessionID, nil
}

// Retrieves the session's data from the sessions table by their identifier.
func GetSession(db *gorm.DB, sessionID string) (session *Session, err error) {
	err = db.Where("session_id = ?", sessionID).First(&session).Error
	return session, err
}

// Updates the session's data in the sessions table.
func UpdateSession(db *gorm.DB, session *Session) error {
	return db.Model(&Session{}).Where("session_id = ?", session.SessionID).Updates(session).Error
}

// Removes a session from the sessions table by their identifier.
func DeleteSession(db *gorm.DB, sessionID string) error {
	return db.Where("session_id = ?", sessionID).Delete(&Session{}).Error
}
