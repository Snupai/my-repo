package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/snupai/cngt-cli/internal/cngt"
	"github.com/snupai/cngt-cli/internal/config"
	"github.com/snupai/cngt-cli/internal/deps"
	"github.com/snupai/cngt-cli/internal/updater"
	"github.com/snupai/cngt-cli/internal/version"
)

var rootCmd = &cobra.Command{
	Use:   "cngt-cli",
	Short: "CLI tool for Custom Nothing Glyph Tools",
	Long: `A cross-platform CLI tool that wraps the custom-nothing-glyph-tools repository,
providing easy installation, dependency management, and usage from any directory.`,
	Version: version.GetVersion(),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Check for updates on any command run
		checkForUpdatesAsync()
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate [args...]",
	Short: "Run GlyphMigrate.py with the given arguments",
	Long:  "Execute the GlyphMigrate.py script from the CNGT repository",
	Run: func(cmd *cobra.Command, args []string) {
		if err := performSetupIfNeeded(); err != nil {
			fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
			os.Exit(1)
		}
		
		if err := cngt.RunScript("GlyphMigrate.py", args); err != nil {
			fmt.Fprintf(os.Stderr, "Error running migrate: %v\n", err)
			os.Exit(1)
		}
	},
}

var modderCmd = &cobra.Command{
	Use:   "modder [args...]",
	Short: "Run GlyphModder.py with the given arguments",
	Long:  "Execute the GlyphModder.py script from the CNGT repository",
	Run: func(cmd *cobra.Command, args []string) {
		if err := performSetupIfNeeded(); err != nil {
			fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
			os.Exit(1)
		}
		
		if err := cngt.RunScript("GlyphModder.py", args); err != nil {
			fmt.Fprintf(os.Stderr, "Error running modder: %v\n", err)
			os.Exit(1)
		}
	},
}

var translatorCmd = &cobra.Command{
	Use:   "translator [args...]",
	Short: "Run GlyphTranslator.py with the given arguments",
	Long:  "Execute the GlyphTranslator.py script from the CNGT repository",
	Run: func(cmd *cobra.Command, args []string) {
		if err := performSetupIfNeeded(); err != nil {
			fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
			os.Exit(1)
		}
		
		if err := cngt.RunScript("GlyphTranslator.py", args); err != nil {
			fmt.Fprintf(os.Stderr, "Error running translator: %v\n", err)
			os.Exit(1)
		}
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the CNGT repository to the latest version",
	Long:  "Pull the latest changes from the custom-nothing-glyph-tools repository",
	Run: func(cmd *cobra.Command, args []string) {
		if err := performSetupIfNeeded(); err != nil {
			fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
			os.Exit(1)
		}
		
		if err := cngt.Update(); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating CNGT: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("CNGT repository updated successfully")
	},
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Update the cngt-cli tool itself",
	Long:  "Check for and install updates to the cngt-cli tool",
	Run: func(cmd *cobra.Command, args []string) {
		if err := updater.SelfUpdate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating CLI: %v\n", err)
			os.Exit(1)
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of CNGT installation and dependencies",
	Long:  "Display information about the CNGT repository and Python dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		status := cngt.GetStatus()
		fmt.Printf("CNGT CLI Version: %s\n", version.GetFullVersion())
		fmt.Printf("CNGT Repository: %s\n", status.RepoStatus)
		fmt.Printf("Python: %s\n", status.PythonStatus)
		fmt.Printf("Dependencies: %s\n", status.DepsStatus)
		
		// If nothing is installed, offer to set up
		if status.RepoStatus == "Not installed" {
			fmt.Println()
			fmt.Println("💡 Tip: Run 'cngt-cli setup' to install CNGT and dependencies interactively")
		}
	},
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup of CNGT repository and dependencies",
	Long:  "Guides you through the installation of the CNGT repository and Python dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		if err := interactiveSetup(); err != nil {
			fmt.Fprintf(os.Stderr, "Setup failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Setup completed successfully!")
		fmt.Println("You can now use cngt-cli commands like 'cngt-cli migrate --help'")
	},
}


func checkForUpdatesAsync() {
	// Check for updates on first setup (weekly)
	go func() {
		if shouldCheckForUpdates() {
			if release, hasUpdate, err := updater.CheckForUpdates(); err == nil && hasUpdate {
				fmt.Printf("\n🔄 New version %s available! Run 'cngt-cli upgrade' to update.\n\n", release.TagName)
			}
		}
	}()
}

func performSetupIfNeeded() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	needsSetup := false
	if !cngt.IsInstalled(cfg.CNGTPath) {
		needsSetup = true
	}
	if !deps.AreInstalled() {
		needsSetup = true
	}

	if needsSetup {
		fmt.Println("📋 CNGT CLI - First Time Setup")
		fmt.Println("This tool requires the CNGT repository and Python dependencies.")
		fmt.Println()
		fmt.Print("Would you like to install everything now? (Y/n): ")
		
		var response string
		fmt.Scanln(&response)
		
		if response == "n" || response == "N" || response == "no" || response == "No" {
			fmt.Println("Setup cancelled. You can run 'cngt-cli setup' anytime to install.")
			return fmt.Errorf("setup required but cancelled by user")
		}
		
		return interactiveSetup()
	}

	return nil
}

func interactiveSetup() error {
	fmt.Println("🚀 Starting CNGT CLI Setup...")
	fmt.Println()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Install CNGT repository
	if !cngt.IsInstalled(cfg.CNGTPath) {
		fmt.Println("📦 Installing CNGT repository...")
		fmt.Println("   Repository: https://github.com/SebiAi/custom-nothing-glyph-tools")
		fmt.Printf("   Location: %s\n", cfg.CNGTPath)
		fmt.Println()
		
		if err := cngt.Install(cfg.CNGTPath); err != nil {
			return fmt.Errorf("failed to install CNGT repository: %w", err)
		}
		fmt.Println("✅ CNGT repository installed successfully")
		fmt.Println()
	} else {
		fmt.Println("✅ CNGT repository already installed")
		fmt.Println()
	}

	// Check and install Python dependencies
	if !deps.AreInstalled() {
		fmt.Println("🐍 Installing Python dependencies...")
		fmt.Println("   Required packages: termcolor, mido, colorama>=0.4.6, cryptography>=42.0.5")
		fmt.Println()
		
		if err := deps.CheckInteractive(); err != nil {
			return fmt.Errorf("failed to install dependencies: %w", err)
		}
		fmt.Println("✅ Python dependencies installed successfully")
		fmt.Println()
	} else {
		fmt.Println("✅ Python dependencies already installed")
		fmt.Println()
	}

	return nil
}

func shouldCheckForUpdates() bool {
	cfg, err := config.Load()
	if err != nil {
		return false
	}

	lastCheckFile := filepath.Join(cfg.DataDir, "last_update_check")
	if info, err := os.Stat(lastCheckFile); err == nil {
		// Check if more than a week has passed
		return time.Since(info.ModTime()) > 7*24*time.Hour
	}

	// Create the file for future checks
	os.WriteFile(lastCheckFile, []byte(time.Now().Format(time.RFC3339)), 0644)
	return true
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(modderCmd)
	rootCmd.AddCommand(translatorCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(setupCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}