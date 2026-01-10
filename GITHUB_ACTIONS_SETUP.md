# GitHub Actions Release Pipeline Setup

This document explains how to set up GitHub Actions for automated builds and releases of Stellar Siege across multiple platforms.

## Prerequisites

1. A GitHub repository for your game
2. GitHub Personal Access Token for Gist leaderboard functionality
3. GitHub Gist ID for storing leaderboard data

## Required GitHub Secrets

You need to configure the following secrets in your GitHub repository:

### 1. GIST_ID
- **Description**: The ID of the GitHub Gist used to store leaderboard data
- **How to get it**: 
  1. Go to https://gist.github.com
  2. Create a new gist (can be private or public)
  3. Name it something like "stellar-siege-leaderboard"
  4. Add initial content: `{"scores": []}`
  5. The Gist ID is in the URL: `https://gist.github.com/username/GIST_ID_HERE`

### 2. GH_GIST_TOKEN
- **Description**: GitHub Personal Access Token with gist permissions
- **Important**: Must be named `GH_GIST_TOKEN` (GitHub reserves names starting with `GITHUB_`)
- **How to create it**:
  1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
  2. Click "Generate new token (classic)"
  3. Give it a descriptive name: "Stellar Siege Gist Access"
  4. Select the **gist** scope (allows creating and modifying gists)
  5. Click "Generate token"
  6. **IMPORTANT**: Copy the token immediately - you won't see it again!

## Setting Up Secrets in GitHub

1. Navigate to your repository on GitHub
2. Go to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add each secret:

   **Secret 1:**
   - Name: `GIST_ID`
   - Value: Your Gist ID (e.g., `a1b2c3d4e5f6g7h8i9j0`)
   
   **Secret 2:**
   - Name: `GH_GIST_TOKEN`
   - Value: Your GitHub Personal Access Token (e.g., `ghp_xxxxxxxxxxxxxxxxxxxx`)
   
   **Note**: The workflow will validate that these secrets are present before building.

## Workflow Files

Two GitHub Actions workflows have been created:

### 1. `.github/workflows/build.yml`
- **Trigger**: Push to main/master/develop branches, or pull requests
- **Purpose**: Continuous Integration - runs tests and linting
- **Platforms**: Tests on Linux, macOS, and Windows
- **Actions**:
  - Runs `go vet`
  - Runs unit tests with race detection
  - Runs golangci-lint
  - Uploads coverage reports

### 2. `.github/workflows/release.yml`
- **Trigger**: 
  - Push tags matching `v*` (e.g., `v1.0.0`)
  - Manual workflow dispatch
- **Purpose**: Build and release binaries for all platforms
- **Platforms**:
  - Linux (amd64)
  - macOS Intel (amd64)
  - macOS Apple Silicon (arm64)
  - Windows (amd64)
- **Actions**:
  - Builds binaries for each platform
  - Creates macOS .app bundles and DMG files
  - Creates tar.gz for Linux
  - Creates zip for Windows
  - Creates GitHub release with all artifacts
  - Includes secrets in .env file during build

## Creating a Release

### Option 1: Using Git Tags (Recommended)

```bash
# Create and push a version tag
git tag v1.0.0
git push origin v1.0.0
```

The release workflow will automatically:
1. Build for all platforms
2. Create a GitHub release
3. Upload all binaries as release assets

### Option 2: Manual Workflow Dispatch

1. Go to your repository on GitHub
2. Click **Actions** tab
3. Select **Release Build** workflow
4. Click **Run workflow**
5. Enter the version (e.g., `v1.0.0`)
6. Click **Run workflow**

## Release Artifacts

After a successful release, the following files will be available:

- `Stellar-Siege-macOS-Intel.dmg` - macOS disk image for Intel Macs
- `Stellar-Siege-macOS-AppleSilicon.dmg` - macOS disk image for M1/M2/M3 Macs
- `stellar-siege-linux-amd64.tar.gz` - Linux tarball with binary and assets
- `stellar-siege-windows-amd64.zip` - Windows zip with exe and assets

Each package includes:
- Compiled binary/executable
- Assets folder (sprites, sounds, etc.)
- Config folder (game configuration)
- `.env.example` file for users to configure their own Gist leaderboard

## Environment Variables in Release

The release workflow automatically creates a `.env` file during the build process with:

```env
GIST_ID=<value from GitHub secret>
GH_GIST_TOKEN=<value from GitHub secret>
GIST_ENABLED=true
```

This `.env` file is then **copied into each release package**:
- **Linux**: Placed in the root of the tar.gz package
- **Windows**: Placed in the root of the zip package  
- **macOS**: Placed in `Stellar Siege.app/Contents/Resources/.env`

The game will automatically load these values when started, enabling the online leaderboard functionality by default.

## Security Notes

1. **Never commit `.env` file**: The `.gitignore` file should include `.env`
2. **Rotate tokens periodically**: Update your GitHub token every 6-12 months
3. **Use least privilege**: The token only needs `gist` scope
4. **Monitor usage**: Check your Gist for unauthorized modifications
5. **Secret validation**: The workflow validates secrets are present before building
6. **Secret naming**: Use `GH_GIST_TOKEN` (GitHub reserves names starting with `GITHUB_`)

## Troubleshooting

### Build fails with "CGO_ENABLED" error
- The workflow sets `CGO_ENABLED=1` which is required for Ebiten
- Ensure all dependencies are installed (handled automatically in workflow)

### macOS build fails
- macOS builds run on `macos-latest` which includes Xcode and required frameworks
- Check that `build_resources/Info.plist` exists or workflow will generate one

### Windows build fails
- Ensure 7z is available (pre-installed on windows-latest runner)
- Check that paths use forward slashes in bash sections

### Release not created
- Verify `GITHUB_TOKEN` has `contents: write` permission (set in workflow)
- Check that the tag starts with `v` (e.g., `v1.0.0` not `1.0.0`)

## Testing Locally

Before pushing tags, test the build locally:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o stellar-siege .

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o stellar-siege .

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o stellar-siege .

# Windows
GOOS=windows GOARCH=amd64 go build -o stellar-siege.exe .
```

## Next Steps

1. Set up the GitHub secrets as described above
2. Test the CI workflow by pushing a commit to main/master
3. Create your first release by tagging a version
4. Download and test the release artifacts
5. Share the release with users!

## Support

For issues with the workflows:
- Check the Actions tab for detailed logs
- Verify all secrets are correctly set
- Ensure your repository has Actions enabled
- Check that you have proper permissions (admin/maintain)
