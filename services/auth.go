package services

import (
	"simpleAuth/config"
	"simpleAuth/models"
	"time"

	"gorm.io/gorm"
)

type UserInfo struct {
	UserID    string
	UserIP    string
	UserAgent string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func SignIn(db *gorm.DB, cfg *config.Config, userDetail UserInfo) (*TokenPair, error) {
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	hashedRefreshToken, err := HashRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	session := models.Session{
		UserID:       userDetail.UserID,
		IP:           userDetail.UserIP,
		UserAgent:    userDetail.UserAgent,
		RefreshToken: hashedRefreshToken,
		ExpireAt:     time.Now().Add(time.Duration(cfg.RefreshTokenExpireMinutes) * time.Minute),
	}

	sessionID, err := models.CreateSession(db, &session)
	if err != nil {
		return nil, err
	}

	accessToken, err := GenerateAccessToken(userDetail.UserID, sessionID, cfg.AccessTokenExpireMinutes, cfg.RSAPrivateKey)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func SignOut(db *gorm.DB, sessionID string) error {
	return models.DeleteSession(db, sessionID)
}

func CheckSessionExists(db *gorm.DB, sessionID string) (bool, error) {
	_, err := models.GetSession(db, sessionID)
	if err != nil {
		return false, err
	}
	return true, nil
}
