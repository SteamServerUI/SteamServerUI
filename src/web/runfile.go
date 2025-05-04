package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v6/src/argmgr"
	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

// APIGameArg is a DTO for GameArg, including RuntimeValue and all fields
type APIGameArg struct {
	Flag          string `json:"flag"`
	DefaultValue  string `json:"default"`
	RuntimeValue  string `json:"runtime_value"`
	Required      bool   `json:"required"`
	RequiresValue bool   `json:"requires_value"`
	Description   string `json:"description"`
	Type          string `json:"type"`
	Special       string `json:"special,omitempty"`
	UILabel       string `json:"ui_label"`
	UIGroup       string `json:"ui_group"`
	Weight        int    `json:"weight"`
	Min           int    `json:"min,omitempty"`
	Max           int    `json:"max,omitempty"`
	Disabled      bool   `json:"disabled"`
}

// APIMeta mirrors argmgr.Meta for API responses
type APIMeta struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// APIRunFile is a DTO for RunFile, using APIGameArg and APIMeta
type APIRunFile struct {
	Meta              APIMeta                 `json:"meta"`
	Architecture      string                  `json:"architecture,omitempty"`
	SteamAppID        string                  `json:"steam_app_id"`
	WindowsExecutable string                  `json:"windows_executable"`
	LinuxExecutable   string                  `json:"linux_executable"`
	Args              map[string][]APIGameArg `json:"args"`
}

// apiResponse is the standard JSON response format
type apiResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error,omitempty"`
}

// toAPIGameArg converts argmgr.GameArg to APIGameArg
func toAPIGameArg(arg argmgr.GameArg) APIGameArg {
	return APIGameArg{
		Flag:          arg.Flag,
		DefaultValue:  arg.DefaultValue,
		RuntimeValue:  arg.RuntimeValue,
		Required:      arg.Required,
		RequiresValue: arg.RequiresValue,
		Description:   arg.Description,
		Type:          arg.Type,
		Special:       arg.Special,
		UILabel:       arg.UILabel,
		UIGroup:       arg.UIGroup,
		Weight:        arg.Weight,
		Min:           arg.Min,
		Max:           arg.Max,
		Disabled:      arg.Disabled,
	}
}

// toAPIRunFile converts argmgr.RunFile to APIRunFile
func toAPIRunFile(rf *argmgr.RunFile) APIRunFile {
	apiArgs := make(map[string][]APIGameArg)
	for category, args := range rf.Args {
		for _, arg := range args {
			apiArgs[category] = append(apiArgs[category], toAPIGameArg(arg))
		}
	}
	return APIRunFile{
		Meta: APIMeta{
			Name:    rf.Meta.Name,
			Version: rf.Meta.Version,
		},
		Architecture:      rf.Architecture,
		SteamAppID:        rf.SteamAppID,
		WindowsExecutable: rf.WindowsExecutable,
		LinuxExecutable:   rf.LinuxExecutable,
		Args:              apiArgs,
	}
}

// writeJSONResponse writes a JSON response with the given status code
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}, errMsg string) {
	resp := apiResponse{Data: data, Error: errMsg}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to write JSON response: %v", err))
	}
}

// HandleRunfileGroups handles GET /api/v2/runfile/groups
func HandleRunfileGroups(w http.ResponseWriter, r *http.Request) {
	logger.Runfile.Debug("GET /api/v2/runfile/groups")
	if argmgr.CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		writeJSONResponse(w, http.StatusInternalServerError, nil, "runfile not loaded")
		return
	}

	groups := argmgr.GetUIGroups()
	logger.Runfile.Info("fetched UI groups")
	writeJSONResponse(w, http.StatusOK, groups, "")
}

// HandleRunfileArgs handles GET /api/v2/runfile/args
func HandleRunfileArgs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	logger.Runfile.Debug(fmt.Sprintf("GET /api/v2/runfile/args group=%s", group))

	if argmgr.CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		writeJSONResponse(w, http.StatusInternalServerError, nil, "runfile not loaded")
		return
	}

	var args []argmgr.GameArg
	if group != "" {
		// Validate group
		validGroups := argmgr.GetUIGroups()
		valid := false
		for _, g := range validGroups {
			if g == group {
				valid = true
				break
			}
		}
		if !valid {
			logger.Runfile.Error(fmt.Sprintf("invalid group: %s", group))
			writeJSONResponse(w, http.StatusBadRequest, nil, fmt.Sprintf("invalid group: %s", group))
			return
		}
		args = argmgr.GetArgsByGroup(group)
	} else {
		args = argmgr.GetAllArgs()
	}

	// Convert to APIGameArg
	apiArgs := make([]APIGameArg, len(args))
	for i, arg := range args {
		apiArgs[i] = toAPIGameArg(arg)
	}

	logger.Runfile.Info(fmt.Sprintf("fetched args for group=%s", group))
	writeJSONResponse(w, http.StatusOK, apiArgs, "")
}

// HandleRunfileArgUpdate handles POST /api/v2/runfile/args/update
func HandleRunfileArgUpdate(w http.ResponseWriter, r *http.Request) {
	logger.Runfile.Debug("POST /api/v2/runfile/args")

	if argmgr.CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		writeJSONResponse(w, http.StatusInternalServerError, nil, "runfile not loaded")
		return
	}

	var req struct {
		Flag  string `json:"flag"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Runfile.Error(fmt.Sprintf("invalid request body: %v", err))
		writeJSONResponse(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	if req.Flag == "" {
		logger.Runfile.Error("flag is required")
		writeJSONResponse(w, http.StatusBadRequest, nil, "flag is required")
		return
	}

	if err := argmgr.SetArgValue(req.Flag, req.Value); err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to set arg %s: %v", req.Flag, err))
		writeJSONResponse(w, http.StatusBadRequest, nil, fmt.Sprintf("failed to set arg: %v", err))
		return
	}

	logger.Runfile.Info(fmt.Sprintf("updated arg %s to %s", req.Flag, req.Value))
	writeJSONResponse(w, http.StatusOK, map[string]string{"flag": req.Flag, "value": req.Value}, "")
}

// HandleRunfile handles GET /api/v2/runfile
func HandleRunfile(w http.ResponseWriter, r *http.Request) {
	logger.Runfile.Debug("GET /api/v2/runfile")

	if argmgr.CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		writeJSONResponse(w, http.StatusInternalServerError, nil, "runfile not loaded")
		return
	}

	apiRunfile := toAPIRunFile(argmgr.CurrentRunfile)
	logger.Runfile.Info("fetched runfile")
	writeJSONResponse(w, http.StatusOK, apiRunfile, "")
}

// HandleRunfileSave handles POST /api/v2/runfile/save
func HandleRunfileSave(w http.ResponseWriter, r *http.Request) {
	logger.Runfile.Debug("POST /api/v2/runfile/save")

	if argmgr.CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		writeJSONResponse(w, http.StatusInternalServerError, nil, "runfile not loaded")
		return
	}

	if err := argmgr.SaveRunfile(); err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to save runfile: %v", err))
		writeJSONResponse(w, http.StatusInternalServerError, nil, fmt.Sprintf("failed to save runfile: %v", err))
		return
	}

	logger.Runfile.Info("runfile saved")
	writeJSONResponse(w, http.StatusOK, "runfile saved", "")
}
