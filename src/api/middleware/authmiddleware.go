package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

var setupReminderCount = 0 // to limit the number of setup reminders shown to the user

// AuthMiddleware protects routes with cookie-based JWT
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request details for debugging
		//logger.Web.Debug("Request Path:" + r.URL.Path) //very spammy

		// Check for first-time setup redirect
		if config.GetIsFirstTimeSetup() {
			totalSetupReminderCount := 3 // Defines how often we redirect the users reqests to the setup page
			if setupReminderCount < totalSetupReminderCount {
				if r.URL.Path == "/" && (r.Referer() == "" || r.Referer() != "/setup") {
					remainingReminderCount := totalSetupReminderCount - setupReminderCount
					logger.Web.Warn("🔍Redirecting to setup page, you should really enable authentication...")
					logger.Web.Warn(fmt.Sprintf("You will be remined %s more times.", strconv.Itoa(remainingReminderCount)))
					http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
					setupReminderCount++
					return
				}
			}
		}

		if !config.GetAuthEnabled() {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("AuthToken")
		if err != nil {
			// Browser redirect check
			accept := r.Header.Get("Accept")
			if accept != "" && strings.Contains(accept, "text/html") {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			// API response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized - No token"})
			return
		}

		valid, err := security.ValidateJWT(cookie.Value)
		if err != nil || !valid {
			// Browser redirect check
			accept := r.Header.Get("Accept")
			if accept != "" && strings.Contains(accept, "text/html") {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			// API response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized - Invalid token"})
			logger.Security.Warn("Unauthorized Request - Invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
