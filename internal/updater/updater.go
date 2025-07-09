package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	
	"github.com/snupai/cngt-cli/internal/version"
)

const (
	githubAPIURL = "https://api.github.com/repos/snupai/cngt-cli/releases/latest"
)

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func CheckForUpdates() (*Release, bool, error) {
	resp, err := http.Get(githubAPIURL)
	if err != nil {
		return nil, false, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, false, fmt.Errorf("failed to decode release info: %w", err)
	}

	currentVersion := version.GetVersion()
	hasUpdate := release.TagName != currentVersion && release.TagName != "v"+currentVersion
	return &release, hasUpdate, nil
}

func Update() error {
	release, hasUpdate, err := CheckForUpdates()
	if err != nil {
		return err
	}

	if !hasUpdate {
		fmt.Println("Already using the latest version")
		return nil
	}

	fmt.Printf("Updating from %s to %s...\n", version.GetVersion(), release.TagName)

	assetName := getBinaryName()
	var downloadURL string
	
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no suitable binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	tmpFile, err := downloadBinary(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer os.Remove(tmpFile)

	if err := replaceBinary(tmpFile); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	fmt.Printf("Successfully updated to %s\n", release.TagName)
	return nil
}

func getBinaryName() string {
	switch runtime.GOOS {
	case "windows":
		return "cngt-cli-windows-amd64.exe"
	case "darwin":
		return "cngt-cli-darwin-amd64"
	case "linux":
		return "cngt-cli-linux-amd64"
	default:
		return "cngt-cli-linux-amd64"
	}
}

func downloadBinary(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpFile, err := os.CreateTemp("", "cngt-cli-update-*")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		os.Remove(tmpFile.Name())
		return "", err
	}

	return tmpFile.Name(), nil
}

func replaceBinary(newBinary string) error {
	currentBinary, err := os.Executable()
	if err != nil {
		return err
	}

	backupBinary := currentBinary + ".backup"
	if err := os.Rename(currentBinary, backupBinary); err != nil {
		return err
	}

	if err := os.Rename(newBinary, currentBinary); err != nil {
		os.Rename(backupBinary, currentBinary)
		return err
	}

	os.Remove(backupBinary)
	return nil
}

func SelfUpdate() error {
	fmt.Println("Checking for CLI updates...")
	
	release, hasUpdate, err := CheckForUpdates()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !hasUpdate {
		fmt.Println("CLI is up to date")
		return nil
	}

	fmt.Printf("New version available: %s\n", release.TagName)
	fmt.Print("Would you like to update? (y/N): ")
	
	var response string
	fmt.Scanln(&response)
	
	if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
		return Update()
	}

	return nil
}