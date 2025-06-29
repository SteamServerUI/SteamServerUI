package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/golang-jwt/jwt/v5"
)

// accessLevelAdminMiddleware restricts access to routes to users with "admin" role.
func accessLevelMiddleware(next http.HandlerFunc, requiredLevel ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if len(requiredLevel) == 0 {
			requiredLevel = []string{"superadmin"}
		}

		// Set CORS headers to match AuthMiddleware
		setCORSHeaders(w, r)

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// If auth is disabled, allow all requests through
		if !config.GetAuthEnabled() {
			next(w, r)
			return
		}

		// Extract token from cookie, header, or query (same as AuthMiddleware)
		var tokenString string
		cookie, err := r.Cookie("AuthToken")
		if err == nil {
			tokenString = cookie.Value
		}
		if tokenString == "" {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}
		if tokenString == "" {
			tokenString = r.URL.Query().Get("token")
		}
		if tokenString == "" {
			returnForbiddenOrRedirect(w, r, "No token provided")
			return
		}

		// Parse JWT to get username
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetJwtKey()), nil
		})
		if err != nil || !token.Valid {
			returnForbiddenOrRedirect(w, r, "Invalid token")
			return
		}
		username, ok := (*claims)["id"].(string)
		if !ok || username == "" {
			returnForbiddenOrRedirect(w, r, "No username in token")
			return
		}

		// Check if user has required Level
		level := config.GetUserLevel(username)
		hasRequiredLevel := slices.Contains(requiredLevel, level)
		if !hasRequiredLevel {
			var formattedLevels string
			if len(requiredLevel) == 1 {
				formattedLevels = requiredLevel[0]
			} else if len(requiredLevel) == 2 {
				formattedLevels = requiredLevel[0] + " or " + requiredLevel[1]
			} else {
				formattedLevels = strings.Join(requiredLevel[:len(requiredLevel)-1], ", ") + ", or " + requiredLevel[len(requiredLevel)-1]
			}
			returnForbiddenOrRedirect(w, r, fmt.Sprintf("Forbidden - %s access level required", formattedLevels))
			return
		}

		// User is authorized for this route
		next(w, r)
	}
}

// returnForbiddenOrRedirect sends a 403 Forbidden for API requests or redirects to /login for browser requests.
func returnForbiddenOrRedirect(w http.ResponseWriter, r *http.Request, message string) {
	accept := r.Header.Get("Accept")
	if accept != "" && strings.Contains(accept, "text/html") {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
