# Installation Guide

This document describes how to install the Static Site Generator on supported platforms.

## Homebrew (macOS - Apple Silicon/M1)

```bash
# Add the tap (first time only)
brew tap mrnickel/ssg https://github.com/mrnickel/StaticSiteGenerator

# Install
brew install ssg
```

**Note:** This formula is optimized for Apple Silicon Macs (M1/M2). For Intel Macs, use the manual installation method below.

## Arch Linux (x86_64)

```bash
# Download the package
wget https://github.com/mrnickel/StaticSiteGenerator/releases/latest/download/ssg-<version>-1-x86_64.pkg.tar.xz

# Install with pacman
sudo pacman -U ssg-<version>-1-x86_64.pkg.tar.xz
```

## Manual Installation

### Download Pre-built Binaries

1. Go to the [Releases page](https://github.com/mrnickel/StaticSiteGenerator/releases)
2. Download the appropriate binary for your platform:

   - `darwin-arm64-<version>.tar.gz` - macOS Apple Silicon (M1/M2)
   - `linux-amd64-<version>.tar.gz` - Linux x86_64

3. Extract and install:

```bash
# Extract the archive
tar -xzf <downloaded-file>.tar.gz

# Move to a directory in your PATH
sudo mv ssg /usr/local/bin/

# Make executable (Linux/macOS)
chmod +x /usr/local/bin/ssg
```

### Build from Source

Requirements:

- Go 1.25 or later

```bash
# Clone the repository
git clone https://github.com/mrnickel/StaticSiteGenerator.git
cd StaticSiteGenerator

# Build and install
make build
sudo make install

# Or build for all platforms
make build-all
```

## Verification

After installation, verify it works:

```bash
ssg version
ssg help
```

## Usage

```bash
# Create a new post
ssg create "My First Post"

# Publish a post
ssg publish "My First Post"

# Start a development server
ssg standup

# Start server on custom port
ssg standup 3000

# View site statistics
ssg stats

# List draft posts
ssg listdrafts

# Regenerate all published posts
ssg regenerate
```

## Updating

### Homebrew

```bash
brew update
brew upgrade ssg
```

### Arch Linux

Download and install the latest package following the installation steps above.

### Manual Installation

Download the latest binary from the releases page and replace your existing installation.
