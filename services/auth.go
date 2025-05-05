package services

import (
	"fmt"
	"simpleAuth/config"
	"simpleAuth/models"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserInfo struct {
	UserID    string
	UserIP    string
	UserAgent string
}

type TokenPair struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
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

func RefreshToken(db *gorm.DB, cfg *config.Config, tokens *TokenPair, userIP string, userAgent string) (*TokenPair, error) {
	payload, err := GetTokenPayload(tokens.AccessToken, cfg.RSAPublicKey, true)
	if err != nil {
		logrus.WithError(err).Error("Failed get payload from Access token")
		return nil, err
	}

	session, err := models.GetSession(db, payload.SID)
	if err != nil {
		logrus.WithError(err).Errorf("Failed get session by session ID %s", payload.SID)
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	if !CompareRefreshToken(session.RefreshToken, tokens.RefreshToken) {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if time.Now().After(session.ExpireAt) {
		return nil, fmt.Errorf("token has expired")
	}

	if session.UserAgent != userAgent {
		models.DeleteSession(db, session.SessionID)
		return nil, fmt.Errorf("user agent not equal")
	}

	if session.IP != userIP {
		notificationPayload := NotificationPayload{
			UserID:    session.UserID,
			SessionID: session.SessionID,
			UserIP:    userIP,
		}

		Notify(cfg, notificationPayload)
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		logrus.WithError(err).Error("Failed generate refresh token")
		return nil, fmt.Errorf("failed generate refresh token")
	}

	accessToken, err := GenerateAccessToken(session.UserID, session.SessionID, cfg.AccessTokenExpireMinutes, cfg.RSAPrivateKey)
	if err != nil {
		logrus.WithError(err).Error("Failed generate access token")
		return nil, fmt.Errorf("failed generate access token")
	}

	session.ExpireAt = time.Now().Add(time.Duration(cfg.RefreshTokenExpireMinutes) * time.Minute)
	session.RefreshToken = refreshToken
	models.UpdateSession(db, session)

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
