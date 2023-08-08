package api

import (
	"net/http"
	"strings"
)

// MattermostAuthorizationRequired checks if request is authorized.
func MattermostAuthorizationRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Mattermost-User-Id")
		if isNoAuthorizationRequest(r) || userID != "" {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Not authorized", http.StatusUnauthorized)
	})
}

func isNoAuthorizationRequest(r *http.Request) bool {
	return strings.Contains(r.URL.Path, ConfigBasePath) ||
		strings.Contains(r.URL.Path, TelemetryBasePath)
}
