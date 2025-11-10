#!/bin/bash

# Script to generate Homebrew formula

if [ $# -ne 2 ]; then
    echo "Usage: $0 <version> <tarball>"
    exit 1
fi

VERSION="$1"
TARBALL="$2"

# Remove 'v' prefix from version
VERSION_NO_V="${VERSION#v}"

# Calculate SHA256
SHA256=$(shasum -a 256 "$TARBALL" | cut -d' ' -f1)

# Generate the formula
cat > ssg.rb << EOF
class Ssg < Formula
  desc "Static Site Generator written in Go"
  homepage "https://github.com/mrnickel/StaticSiteGenerator"
  version "$VERSION_NO_V"
  url "https://github.com/mrnickel/StaticSiteGenerator/releases/download/$VERSION/darwin-arm64-$VERSION.tar.gz"
  sha256 "$SHA256"

  def install
    bin.install "ssg"
  end

  test do
    system "#{bin}/ssg", "help"
  end
end
EOF

echo "Generated ssg.rb for version $VERSION_NO_V"
echo "SHA256: $SHA256"