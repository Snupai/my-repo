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
	cmd := exec.Command("python", "--version")
	return cmd.Run() == nil
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
	
	cmd := exec.Command("python", "-c", fmt.Sprintf("import %s", pkgName))
	return cmd.Run() == nil
}

func CheckInteractive() error {
	if isPythonInstalled() {
		fmt.Println("✓ Python is installed")
	} else {
		fmt.Println("✗ Python is not installed")
		fmt.Println("Please install Python from https://python.org or your package manager")
		return fmt.Errorf("Python is required")
	}

	missing := []string{}
	for _, pkg := range requiredPackages {
		if isPackageInstalled(pkg) {
			fmt.Printf("✓ %s is installed\n", pkg)
		} else {
			fmt.Printf("✗ %s is missing\n", pkg)
			missing = append(missing, pkg)
		}
	}

	if len(missing) > 0 {
		fmt.Printf("\nMissing %d dependencies. Would you like to install them? (y/N): ", len(missing))
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return Install()
		} else {
			return fmt.Errorf("dependencies are required to run CNGT")
		}
	}

	fmt.Println("All dependencies are installed!")
	return nil
}