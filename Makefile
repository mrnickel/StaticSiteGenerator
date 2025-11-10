# Static Site Generator Makefile

VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
LDFLAGS = -ldflags="-s -w -X main.version=$(VERSION)"

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	go build $(LDFLAGS) -o ssg .

# Build for release platforms
.PHONY: build-release
build-release: clean
	mkdir -p dist/release
	
	# macOS ARM64 (M1 Macs)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/ssg-darwin-arm64 .
	
	# Linux x86_64 (Intel - for Arch Linux)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/ssg-linux-amd64 .

# Create release packages
.PHONY: package
package: build-release
	@echo "Creating release packages..."
	
	# Create macOS tarball for Homebrew
	mkdir -p dist/release/darwin-arm64
	cp dist/ssg-darwin-arm64 dist/release/darwin-arm64/ssg
	chmod +x dist/release/darwin-arm64/ssg
	cd dist/release && tar -czf ../../darwin-arm64-$(VERSION).tar.gz -C darwin-arm64 ssg
	
	# Create Linux tarball for Arch
	mkdir -p dist/release/linux-amd64
	cp dist/ssg-linux-amd64 dist/release/linux-amd64/ssg
	chmod +x dist/release/linux-amd64/ssg
	cd dist/release && tar -czf ../../linux-amd64-$(VERSION).tar.gz -C linux-amd64 ssg

# Create Arch Linux package
.PHONY: arch-package
arch-package: build-release
	@echo "Creating Arch Linux package..."
	./scripts/create-arch-package.sh $(VERSION)

# Create Homebrew formula
.PHONY: homebrew-formula
homebrew-formula: package
	@echo "Creating Homebrew formula..."
	./scripts/generate-homebrew-formula.sh $(VERSION) darwin-arm64-$(VERSION).tar.gz

# Full release build
.PHONY: release
release: package arch-package homebrew-formula
	@echo "Release packages created:"
	@echo "  - darwin-arm64-$(VERSION).tar.gz (Homebrew)"
	@echo "  - linux-amd64-$(VERSION).tar.gz (Manual install)"
	@echo "  - ssg-$(VERSION:v%=%)-1-x86_64.pkg.tar.xz (Arch Linux)"
	@echo "  - ssg.rb (Homebrew formula)"

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf dist/
	rm -f ssg
	rm -f *.tar.gz
	rm -f *.pkg.tar.xz
	rm -f ssg.rb

# Test
.PHONY: test
test:
	go test -v ./...

# Install locally
.PHONY: install
install: build
	cp ssg /usr/local/bin/

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build           - Build for current platform"
	@echo "  build-release   - Build for release platforms (macOS ARM64, Linux x86_64)"
	@echo "  package         - Create release tarballs"
	@echo "  arch-package    - Create Arch Linux package"
	@echo "  homebrew-formula - Create Homebrew formula"
	@echo "  release         - Create all release packages"
	@echo "  clean           - Clean build artifacts"
	@echo "  test            - Run tests"
	@echo "  install         - Install locally to /usr/local/bin/"
	@echo "  help            - Show this help"