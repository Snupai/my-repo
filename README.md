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
- `cngt-cli self-update` - Update the CLI tool itself
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
cngt-cli self-update
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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.