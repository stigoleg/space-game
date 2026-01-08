#!/bin/bash
# Build script for Stellar Siege
# Suppresses macOS 15 deprecation warnings from Ebiten's Metal driver

echo "Building Stellar Siege..."
CGO_CFLAGS="-Wno-deprecated-declarations" go build -o stellar-siege .

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
    echo "Run with: ./stellar-siege"
else
    echo "✗ Build failed"
    exit 1
fi
