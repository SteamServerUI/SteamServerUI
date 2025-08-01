package web

import (
	"io/fs"
	"net/http"
	"text/template"

	"github.com/SteamServerUI/SteamServerUI/v6/src/config"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

func ServeTwoBoxFormTemplate(w http.ResponseWriter, r *http.Request) {
	type Step struct {
		ID                       string
		Title                    string
		HeaderTitle              string
		StepMessage              string
		PrimaryLabel             string
		SecondaryLabel           string
		SecondaryLabelType       string
		SubmitButtonText         string
		SkipButtonText           string
		PrimaryPlaceholderText   string
		SecondaryPlaceholderText string
		ConfigField              string // Field name to save in config
		NextStep                 string // ID of the next step
	}

	type TemplateData struct {
		IsFirstTimeSetup         bool
		Path                     string
		Title                    string
		HeaderTitle              string
		StepMessage              string
		PrimaryLabel             string
		SecondaryLabel           string
		SecondaryLabelType       string
		SubmitButtonText         string
		SkipButtonText           string
		Mode                     string
		ShowExtraButtons         bool
		FooterText               string
		Step                     string
		ConfigField              string
		NextStep                 string
		PrimaryPlaceholderText   string
		SecondaryPlaceholderText string
	}

	// Get the sub-filesystem
	twoboxformAssetsFS, err := fs.Sub(config.GetTWOBOXFS(), "UIMod/onboard_bundled/twoboxform")
	if err != nil {
		logger.Web.Error("Failed to get twoboxform FS")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse the template directly from the filesystem
	tmpl, err := template.ParseFS(twoboxformAssetsFS, "twoboxform.html")
	if err != nil {
		logger.Web.Error("Failed to parse 2BoxForm template")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	path := r.URL.Path
	stepID := r.URL.Query().Get("step")
	if stepID == "" && path == "/setup" {
		stepID = "welcome" // Start with welcome page for first-time setup
	}

	// Define all steps in a map for easy access and modification
	steps := map[string]Step{
		"welcome": {
			ID:                 "welcome",
			Title:              "Steam Server UI",
			HeaderTitle:        "",
			StepMessage:        "",
			PrimaryLabel:       "",
			SecondaryLabel:     "",
			SecondaryLabelType: "hidden",
			SubmitButtonText:   "Start Setup",
			SkipButtonText:     "Skip Setup",
			NextStep:           "beta_warning",
		},
		"beta_warning": {
			ID:                 "beta_warning",
			Title:              "Beta Software",
			HeaderTitle:        "",
			StepMessage:        "This is a beta version of Steam Server UI. Please report any bugs or issues you encounter.",
			PrimaryLabel:       "",
			SecondaryLabel:     "",
			SecondaryLabelType: "hidden",
			SubmitButtonText:   "I understand",
			SkipButtonText:     "I understand in gray",
			NextStep:           "telemetry",
		},
		"telemetry": { // unused, but kept for reference
			ID:                     "telemetry",
			Title:                  "Steam Server UI",
			HeaderTitle:            "Help us improve Steam Server UI",
			StepMessage:            "by enabling anonymous telemetry. All Data is anonymous and will not be shared with any third parties.",
			PrimaryPlaceholderText: "",
			PrimaryLabel:           "",
			SecondaryLabel:         "",
			SecondaryLabelType:     "hidden",
			SubmitButtonText:       "Enable Telemetry",
			SkipButtonText:         "Skip",
			ConfigField:            "IsTelemetryEnabled",
			NextStep:               "admin_account",
		},
		"admin_account": {
			ID:                       "admin_account",
			Title:                    "Steam Server UI",
			HeaderTitle:              "Admin Account Setup",
			StepMessage:              "Set up your admin account.",
			PrimaryPlaceholderText:   "Username",
			PrimaryLabel:             "Username",
			SecondaryLabel:           "Password",
			SecondaryPlaceholderText: "Password",
			SecondaryLabelType:       "password",
			SubmitButtonText:         "Save & Continue",
			SkipButtonText:           "Skip Authentication",
			ConfigField:              "", // Special handling for admin account
			NextStep:                 "finalize",
		},
		"finalize": {
			ID:                 "finalize",
			Title:              "Finalize Setup",
			HeaderTitle:        "",
			StepMessage:        "Setup felt too short? More setup steps will be added soon™. To change settings, click the Settings tab from the main UI.",
			PrimaryLabel:       "",
			SecondaryLabel:     "",
			SecondaryLabelType: "hidden",
			SubmitButtonText:   "Return to Start",
			SkipButtonText:     "Skip Authentication",
			NextStep:           "welcome", // Return to first step if "Return to Setup" is clicked
		},
	}

	data := TemplateData{
		IsFirstTimeSetup: config.GetIsFirstTimeSetup(),
		Path:             path,
		Step:             stepID,
		FooterText:       "Need help? Check the Steam Server UI Github Wiki.",
	}

	switch {

	case path == "/login" && !config.GetAuthEnabled():
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return

	case path == "/setup":
		data.Mode = "setup"
		data.ShowExtraButtons = true

		// Find the current step in our map
		if step, exists := steps[stepID]; exists {
			data.Title = step.Title
			data.HeaderTitle = step.HeaderTitle
			data.StepMessage = step.StepMessage
			data.PrimaryLabel = step.PrimaryLabel
			data.SecondaryLabel = step.SecondaryLabel
			data.SecondaryLabelType = step.SecondaryLabelType
			data.SubmitButtonText = step.SubmitButtonText
			data.SkipButtonText = step.SkipButtonText
			data.ConfigField = step.ConfigField
			data.NextStep = step.NextStep
			data.PrimaryPlaceholderText = step.PrimaryPlaceholderText
			data.SecondaryPlaceholderText = step.SecondaryPlaceholderText
		} else {
			// Default to welcome page if step is invalid
			welcomeStep := steps["welcome"]
			data.Title = welcomeStep.Title
			data.HeaderTitle = welcomeStep.HeaderTitle
			data.StepMessage = welcomeStep.StepMessage
			data.PrimaryLabel = welcomeStep.PrimaryLabel
			data.SecondaryLabel = welcomeStep.SecondaryLabel
			data.SecondaryLabelType = welcomeStep.SecondaryLabelType
			data.SubmitButtonText = welcomeStep.SubmitButtonText
			data.SkipButtonText = welcomeStep.SkipButtonText
			data.ConfigField = welcomeStep.ConfigField
			data.NextStep = welcomeStep.NextStep
			data.Step = "welcome"
		}

	default:
		data.Title = "Steam Server UI"
		data.HeaderTitle = ""
		data.PrimaryLabel = "Username"
		data.SecondaryLabel = "Password"
		data.PrimaryPlaceholderText = "Enter Username"
		data.SecondaryPlaceholderText = "Enter Password"
		data.SecondaryLabelType = "password"
		data.SubmitButtonText = "Login"
		data.Mode = "login"
		data.ShowExtraButtons = false
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		logger.Web.Error("Failed to execute 2BoxForm template: %v" + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
