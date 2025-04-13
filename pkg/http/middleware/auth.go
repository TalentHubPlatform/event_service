package middleware

import (
	"encoding/json"
	"net/http"
)

type ExternalResponse struct {
	Status string `json:"status"`
}

func JWTAuthMiddleware(authURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, authURL, nil)
			if err != nil {
				http.Error(w, "Failed to create request to external service", http.StatusInternalServerError)
				return
			}

			req.Header.Set("Authorization", authHeader)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				http.Error(w, "External service unavailable", http.StatusServiceUnavailable)
				return
			}
			defer resp.Body.Close()

			var extResp ExternalResponse
			if err := json.NewDecoder(resp.Body).Decode(&extResp); err != nil {
				http.Error(w, "Failed to decode external response", http.StatusInternalServerError)
				return
			}

			if extResp.Status == "" {
				http.Error(w, "Access denied by external service", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
