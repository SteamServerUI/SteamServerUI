package httpauth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/middleware"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/google/uuid"
)

func RegisterAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	middleware.HandleCORS(w, r)

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Allow only GET or POST methods
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method Not Allowed"})
		return
	}

	// Set default duration for GET requests, require duration for POST
	durationMonths := 1
	if r.Method == http.MethodPost {
		var reqBody struct {
			DurationMonths *int `json:"durationMonths"` // Use pointer to distinguish between 0 and unspecified
		}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - Invalid JSON"})
			return
		}
		if reqBody.DurationMonths == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - durationMonths is required for POST"})
			return
		}
		if *reqBody.DurationMonths <= 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - Duration must be positive"})
			return
		}
		if *reqBody.DurationMonths > 120 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Bad Request - Duration must be less than 10 years"})
			return
		}
		durationMonths = *reqBody.DurationMonths
	}

	var creds security.UserCredentials

	// Generate a random UUID as the username
	creds.Username = "apikey-" + uuid.NewString()

	// Hash a random UUID as the password
	hashedPassword, err := security.HashPassword(uuid.NewString())
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	apikey, err := security.GenerateJWT(creds.Username, durationMonths)
	expires := time.Now().AddDate(0, durationMonths, 0)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal Server Error"})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "APIKey registered successfully",
		"apikey":  apikey,
		"expires": expires.Format(time.RFC3339),
	})
	logger.Security.Infof("APIKey %s registered successfully. Expires: %s ", creds.Username, expires.Format(time.RFC3339))
}
