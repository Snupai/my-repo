package deps

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	for _, pkg := range requiredPackages {
		if !isPackageInstalled(pkg) {
			return false
		}
	}
	return true
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
	
	if isUvAvailable() {
		return installWithUv(cfg)
	} else {
		return installWithPip(cfg)
	}
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

func installWithUv(cfg *config.Config) error {
	reqFile := filepath.Join(cfg.CNGTPath, "requirements.txt")
	if _, err := os.Stat(reqFile); os.IsNotExist(err) {
		return installPackagesDirectly("uv", "add")
	}

	cmd := exec.Command("uv", "pip", "install", "-r", reqFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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