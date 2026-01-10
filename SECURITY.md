# Security & Installation Guide

## Why Am I Seeing Security Warnings?

When you download and run Stellar Siege, you may see security warnings from your operating system. **This is normal and expected behavior** for applications that are not code-signed with paid developer certificates.

Stellar Siege is an open-source game built from publicly available source code. The warnings appear because:

1. **macOS Gatekeeper**: Apple requires a $99/year Developer Program membership to sign apps
2. **Windows SmartScreen**: Microsoft requires a $100-400/year code signing certificate
3. **Community Distribution**: We currently distribute builds without these paid certificates to keep the game free

**The application is safe** - you can verify this by:
- Reviewing the [source code](https://github.com/sogud/stellar-siege) yourself
- Checking the SHA256 checksums (see below)
- Building from source (instructions below)

---

## Verifying Your Download (Recommended)

Before running the game, verify the integrity of your download using SHA256 checksums:

### macOS/Linux:
```bash
# Download checksums.txt from the release
# Navigate to your download folder
cd ~/Downloads

# Verify the checksum
shasum -a 256 Stellar-Siege-macOS-Intel.dmg
# Compare output with checksums.txt
```

### Windows (PowerShell):
```powershell
# Navigate to your download folder
cd $env:USERPROFILE\Downloads

# Verify the checksum
Get-FileHash stellar-siege-windows-amd64.zip -Algorithm SHA256
# Compare output with checksums.txt
```

**Important**: If the checksums don't match, **do not run the file** - re-download it and verify again.

---

## Installing on macOS

### First-Time Installation

1. **Download** the appropriate DMG file:
   - `Stellar-Siege-macOS-Intel.dmg` (Intel Macs)
   - `Stellar-Siege-macOS-AppleSilicon.dmg` (M1/M2/M3/M4 Macs)

2. **Verify checksum** (recommended - see above)

3. **Mount the DMG** by double-clicking it

4. **Drag Stellar Siege** to your Applications folder

5. **First Launch** - You will see a Gatekeeper warning:

   ```
   "Stellar Siege" cannot be opened because it is from an unidentified developer.
   ```

### Bypassing Gatekeeper (Safe Method)

**Option 1: Right-Click Method (Easiest)**
1. Locate `Stellar Siege.app` in your Applications folder
2. **Right-click** (or Control-click) on the app
3. Select **"Open"** from the menu
4. Click **"Open"** in the confirmation dialog
5. macOS will remember your choice - future launches work normally

**Option 2: System Settings Method**
1. Try to open the app normally (you'll get the warning)
2. Open **System Settings** → **Privacy & Security**
3. Scroll down to the **Security** section
4. You'll see: `"Stellar Siege" was blocked from use because it is not from an identified developer`
5. Click **"Open Anyway"**
6. Confirm by clicking **"Open"**

**Option 3: Command Line Method**
```bash
# Remove the quarantine attribute
xattr -d com.apple.quarantine "/Applications/Stellar Siege.app"
```

### Why This Works
These methods tell macOS "I trust this application" - they don't bypass security, they grant explicit permission for this specific app.

---

## Installing on Windows

### First-Time Installation

1. **Download** `stellar-siege-windows-amd64.zip`

2. **Verify checksum** (recommended - see above)

3. **Extract** the ZIP file to a location of your choice (e.g., `C:\Games\Stellar-Siege`)

4. **First Launch** - You may see a SmartScreen warning:

   ```
   Windows protected your PC
   Microsoft Defender SmartScreen prevented an unrecognized app from starting.
   ```

### Bypassing SmartScreen (Safe Method)

1. Click **"More info"** in the SmartScreen dialog
2. Click **"Run anyway"** button
3. Windows will remember your choice - future launches work normally

### Alternative: Disable SmartScreen (Not Recommended)
If you frequently run unsigned applications, you can disable SmartScreen:
1. Open **Windows Security** → **App & browser control**
2. Click **"Reputation-based protection settings"**
3. Toggle off **"Check apps and files"**

**Warning**: Only disable SmartScreen if you understand the security implications and trust the software you download.

---

## Installing on Linux

### Installation

1. **Download** `stellar-siege-linux-amd64.tar.gz`

2. **Verify checksum** (recommended - see above)

3. **Extract** the archive:
   ```bash
   tar -xzf stellar-siege-linux-amd64.tar.gz
   cd stellar-siege-linux
   ```

4. **Make executable** (if needed):
   ```bash
   chmod +x stellar-siege
   ```

5. **Run**:
   ```bash
   ./stellar-siege
   ```

Linux typically does not show security warnings for unsigned binaries, but you should still verify checksums.

---

## Building from Source

If you prefer to build the game yourself, you can verify every line of code:

### Prerequisites
- Go 1.24.0 or later
- Platform-specific dependencies:
  - **Linux**: `libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config`
  - **macOS**: Xcode Command Line Tools
  - **Windows**: GCC (MinGW-w64)

### Build Instructions

```bash
# Clone the repository
git clone https://github.com/sogud/stellar-siege.git
cd stellar-siege

# Install dependencies
go mod download

# Build the game
go build -ldflags="-s -w" -o stellar-siege .

# Run
./stellar-siege  # macOS/Linux
# or
stellar-siege.exe  # Windows
```

---

## Frequently Asked Questions

### Q: Is this game safe to run?
**A**: Yes. The game is open-source - you can review the code yourself. The security warnings appear because we don't pay for code signing certificates, not because the software is malicious.

### Q: Why don't you just sign the application?
**A**: Code signing requires paid developer certificates:
- Apple Developer Program: $99/year
- Windows Code Signing Certificate: $100-400/year

As an open-source project, we prioritize keeping the game free. If funding becomes available, we may add code signing in the future.

### Q: Can I trust this application?
**A**: You can verify trust by:
1. Reviewing the [source code](https://github.com/sogud/stellar-siege)
2. Checking SHA256 checksums match the official release
3. Building from source yourself
4. Checking the GitHub repository's history and contributors

### Q: What data does the game collect?
**A**: Stellar Siege optionally submits high scores to a GitHub Gist leaderboard (requires manual configuration). No other data is collected or transmitted. The game runs entirely locally.

### Q: Will you add code signing in the future?
**A**: If the project receives sufficient community support or funding, we will consider purchasing code signing certificates to eliminate these warnings.

---

## Reporting Security Issues

If you discover a security vulnerability in Stellar Siege, please report it via:
- GitHub Issues: https://github.com/sogud/stellar-siege/issues
- Email: [maintainer email here]

We take security seriously and will address legitimate issues promptly.

---

## Additional Resources

- **Source Code**: https://github.com/sogud/stellar-siege
- **Latest Release**: https://github.com/sogud/stellar-siege/releases/latest
- **Report Issues**: https://github.com/sogud/stellar-siege/issues
- **Build Instructions**: See README.md

---

*Last Updated: January 2026*
