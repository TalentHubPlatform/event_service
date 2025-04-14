package middleware

import (
	"event_service/pkg/utils"
	"log/slog"
	"net/http"
)

type ExternalResponse struct {
	Status string `json:"status"`
}

func JWTAuthMiddleware(log *slog.Logger, authURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			isAuthorized, err := utils.CheckAuthorization(log, &utils.AuthRequest{AuthURL: authURL, JwtToken: authHeader})
			if err != nil || !isAuthorized {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
