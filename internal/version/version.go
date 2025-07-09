package version

// Version information - can be overridden at build time
var (
	Version   = "1.0.0"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

// GetVersion returns the current version
func GetVersion() string {
	return Version
}

// GetFullVersion returns detailed version information
func GetFullVersion() string {
	return Version + " (commit: " + GitCommit + ", built: " + BuildTime + ")"
}