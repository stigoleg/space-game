#!/bin/bash
# Build script for Stellar Siege
# Creates both standalone binary and macOS .app bundle
# Suppresses macOS 15 deprecation warnings from Ebiten's Metal driver

set -e  # Exit on error

APP_NAME="Stellar Siege"
BINARY_NAME="stellar-siege"
BUNDLE_ID="com.stellarsiege.game"
VERSION="1.0.0"

echo "🚀 Building Stellar Siege..."

# Clean previous builds
rm -f "$BINARY_NAME"
rm -rf "$APP_NAME.app"

# Build the binary with macOS fix
echo "📦 Compiling binary..."
CGO_CFLAGS="-Wno-deprecated-declarations" go build -o "$BINARY_NAME" .

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Binary compiled successfully!"

# Create .app bundle structure
echo "🎨 Creating macOS app bundle..."

APP_DIR="$APP_NAME.app/Contents"
mkdir -p "$APP_DIR/MacOS"
mkdir -p "$APP_DIR/Resources"
mkdir -p "$APP_DIR/Resources/assets"
mkdir -p "$APP_DIR/Resources/config"
mkdir -p "$APP_DIR/Resources/data"

# Copy binary to app bundle (keeping one in root for terminal use)
cp "$BINARY_NAME" "$APP_DIR/MacOS/"

# Copy assets if they exist
if [ -d "assets" ]; then
    echo "📁 Copying assets..."
    cp -r assets/* "$APP_DIR/Resources/assets/" 2>/dev/null || true
fi

if [ -d "config" ]; then
    echo "📁 Copying config..."
    cp -r config/* "$APP_DIR/Resources/config/" 2>/dev/null || true
fi

# Note: User data (leaderboard, progression, achievements) will be saved to
# ~/Library/Application Support/StellarSiege/ when running from .app
echo "💾 User data will be saved to ~/Library/Application Support/StellarSiege/"

# Copy or create Info.plist
if [ -f "build_resources/Info.plist" ]; then
    cp "build_resources/Info.plist" "$APP_DIR/"
    echo "📝 Using Info.plist from build_resources/"
else
    echo "📝 Generating Info.plist..."
    cat > "$APP_DIR/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>$BINARY_NAME</string>
    <key>CFBundleName</key>
    <string>$APP_NAME</string>
    <key>CFBundleDisplayName</key>
    <string>$APP_NAME</string>
    <key>CFBundleIdentifier</key>
    <string>$BUNDLE_ID</string>
    <key>CFBundleVersion</key>
    <string>$VERSION</string>
    <key>CFBundleShortVersionString</key>
    <string>$VERSION</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleSignature</key>
    <string>STSG</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.15</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>LSApplicationCategoryType</key>
    <string>public.app-category.games</string>
</dict>
</plist>
EOF
fi

# Copy icon if exists
if [ -f "build_resources/icon.icns" ]; then
    cp "build_resources/icon.icns" "$APP_DIR/Resources/"
    echo "🎨 App icon added"
fi

# Create PkgInfo
echo "APPLSTSG" > "$APP_DIR/PkgInfo"

# Make binary executable
chmod +x "$APP_DIR/MacOS/$BINARY_NAME"

# Remove quarantine attribute to prevent "damaged" warning
xattr -cr "$APP_NAME.app" 2>/dev/null || true

echo ""
echo "✅ Build complete!"
echo ""
echo "📦 Created:"
echo "   - $BINARY_NAME (terminal executable)"
echo "   - $APP_NAME.app (double-click app bundle)"
echo ""
echo "🎮 To run:"
echo "   - Double-click '$APP_NAME.app' in Finder"
echo "   - Or: open '$APP_NAME.app'"
echo "   - Or: ./$BINARY_NAME (terminal)"
echo ""
echo "📦 To distribute:"
echo "   - Zip: zip -r Stellar-Siege-v$VERSION.zip '$APP_NAME.app'"
echo "   - DMG: hdiutil create -volname '$APP_NAME' -srcfolder '$APP_NAME.app' -ov -format UDZO Stellar-Siege-v$VERSION.dmg"
echo ""
