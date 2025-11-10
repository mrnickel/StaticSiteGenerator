#!/bin/bash

# Script to manually create Arch Linux package

if [ $# -ne 1 ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

VERSION="$1"
VERSION_NO_V="${VERSION#v}"
PKGNAME="ssg"

echo "Creating Arch Linux package for $PKGNAME version $VERSION_NO_V"

# Create package directory structure
mkdir -p "pkg/usr/local/bin"
cp "dist/ssg-linux-amd64" "pkg/usr/local/bin/ssg"
chmod +x "pkg/usr/local/bin/ssg"

# Get file size in a portable way
if command -v du >/dev/null 2>&1; then
    if du -sb pkg/usr >/dev/null 2>&1; then
        # Linux du command
        SIZE=$(du -sb pkg/usr | cut -f1)
    else
        # macOS/BSD du command
        SIZE=$(du -sk pkg/usr | cut -f1)
        SIZE=$((SIZE * 1024))  # Convert KB to bytes
    fi
else
    SIZE="0"
fi

# Create .PKGINFO file
cat > pkg/.PKGINFO << EOF
pkgname = $PKGNAME
pkgbase = $PKGNAME
pkgver = $VERSION_NO_V-1
pkgdesc = Static Site Generator written in Go
url = https://github.com/mrnickel/StaticSiteGenerator
builddate = $(date +%s)
packager = mrnickel <your-email@example.com>
size = $SIZE
arch = x86_64
license = MIT
EOF

# Create simplified .MTREE file
cat > pkg/.MTREE << 'EOF'
#mtree
/set type=file uid=0 gid=0 mode=644
.PKGINFO size=249 time=1699000000.000000000 mode=644 type=file
usr type=dir mode=755
usr/local type=dir mode=755
usr/local/bin type=dir mode=755
usr/local/bin/ssg mode=755 type=file
EOF

# Create the package archive
cd pkg
tar -cf "../$PKGNAME-$VERSION_NO_V-1-x86_64.pkg.tar" .PKGINFO .MTREE usr/
cd ..

# Compress with xz
if command -v xz >/dev/null 2>&1; then
    xz "$PKGNAME-$VERSION_NO_V-1-x86_64.pkg.tar"
    echo "Created: $PKGNAME-$VERSION_NO_V-1-x86_64.pkg.tar.xz"
else
    # Fallback to gzip if xz is not available
    gzip "$PKGNAME-$VERSION_NO_V-1-x86_64.pkg.tar"
    mv "$PKGNAME-$VERSION_NO_V-1-x86_64.pkg.tar.gz" "$PKGNAME-$VERSION_NO_V-1-x86_64.pkg.tar.xz"
    echo "Created: $PKGNAME-$VERSION_NO_V-1-x86_64.pkg.tar.xz (using gzip fallback)"
fi

# Clean up
rm -rf pkg

echo "Arch Linux package created successfully!"