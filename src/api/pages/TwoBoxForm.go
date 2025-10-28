package pages

import (
	"io/fs"
	"net/http"
	"text/template"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/localization"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
)

func ServeTwoBoxFormTemplate(w http.ResponseWriter, r *http.Request) {
	type Step struct {
		ID                 string
		Title              string
		HeaderTitle        string
		StepMessage        string
		PrimaryLabel       string
		SecondaryLabel     string
		SecondaryLabelType string
		SecondaryOptions   []struct {
			Display string // Text shown in dropdown
			Value   string // Value sent on submission
		}
		SubmitButtonText         string
		SkipButtonText           string
		PrimaryPlaceholderText   string
		SecondaryPlaceholderText string
		ConfigField              string // Field name to save in config
		NextStep                 string // ID of the next step
	}

	type TemplateData struct {
		IsFirstTimeSetup   bool
		Path               string
		Title              string
		HeaderTitle        string
		StepMessage        string
		PrimaryLabel       string
		SecondaryLabel     string
		SecondaryLabelType string
		SecondaryOptions   []struct {
			Display string
			Value   string
		}
		SubmitButtonText         string
		SkipButtonText           string
		Mode                     string
		ShowExtraButtons         bool
		FooterText               string
		FooterTextInfo           string
		Step                     string
		ConfigField              string
		NextStep                 string
		PrimaryPlaceholderText   string
		SecondaryPlaceholderText string
		Steps                    []Step
	}

	twoboxformAssetsFS, err := fs.Sub(config.GetV1UIFS(), "SSUI/onboard_bundled/twoboxform")
	if err != nil {
		logger.Web.Error("Failed to get bundled FS")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

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
			Title:              localization.GetString("UIText_Welcome_Title"),
			HeaderTitle:        localization.GetString("UIText_Welcome_HeaderTitle"),
			StepMessage:        "",
			PrimaryLabel:       "",
			SecondaryLabel:     "",
			SecondaryLabelType: "hidden",
			SubmitButtonText:   localization.GetString("UIText_Welcome_SubmitButton"),
			SkipButtonText:     localization.GetString("UIText_Welcome_SkipButton"),
			NextStep:           "pls_read",
		},
		"pls_read": {
			ID:                 "pls_read",
			Title:              localization.GetString("UIText_PlsRead_Title"),
			HeaderTitle:        localization.GetString("UIText_PlsRead_HeaderTitle"),
			StepMessage:        localization.GetString("UIText_PlsRead_StepMessage"),
			PrimaryLabel:       "",
			SecondaryLabel:     "",
			SecondaryLabelType: "hidden",
			SubmitButtonText:   localization.GetString("UIText_PlsRead_SubmitButton"),
			SkipButtonText:     localization.GetString("UIText_PlsRead_SkipButton"),
			NextStep:           "admin_account",
		},
		"sample_step": {
			ID:                       "sample_step",
			Title:                    localization.GetString("UIText_SampleStep_Title"),
			HeaderTitle:              localization.GetString("UIText_SampleStep_HeaderTitle"),
			StepMessage:              localization.GetString("UIText_SampleStep_StepMessage"),
			PrimaryPlaceholderText:   localization.GetString("UIText_SampleStep_PrimaryPlaceholder"),
			PrimaryLabel:             localization.GetString("UIText_SampleStep_PrimaryLabel"),
			SecondaryPlaceholderText: localization.GetString("UIText_SampleStep_SecondaryPlaceholder"),
			SecondaryLabel:           "",
			SecondaryLabelType:       "hidden",
			SubmitButtonText:         localization.GetString("UIText_SampleStep_SubmitButton"),
			SkipButtonText:           localization.GetString("UIText_SampleStep_SkipButton"),
			ConfigField:              "SampleStep",
			NextStep:                 "admin_account",
		},
		"admin_account": {
			ID:                       "admin_account",
			Title:                    localization.GetString("UIText_AdminAccount_Title"),
			HeaderTitle:              localization.GetString("UIText_AdminAccount_HeaderTitle"),
			StepMessage:              localization.GetString("UIText_AdminAccount_StepMessage"),
			PrimaryPlaceholderText:   localization.GetString("UIText_AdminAccount_PrimaryPlaceholder"),
			PrimaryLabel:             localization.GetString("UIText_AdminAccount_PrimaryLabel"),
			SecondaryLabel:           localization.GetString("UIText_AdminAccount_SecondaryLabel"),
			SecondaryPlaceholderText: localization.GetString("UIText_AdminAccount_SecondaryPlaceholder"),
			SecondaryLabelType:       "password",
			SubmitButtonText:         localization.GetString("UIText_AdminAccount_SubmitButton"),
			SkipButtonText:           localization.GetString("UIText_AdminAccount_SkipButton"),
			ConfigField:              "",
			NextStep:                 "finalize",
		},
		"finalize": {
			ID:               "finalize",
			Title:            localization.GetString("UIText_Finalize_Title"),
			HeaderTitle:      localization.GetString("UIText_Finalize_HeaderTitle"),
			StepMessage:      localization.GetString("UIText_Finalize_StepMessage"),
			PrimaryLabel:     "",
			SecondaryLabel:   "",
			SubmitButtonText: localization.GetString("UIText_Finalize_SubmitButton"),
			SkipButtonText:   localization.GetString("UIText_Finalize_SkipButton"),
			NextStep:         "welcome", // Return to first step if "Return to Setup" is clicked
		},
	}

	data := TemplateData{
		IsFirstTimeSetup: config.GetIsFirstTimeSetup(),
		Path:             path,
		Step:             stepID,
		FooterText:       localization.GetString("UIText_FooterText"),
		FooterTextInfo:   localization.GetString("UIText_FooterTextInfo"),
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
			data.SecondaryOptions = step.SecondaryOptions
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
		stepOrder := []string{
			"welcome", "pls_read",
			//"game_branch",
			//"server_name", "save_name", "world_id", "max_players",
			//"server_password", "discord_enabled", "discord_token", "control_panel_channel", "save_channel",
			//"log_channel", "connection_list_channel", "status_channel", "control_channel",
			//"network_config_choice", "game_port", "update_port", "upnp_enabled",
			//"local_ip_address",
			"admin_account", "finalize",
		}
		var stepSlice []Step
		for _, id := range stepOrder {
			if step, exists := steps[id]; exists {
				stepSlice = append(stepSlice, step)
			}
		}
		data.Steps = stepSlice

	case path == "/changeuser":
		data.Title = localization.GetString("UIText_ChangeUser_Title")
		data.HeaderTitle = localization.GetString("UIText_ChangeUser_HeaderTitle")
		data.PrimaryLabel = localization.GetString("UIText_ChangeUser_PrimaryLabel")
		data.SecondaryLabel = localization.GetString("UIText_ChangeUser_SecondaryLabel")
		data.SecondaryPlaceholderText = localization.GetString("UIText_ChangeUser_SecondaryPlaceholder")
		data.SecondaryLabelType = "password"
		data.SubmitButtonText = localization.GetString("UIText_ChangeUser_SubmitButton")
		data.Mode = "changeuser"
		data.ShowExtraButtons = false

	default:
		data.Title = localization.GetString("UIText_Login_Title")
		data.HeaderTitle = localization.GetString("UIText_Login_HeaderTitle") + config.GetBackendName()
		data.PrimaryLabel = localization.GetString("UIText_Login_PrimaryLabel")
		data.SecondaryLabel = localization.GetString("UIText_Login_SecondaryLabel")
		data.PrimaryPlaceholderText = localization.GetString("UIText_Login_PrimaryPlaceholder")
		data.SecondaryPlaceholderText = localization.GetString("UIText_Login_SecondaryPlaceholder")
		data.SecondaryLabelType = "password"
		data.SubmitButtonText = localization.GetString("UIText_Login_SubmitButton")
		data.Mode = "login"
		data.Step = ""
		data.ShowExtraButtons = false
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, data); err != nil {
		logger.Web.Error("Failed to execute 2BoxForm template: %v" + err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
