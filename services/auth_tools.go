package services

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWT Payload
type CustomClaims struct {
	Subject string `json:"sub"` // User ID
	SID     string `json:"sid"` // Session ID
	jwt.RegisteredClaims
}

// Creates a new refresh token using random bytes and encodes it in base64.
func GenerateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(tokenBytes)
	return refreshToken, nil
}

// Creates a new access token for a user with a specified expiration time.
func GenerateAccessToken(userID string, sessionID string, accessTokenExpireMinutes int16, privateKey *rsa.PrivateKey) (string, error) {
	claims := CustomClaims{
		Subject: userID,
		SID:     sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(accessTokenExpireMinutes) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Checks the validity of the provided token string using the public key.
func ValidateToken(tokenString string, publicKey *rsa.PublicKey) (*CustomClaims, error) {
	return GetTokenPayload(tokenString, publicKey, false)
}

// Parses the token string and retrieves the claims, optionally skipping validation.
func GetTokenPayload(tokenString string, publicKey *rsa.PublicKey, skipValidation bool) (*CustomClaims, error) {
	var options []jwt.ParserOption

	if skipValidation {
		options = append(options, jwt.WithoutClaimsValidation())
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	}, options...)

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func HashRefreshToken(token string) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedToken), nil
}

func CompareRefreshToken(hashedToken string, token string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedToken), []byte(token))
	return err == nil
}
