package cngt

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/snupai/cngt-cli/internal/config"
	"github.com/snupai/cngt-cli/internal/deps"
)

const (
	repoURL = "https://github.com/SebiAi/custom-nothing-glyph-tools.git"
)

type Status struct {
	RepoStatus    string
	PythonStatus  string
	DepsStatus    string
}

func IsInstalled(path string) bool {
	_, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil
}

func Install(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	return nil
}

func Update() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	repo, err := git.PlainOpen(cfg.CNGTPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to pull updates: %w", err)
	}

	return nil
}

func RunScript(scriptName string, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	scriptPath := filepath.Join(cfg.CNGTPath, scriptName)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script %s not found", scriptName)
	}

	// Find the best Python command
	pythonCmd := findPythonCommand()
	if pythonCmd == "" {
		return fmt.Errorf("Python is not installed or not found in PATH")
	}
	
	cmdArgs := append([]string{scriptPath}, args...)
	cmd := exec.Command(pythonCmd, cmdArgs...)
	cmd.Dir = cfg.CNGTPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func findPythonCommand() string {
	pythonCommands := []string{"python", "python3", "py"}
	
	for _, pythonCmd := range pythonCommands {
		cmd := exec.Command(pythonCmd, "--version")
		if cmd.Run() == nil {
			return pythonCmd
		}
	}
	return ""
}

func GetStatus() Status {
	cfg, err := config.Load()
	if err != nil {
		return Status{
			RepoStatus:   "Error loading config",
			PythonStatus: "Unknown",
			DepsStatus:   "Unknown",
		}
	}

	status := Status{}

	if IsInstalled(cfg.CNGTPath) {
		repo, err := git.PlainOpen(cfg.CNGTPath)
		if err != nil {
			status.RepoStatus = "Error opening repository"
		} else {
			ref, err := repo.Head()
			if err != nil {
				status.RepoStatus = "Installed (unknown commit)"
			} else {
				commit, err := repo.CommitObject(ref.Hash())
				if err != nil {
					status.RepoStatus = "Installed (unknown commit)"
				} else {
					status.RepoStatus = fmt.Sprintf("Installed (commit: %s)", commit.Hash.String()[:7])
				}
			}
		}
	} else {
		status.RepoStatus = "Not installed"
	}

	// Try different Python command variations
	pythonCommands := []string{"python", "python3", "py"}
	status.PythonStatus = "Not found"
	
	for _, pythonCmd := range pythonCommands {
		cmd := exec.Command(pythonCmd, "--version")
		if out, err := cmd.Output(); err == nil {
			status.PythonStatus = strings.TrimSpace(string(out))
			break
		}
	}

	if deps.AreInstalled() {
		status.DepsStatus = "Installed"
	} else {
		status.DepsStatus = "Missing dependencies"
	}

	return status
}