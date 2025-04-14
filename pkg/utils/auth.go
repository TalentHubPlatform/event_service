package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type AuthRequest struct {
	AuthURL  string
	JwtToken string
}

type AuthResponsePayload struct {
	Status string `json:"status" validate:"required" example:"ok"`
	Id     int    `json:"id" validate:"required" example:"52"`
}

func GetAuthorizationResponse(log *slog.Logger, request *AuthRequest) (*AuthResponsePayload, error) {
	req, err := http.NewRequest(http.MethodGet, request.AuthURL, nil)
	if err != nil {
		log.Error("Failed to create request to external service")
		return nil, nil
	}

	req.Header.Set("Authorization", request.JwtToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("External service unavailable")
		return nil, nil
	}
	defer resp.Body.Close()

	var extResp AuthResponsePayload
	if err := json.NewDecoder(resp.Body).Decode(&extResp); err != nil {
		log.Error("Failed to decode external response")
		return nil, nil
	}

	return &extResp, nil
}

func CheckAuthorization(log *slog.Logger, request *AuthRequest) (bool, error) {
	extResp, err := GetAuthorizationResponse(log, request)
	if err != nil {
		return false, err
	}

	return extResp.Status == "ok", nil
}
