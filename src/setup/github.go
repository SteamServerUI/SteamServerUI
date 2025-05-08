package setup

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v6/src/logger"
)

// downloadAllFiles downloads all files in the provided map concurrently
func downloadAllFiles(files map[string]string) {
	var wg sync.WaitGroup
	for filepath, url := range files {
		wg.Add(1)
		go func(filepath, url string) {
			defer wg.Done()
			fileName := filepath[strings.LastIndex(filepath, "/")+1:]
			logger.Install.Info("Downloading " + fileName + "...")
			err := downloadFile(filepath, url)
			if err != nil {
				logger.Install.Error("❌Error downloading " + fileName + ": " + err.Error())
			} else {
				logger.Install.Info("✅Downloaded " + fileName + " successfully from branch " + downloadBranch)
			}
		}(filepath, url)
	}
	wg.Wait()
	logger.Install.Info("✅All files downloaded successfully.")
}

// updateFilesIfDifferent checks for differences and updates files if necessary using concurrency
func updateFilesIfDifferent(files map[string]string) {
	var wg sync.WaitGroup
	for filepath, url := range files {
		fileName := filepath[strings.LastIndex(filepath, "/")+1:]
		if fileName == "config.json" {
			continue // Skip updating config.json to preserve local changes
		}
		wg.Add(1)
		go func(filepath, url string) {
			defer wg.Done()
			checkAndUpdateFile(filepath, url)
		}(filepath, url)
	}
	wg.Wait()
	logger.Install.Info("✅File validation and update check complete.")
}

// checkAndUpdateFile checks if a file needs updating by comparing its SHA-1 hash with the remote ETag
func checkAndUpdateFile(filepath, url string) {
	fileName := filepath[strings.LastIndex(filepath, "/")+1:]

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// File doesn't exist locally, download it
		logger.Install.Info("Downloading " + fileName + "...")
		err := downloadFile(filepath, url)
		if err != nil {
			logger.Install.Error("❌Error downloading " + fileName + ": " + err.Error())
		} else {
			logger.Install.Debug("✅Downloaded " + fileName + " successfully")
		}
	} else {
		// File exists, check if it needs updating
		localHash, err := computeGitBlobSHA1(filepath)
		if err != nil {
			logger.Install.Error("❌Error computing hash for " + fileName + ": " + err.Error())
			return
		}

		// Extract the necessary parts from URL to build the API call
		// Example URL: https://raw.githubusercontent.com/SteamServerUI/SteamServerUI/main/UIMod/index.html
		urlParts := strings.Split(url, "/")
		if len(urlParts) < 7 {
			logger.Install.Error("❌Invalid URL format: " + url)
			return
		}

		repoOwner := urlParts[3]
		repoName := urlParts[4]
		branch := urlParts[5]
		filePath := strings.Join(urlParts[6:], "/")

		remoteHash, err := getFileHash(repoOwner, repoName, branch, filePath)
		if err != nil {
			return
		}

		logger.Install.Debug("Local hash for " + fileName + ": " + localHash)
		logger.Install.Debug("Remote hash for " + fileName + ": " + remoteHash)

		if localHash != remoteHash {
			logger.Install.Debug("🔄Updating " + fileName + " due to differences...")
			err := downloadFile(filepath, url)
			if err != nil {
				logger.Install.Error("❌Error updating " + fileName + ": " + err.Error())
			} else {
				logger.Install.Info("✅Updated " + fileName + " successfully from branch " + downloadBranch)
			}
		} else {
			logger.Install.Debug("✅" + fileName + " is up-to-date.")
		}
	}
}

func computeGitBlobSHA1(filepath string) (string, error) {
	// Read the file content
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	// Create the blob header (format: "blob " + content length + null byte)
	header := fmt.Sprintf("blob %d\x00", len(content))

	// Combine header and content
	blobData := append([]byte(header), content...)

	// Compute SHA-1 hash
	hash := sha1.New()
	hash.Write(blobData)

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// getFileHash fetches the file hash from GitHub API
func getFileHash(repoOwner, repoName, branch, filePath string) (string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s",
		repoOwner, repoName, filePath, branch)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	// Optional: req.Header.Add("Authorization", "token YOUR_GITHUB_TOKEN")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if err := handleRateLimitErrors(resp); err != nil {
			return "", err
		}
		logger.Install.Error(fmt.Sprintf("❌Request failed with status: %s", resp.Status))
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	var fileInfo struct {
		SHA string `json:"sha"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&fileInfo); err != nil {
		return "", err
	}

	return fileInfo.SHA, nil
}

// handleRateLimitErrors processes GitHub API rate limit errors
func handleRateLimitErrors(resp *http.Response) error {
	if resp.StatusCode != http.StatusForbidden {
		return nil
	}

	// Check secondary rate limit
	if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
		return handleSecondaryRateLimit(retryAfter, resp.Status)
	}

	// Check primary rate limit
	if remaining := resp.Header.Get("x-ratelimit-remaining"); remaining == "0" {
		return handlePrimaryRateLimit(resp.Header.Get("x-ratelimit-reset"), resp.Status)
	}

	// Generic 403 error
	logger.Install.Error("❌Forbidden request, no specific rate limit info available")
	return fmt.Errorf("bad status: %s", resp.Status)
}

// handleSecondaryRateLimit processes secondary rate limit exceeded errors
func handleSecondaryRateLimit(retryAfter, status string) error {
	waitSeconds, err := strconv.Atoi(retryAfter)
	if err != nil {
		logger.Install.Error(fmt.Sprintf("❌Failed to parse Retry-After header: %v", err))
		return fmt.Errorf("bad status: %s, failed to parse Retry-After: %v", status, err)
	}
	errMsg := fmt.Sprintf("Github API secondary rate limit exceeded. (too many requests in a short time). Retry after %d seconds", waitSeconds)
	logger.Install.Error("❌" + errMsg)
	return fmt.Errorf("bad status: %s, %s", status, errMsg)
}

// handlePrimaryRateLimit processes primary rate limit exceeded errors
func handlePrimaryRateLimit(reset, status string) error {
	resetInt, err := strconv.ParseInt(reset, 10, 64)
	if err != nil {
		logger.Install.Error(fmt.Sprintf("❌Failed to parse x-ratelimit-reset header: %v", err))
		return fmt.Errorf("bad status: %s, failed to parse x-ratelimit-reset: %v", status, err)
	}
	resetTime := time.Unix(resetInt, 0).UTC()
	errMsg := fmt.Sprintf("Update check Github ratelimit exceeded: hourly request quota of 60 calls reached. Resets at %s", resetTime.Format(time.RFC1123))
	logger.Install.Debug("🧱" + errMsg)
	return fmt.Errorf("bad status: %s, %s", status, errMsg)
}
