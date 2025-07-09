package deps

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/snupai/cngt-cli/internal/config"
)

var requiredPackages = []string{
	"termcolor",
	"mido",
	"colorama>=0.4.6",
	"cryptography>=42.0.5",
}

func AreInstalled() bool {
	cfg, err := config.Load()
	if err != nil {
		return false
	}
	
	// Check if using uv environment
	if isUvAvailable() && hasUvProject(cfg.CNGTPath) {
		return arePackagesInstalledInUv(cfg.CNGTPath)
	}
	
	// Fallback to system Python check
	for _, pkg := range requiredPackages {
		if !isPackageInstalled(pkg) {
			return false
		}
	}
	return true
}

func arePackagesInstalledInUv(projectPath string) bool {
	for _, pkg := range requiredPackages {
		pkgName := strings.Split(pkg, ">=")[0]
		pkgName = strings.Split(pkgName, "==")[0]
		
		cmd := exec.Command("uv", "run", "python", "-c", fmt.Sprintf("import %s", pkgName))
		cmd.Dir = projectPath
		if cmd.Run() != nil {
			return false
		}
	}
	return true
}

func hasUvProject(path string) bool {
	_, err := os.Stat(filepath.Join(path, "pyproject.toml"))
	return err == nil
}

func Install() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !isPythonInstalled() {
		return fmt.Errorf("Python is not installed. Please install Python first")
	}

	fmt.Println("Installing Python dependencies...")
	
	// Try to install uv if not available
	if !isUvAvailable() {
		fmt.Println("   Installing uv (modern Python package manager)...")
		if err := installUv(); err != nil {
			fmt.Printf("   Failed to install uv, falling back to pip: %v\n", err)
			return installWithPip(cfg)
		}
		fmt.Println("   ✓ uv installed successfully")
	}
	
	return installWithUv(cfg)
}

func isPythonInstalled() bool {
	// Try different Python command variations
	pythonCommands := []string{"python", "python3", "py"}
	
	for _, pythonCmd := range pythonCommands {
		cmd := exec.Command(pythonCmd, "--version")
		if cmd.Run() == nil {
			return true
		}
	}
	return false
}

func isUvAvailable() bool {
	cmd := exec.Command("uv", "--version")
	return cmd.Run() == nil
}

func installUv() error {
	// Install uv using the official installer
	var cmd *exec.Cmd
	
	// Different installation methods for different platforms
	switch runtime.GOOS {
	case "windows":
		// Use PowerShell to install uv on Windows
		cmd = exec.Command("powershell", "-Command", "irm https://astral.sh/uv/install.ps1 | iex")
	default:
		// Use curl for Unix-like systems
		cmd = exec.Command("sh", "-c", "curl -LsSf https://astral.sh/uv/install.sh | sh")
	}
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installWithUv(cfg *config.Config) error {
	// Initialize uv project in the CNGT directory
	fmt.Println("   Setting up Python environment...")
	
	// Store current directory to restore later
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)
	
	// Change to CNGT directory
	if err := os.Chdir(cfg.CNGTPath); err != nil {
		return fmt.Errorf("failed to change to CNGT directory: %w", err)
	}
	
	// Initialize uv project if not already initialized
	pyprojectFile := filepath.Join(cfg.CNGTPath, "pyproject.toml")
	if _, err := os.Stat(pyprojectFile); os.IsNotExist(err) {
		cmd := exec.Command("uv", "init", "--no-readme")
		cmd.Dir = cfg.CNGTPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to initialize uv project: %w", err)
		}
	}
	
	// Install packages using uv
	reqFile := filepath.Join(cfg.CNGTPath, "requirements.txt")
	if _, err := os.Stat(reqFile); os.IsNotExist(err) {
		// Install packages individually if no requirements.txt
		for _, pkg := range requiredPackages {
			fmt.Printf("   Installing %s...\n", pkg)
			cmd := exec.Command("uv", "add", pkg)
			cmd.Dir = cfg.CNGTPath
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install %s: %w", pkg, err)
			}
		}
	} else {
		// Install from requirements.txt
		fmt.Println("   Installing from requirements.txt...")
		
		// First, add the requirements to the project
		cmd := exec.Command("uv", "add", "-r", reqFile)
		cmd.Dir = cfg.CNGTPath
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add requirements from requirements.txt: %w", err)
		}
	}
	
	return nil
}

func installWithPip(cfg *config.Config) error {
	reqFile := filepath.Join(cfg.CNGTPath, "requirements.txt")
	if _, err := os.Stat(reqFile); os.IsNotExist(err) {
		return installPackagesDirectly("pip", "install")
	}

	cmd := exec.Command("pip", "install", "-r", reqFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installPackagesDirectly(tool, action string) error {
	for _, pkg := range requiredPackages {
		fmt.Printf("Installing %s...\n", pkg)
		cmd := exec.Command(tool, action, pkg)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
	}
	return nil
}

func isPackageInstalled(pkg string) bool {
	pkgName := strings.Split(pkg, ">=")[0]
	pkgName = strings.Split(pkgName, "==")[0]
	
	// Try different Python command variations
	pythonCommands := []string{"python", "python3", "py"}
	
	for _, pythonCmd := range pythonCommands {
		cmd := exec.Command(pythonCmd, "-c", fmt.Sprintf("import %s", pkgName))
		if cmd.Run() == nil {
			return true
		}
	}
	return false
}

func CheckInteractive() error {
	if isPythonInstalled() {
		fmt.Println("   ✓ Python is installed")
	} else {
		fmt.Println("   ✗ Python is not installed")
		fmt.Println()
		fmt.Println("   Python is required to run CNGT tools.")
		fmt.Println("   Please install Python from https://python.org")
		fmt.Println("   On Windows, you can also use: winget install Python.Python.3")
		fmt.Println()
		return fmt.Errorf("Python is required but not installed")
	}

	missing := []string{}
	for _, pkg := range requiredPackages {
		if isPackageInstalled(pkg) {
			fmt.Printf("   ✓ %s is installed\n", pkg)
		} else {
			fmt.Printf("   ✗ %s is missing\n", pkg)
			missing = append(missing, pkg)
		}
	}

	if len(missing) > 0 {
		fmt.Printf("\n   Missing %d Python packages. Install them automatically? (Y/n): ", len(missing))
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "n" || response == "no" {
			fmt.Println("   Setup cancelled. You can install packages manually with:")
			for _, pkg := range missing {
				fmt.Printf("   pip install %s\n", pkg)
			}
			return fmt.Errorf("Python dependencies are required but installation was cancelled")
		}
		
		fmt.Println("   Installing Python packages...")
		if err := Install(); err != nil {
			return fmt.Errorf("failed to install Python packages: %w", err)
		}
		fmt.Println("   ✓ Python packages installed successfully")
	} else {
		fmt.Println("   ✓ All Python packages are already installed")
	}

	return nil
}