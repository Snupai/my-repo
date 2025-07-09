# CNGT CLI

A cross-platform CLI tool that wraps the [custom-nothing-glyph-tools](https://github.com/SebiAi/custom-nothing-glyph-tools) repository, providing easy installation, dependency management, and usage from any directory.

## Features

- **Easy Installation**: One-command setup of CNGT repository and dependencies
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Dependency Management**: Automatically installs Python dependencies using uv or pip
- **Auto-Updates**: Keep both the CLI tool and CNGT repository up to date
- **Simple Interface**: Use CNGT tools from any directory

## Installation

### Quick Install (Recommended)

```bash
# Linux/macOS
curl -sSL https://github.com/snupai/cngt-cli/releases/latest/download/cngt-cli-linux-amd64 -o cngt-cli
chmod +x cngt-cli
sudo mv cngt-cli /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/snupai/cngt-cli/releases/latest/download/cngt-cli-windows-amd64.exe" -OutFile "cngt-cli.exe"
```

### Manual Installation

1. Download the appropriate binary for your platform from the [releases page](https://github.com/snupai/cngt-cli/releases)
2. Make it executable and add to your PATH

## Usage

### First Run

The first time you run any CNGT command, the tool will automatically:
1. Download the CNGT repository
2. Check for Python installation
3. Install required Python dependencies

```bash
cngt-cli migrate --help
```

### Available Commands

- `cngt-cli migrate [args...]` - Run GlyphMigrate.py
- `cngt-cli modder [args...]` - Run GlyphModder.py  
- `cngt-cli translator [args...]` - Run GlyphTranslator.py
- `cngt-cli update` - Update CNGT repository
- `cngt-cli upgrade` - Update the CLI tool itself
- `cngt-cli status` - Show installation status
- `cngt-cli --help` - Show help information

### Examples

```bash
# Check status
cngt-cli status

# Run GlyphModder
cngt-cli modder input.glyph output.glyph

# Run GlyphTranslator
cngt-cli translator --format json input.glyph

# Update everything
cngt-cli update
cngt-cli upgrade
```

## Requirements

- Python 3.x
- Internet connection for initial setup and updates

The tool will automatically install these Python packages:
- termcolor
- mido
- colorama>=0.4.6
- cryptography>=42.0.5

## Building from Source

```bash
# Clone the repository
git clone https://github.com/snupai/cngt-cli.git
cd cngt-cli

# Install dependencies
make install

# Build for current platform
make dev

# Build for all platforms
make build
```

## Configuration

The tool stores data in platform-specific locations:
- **Linux**: `~/.local/share/cngt-cli/`
- **macOS**: `~/Library/Application Support/cngt-cli/`
- **Windows**: `%APPDATA%\cngt-cli\`

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Disclaimer

This tool is a wrapper around the custom-nothing-glyph-tools repository. It is not affiliated with Nothing Technology Limited. The underlying tools are provided as-is without warranty.

## Releasing New Versions

### For Maintainers

This project uses automated GitHub Actions for releasing new versions. There are multiple ways to trigger a new release:

#### Method 1: Using the Release Script (Recommended)

```bash
# Create and push a new release tag
make release VERSION=1.0.1
```

This will:
1. Validate the version format (semantic versioning)
2. Update the version in source code
3. Create a git commit and tag
4. Push to GitHub, triggering the automated release

#### Method 2: Manual Git Tagging

```bash
# Update version in source code
sed -i 's/Version   = ".*"/Version   = "1.0.1"/' internal/version/version.go

# Commit the version change
git add internal/version/version.go
git commit -m "chore: bump version to 1.0.1"

# Create and push tag
git tag -a v1.0.1 -m "Release version 1.0.1"
git push origin main
git push origin v1.0.1
```

#### Method 3: GitHub Actions Manual Trigger

1. Go to the [Actions tab](https://github.com/snupai/cngt-cli/actions)
2. Select the "Release" workflow
3. Click "Run workflow"
4. Enter the version number (e.g., `1.0.1`)
5. Click "Run workflow"

### What Happens During Release

The automated release process:

1. **Builds binaries** for all supported platforms:
   - Windows AMD64
   - Linux AMD64/ARM64
   - macOS AMD64/ARM64

2. **Generates changelog** from git commits since the last release

3. **Creates GitHub release** with:
   - Release notes and changelog
   - Installation instructions
   - All platform binaries
   - Checksums for verification

4. **Updates self-update mechanism** so users can upgrade automatically

### Release Requirements

- Version must follow semantic versioning (e.g., `1.0.1`, `2.1.0`)
- Must be on `main` branch with clean working directory
- All tests must pass
- GitHub repository must have appropriate permissions

### Testing a Release

After creating a release:

```bash
# Test the release binaries
curl -sSL https://github.com/snupai/cngt-cli/releases/latest/download/cngt-cli-linux-amd64 -o cngt-cli-test
chmod +x cngt-cli-test
./cngt-cli-test --version
./cngt-cli-test --help
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test -v ./...`
5. Build and test: `make dev`
6. Submit a Pull Request

### Code Style

- Follow Go conventions
- Run `go fmt` before committing
- Add tests for new functionality
- Update documentation as needed