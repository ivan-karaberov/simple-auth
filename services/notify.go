package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"simpleAuth/config"

	"github.com/sirupsen/logrus"
)

type NotificationPayload struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	UserIP    string `json:"user_ip"`
}

func Notify(cfg *config.Config, payload NotificationPayload) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		logrus.WithError(err).Error("Failed conversion to json")
		return
	}

	resp, err := http.Post(cfg.WebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		logrus.WithError(err).Error("Failed send notification")
		return
	}
	defer resp.Body.Close()
}
