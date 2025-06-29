// handlers.go
package web

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/loader"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v6/src/security"
)

var setupReminderCount = 0 // to limit the number of setup reminders shown to the user

// LoginHandler issues a JWT cookie
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for cross-origin requests
	setCORSHeaders(w, r)

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var creds security.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - Invalid JSON"})
		return
	}

	// Check credentials using security package
	valid, err := security.ValidateCredentials(creds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}
	if !valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized - Invalid credentials"})
		return
	}

	// Generate JWT
	tokenString, err := security.GenerateJWT(creds.Username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	// Set cookie - modify to work with cross-origin requests
	http.SetCookie(w, &http.Cookie{
		Name:     "AuthToken",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Duration(config.GetAuthTokenLifetime()) * time.Minute),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode, // Change to None to allow cross-origin
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// AuthMiddleware protects routes with cookie-based JWT and adds CORS headers
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for cross-origin requests
		setCORSHeaders(w, r)

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Existing authentication logic
		if !config.GetAuthEnabled() {
			if config.GetIsFirstTimeSetup() {
				if setupReminderCount < 1 {
					http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
					setupReminderCount++
				}
			}
			next.ServeHTTP(w, r)
			return
		}

		// Check for token from multiple sources
		var tokenString string

		// 1. Check cookie first
		cookie, err := r.Cookie("AuthToken")
		if err == nil {
			tokenString = cookie.Value
		}

		// 2. If no cookie, check Authorization header
		if tokenString == "" {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// 3. Finally check query parameter (for EventSource)
		if tokenString == "" {
			tokenString = r.URL.Query().Get("token")
		}

		// No token found in any location
		if tokenString == "" {
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

		// Validate token
		valid, err := security.ValidateJWT(tokenString)
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

// Helper function to set CORS headers consistently across handlers
func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for cross-origin requests
	setCORSHeaders(w, r)

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Clear the cookie by setting it with an expired time
	http.SetCookie(w, &http.Cookie{
		Name:     "AuthToken",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Set to past time to expire immediately
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteNoneMode, // Changed to None for cross-origin
	})

	accept := r.Header.Get("Accept")
	if accept != "" && strings.Contains(accept, "text/html") {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	// For API requests, return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully logged out",
	})
}

// Add a special handler for auth check
func AuthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for cross-origin requests
	setCORSHeaders(w, r)

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// The AuthMiddleware will have already verified the token
	// If we get here, the user is authenticated

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "authenticated",
	})
}

// OLD

// RegisterUserHandler registers new users
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {

	setCORSHeaders(w, r)

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var creds security.UserCredentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - Invalid JSON"})
		return
	}

	// Hash the password
	hashedPassword, err := security.HashPassword(creds.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}

	// Initialize Users map if nil
	if config.GetUsers() == nil {
		config.SetUsers(make(map[string]string))
	}

	// Add or update the user
	config.SetUsers(map[string]string{creds.Username: hashedPassword})
	config.SetUserLevels(map[string]string{creds.Username: config.GetDefaultUserLevel()}) // TODO: remove default user level

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "User registered successfully",
		"username": creds.Username,
	})
}

// SetupFinalizeHandler marks setup as complete
func SetupFinalizeHandler(w http.ResponseWriter, r *http.Request) {

	setCORSHeaders(w, r)

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//check if users map is nil or empty
	if config.GetUsers() == nil || len(config.GetUsers()) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - No users registered - cannot finalize setup at this time"})
		return
	}

	// Mark setup as complete and enable auth
	config.SetIsFirstTimeSetup(false)
	config.SetAuthEnabled(true)

	logger.Web.Core("User Setup finalized successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":      "Setup finalized successfully",
		"restart_hint": "You will be redirected to the login page...",
	})
	loader.ReloadConfig()
}
