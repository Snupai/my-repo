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
}

var migrateCmd = &cobra.Command{
	Use:   "migrate [args...]",
	Short: "Run GlyphMigrate.py with the given arguments",
	Long:  "Execute the GlyphMigrate.py script from the CNGT repository",
	Run: func(cmd *cobra.Command, args []string) {
		if err := ensureSetup(); err != nil {
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
		if err := ensureSetup(); err != nil {
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
		if err := ensureSetup(); err != nil {
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
		if err := cngt.Update(); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating CNGT: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("CNGT repository updated successfully")
	},
}

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
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
	},
}

func ensureSetup() error {
	// Check for updates on first setup (weekly)
	go func() {
		if shouldCheckForUpdates() {
			if release, hasUpdate, err := updater.CheckForUpdates(); err == nil && hasUpdate {
				fmt.Printf("\n🔄 New version %s available! Run 'cngt-cli self-update' to upgrade.\n\n", release.TagName)
			}
		}
	}()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !cngt.IsInstalled(cfg.CNGTPath) {
		fmt.Println("CNGT repository not found. Installing...")
		if err := cngt.Install(cfg.CNGTPath); err != nil {
			return fmt.Errorf("failed to install CNGT: %w", err)
		}
		fmt.Println("CNGT repository installed successfully")
	}

	if !deps.AreInstalled() {
		fmt.Println("Python dependencies not found. Installing...")
		if err := deps.Install(); err != nil {
			return fmt.Errorf("failed to install dependencies: %w", err)
		}
		fmt.Println("Dependencies installed successfully")
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
	rootCmd.AddCommand(selfUpdateCmd)
	rootCmd.AddCommand(statusCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}