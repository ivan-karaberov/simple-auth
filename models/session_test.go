package models

import (
	"simpleAuth/logger"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.CustomGormLogger(),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Session{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestCreateSession(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	session := &Session{
		UserID:       uuid.New().String(),
		IP:           "192.168.1.1",
		UserAgent:    "test-agent",
		RefreshToken: "test-refresh-token",
		ExpireAt:     time.Now().Add(24 * time.Hour),
	}

	sessionID, err := CreateSession(db, session)
	assert.NoError(t, err)
	assert.NotEmpty(t, sessionID)

	var createdSession Session
	err = db.First(&createdSession, "session_id = ?", sessionID).Error
	assert.NoError(t, err)
	assert.Equal(t, session.UserID, createdSession.UserID)
	assert.Equal(t, session.IP, createdSession.IP)
	assert.Equal(t, session.UserAgent, createdSession.UserAgent)
}

func TestGetSession(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	session := &Session{
		SessionID:    uuid.New().String(),
		UserID:       uuid.New().String(),
		IP:           "192.168.1.1",
		UserAgent:    "test-agent",
		RefreshToken: "test-refresh-token",
		ExpireAt:     time.Now().Add(24 * time.Hour),
	}

	db.Create(session)

	retrievedSession, err := GetSession(db, session.SessionID)
	assert.NoError(t, err)
	assert.Equal(t, session.SessionID, retrievedSession.SessionID)
	assert.Equal(t, session.IP, retrievedSession.IP)
}

func TestUpdateClient(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	session := &Session{
		SessionID:    uuid.New().String(),
		UserID:       uuid.New().String(),
		IP:           "192.168.1.1",
		UserAgent:    "test-agent",
		RefreshToken: "test-refresh-token",
		ExpireAt:     time.Now().Add(24 * time.Hour),
	}

	db.Create(session)

	session.UserAgent = "updated-agent"
	err = UpdateClient(db, session)
	assert.NoError(t, err)

	retrievedSession, err := GetSession(db, session.SessionID)
	assert.NoError(t, err)
	assert.Equal(t, "updated-agent", retrievedSession.UserAgent)
}

func TestDeleteSession(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	session := &Session{
		SessionID:    uuid.New().String(),
		UserID:       uuid.New().String(),
		IP:           "192.168.1.1",
		UserAgent:    "test-agent",
		RefreshToken: "test-refresh-token",
		ExpireAt:     time.Now().Add(24 * time.Hour),
	}

	db.Create(session)

	err = DeleteSession(db, session.SessionID)
	assert.NoError(t, err)

	var deletedSession Session
	err = db.First(&deletedSession, "session_id = ?", session.SessionID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
