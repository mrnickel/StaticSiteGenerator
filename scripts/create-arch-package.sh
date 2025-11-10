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

# Get file size (installed size in bytes)
if du -sb pkg/usr >/dev/null 2>&1; then
    SIZE=$(du -sb pkg/usr | cut -f1)
else
    # Fallback for systems without -b option
    SIZE=$(du -sk pkg/usr | cut -f1)
    SIZE=$((SIZE * 1024))
fi

# Get binary size and hash for MTREE
BINSIZE=$(stat -c%s "pkg/usr/local/bin/ssg" 2>/dev/null || stat -f%z "pkg/usr/local/bin/ssg")
BUILDDATE=$(date +%s)

# Create .PKGINFO file
cat > pkg/.PKGINFO << EOF
pkgname = $PKGNAME
pkgbase = $PKGNAME
pkgver = $VERSION_NO_V-1
pkgdesc = Static Site Generator written in Go
url = https://github.com/mrnickel/StaticSiteGenerator
builddate = $BUILDDATE
packager = mrnickel
size = $SIZE
arch = x86_64
EOF

# Create .MTREE file with proper metadata
cat > pkg/.MTREE << EOF
#mtree
/set type=file uid=0 gid=0 mode=644
./.PKGINFO time=$BUILDDATE.0 mode=644 type=file
./usr time=$BUILDDATE.0 mode=755 type=dir
./usr/local time=$BUILDDATE.0 mode=755 type=dir
./usr/local/bin time=$BUILDDATE.0 mode=755 type=dir
./usr/local/bin/ssg size=$BINSIZE time=$BUILDDATE.0 mode=755 type=file
EOF

# Create the package archive with proper ordering
cd pkg
tar -czf "../dist/$PKGNAME-$VERSION_NO_V-1-x86_64.pkg.tar.zst" --format=gnu --sort=name .PKGINFO .MTREE usr/
cd ..

# Clean up
rm -rf pkg

echo "Created: dist/$PKGNAME-$VERSION_NO_V-1-x86_64.pkg.tar.zst"
echo "Arch Linux package created successfully!"