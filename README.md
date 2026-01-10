# Stellar Siege

**Defend the Frontier** - A fast-paced space combat game built with Go and Ebiten.

## Overview

Stellar Siege is a top-down space shooter where you defend against waves of increasingly challenging enemies. Features include multiple weapon types, power-ups, boss battles, progressive difficulty, achievements, and an optional online leaderboard.

## Features

- **Dynamic Combat System**
  - 5 weapon types: Minigun, Laser, Plasma, Rockets, Spread Shot
  - Multiple enemy types with unique behaviors
  - Epic boss battles with multiple phases
  - Power-ups and mystery boxes

- **Progression System**
  - Wave-based gameplay with increasing difficulty
  - Achievement system tracking your accomplishments
  - Local and optional online leaderboards
  - Persistent player progression

- **Polish & Effects**
  - Particle effects and explosions
  - Dynamic starfield background
  - Floating damage text
  - Announcements and visual feedback
  - Procedurally generated sound effects

## Installation

### Downloading Pre-built Releases

**Pre-built releases are available for Linux, macOS (Intel & Apple Silicon), and Windows.**

Download the latest release from: https://github.com/sogud/stellar-siege/releases/latest

#### Security Notice

**First-time users may see security warnings** - this is normal for unsigned applications. The game is safe to run, but operating systems display warnings because we don't pay for code signing certificates ($99-400/year).

**How to install safely:**
1. Download the appropriate file for your platform
2. Verify the SHA256 checksum (recommended) - compare with `checksums.txt`
3. Follow the security bypass instructions below

#### Platform-Specific Instructions

- **macOS**: Right-click the app → "Open" → "Open" (bypasses Gatekeeper)
- **Windows**: Click "More info" → "Run anyway" (bypasses SmartScreen)
- **Linux**: No warnings typically appear

**For detailed instructions with screenshots, see [SECURITY.md](SECURITY.md)**

---

## Quick Start

### Prerequisites

- Go 1.24.0 or later
- Platform-specific dependencies:
  - **Linux**: `libc6-dev`, `libgl1-mesa-dev`, `libxcursor-dev`, `libxi-dev`, `libxinerama-dev`, `libxrandr-dev`, `libxxf86vm-dev`, `libasound2-dev`, `pkg-config`
  - **macOS**: Xcode Command Line Tools
  - **Windows**: GCC (e.g., via TDM-GCC or mingw-w64)

### Building from Source

```bash
# Clone the repository
git clone <your-repo-url>
cd space-game

# Install dependencies
go mod download

# Build the game
go build -o stellar-siege .

# Run the game
./stellar-siege
```

### macOS App Bundle

```bash
# Build macOS .app bundle
./build.sh

# Run from Finder
open "Stellar Siege.app"

# Or from terminal
open "Stellar Siege.app"
```

## Controls

- **Arrow Keys** or **WASD**: Move your ship
- **Mouse**: Aim your weapons
- **Left Click**: Fire weapons
- **Space Bar**: Use special ability (when available)
- **ESC**: Pause game / Return to menu

## Game Mechanics

### Weapons

- **Minigun**: Rapid-fire standard bullets
- **Laser**: Continuous beam that pierces enemies
- **Plasma**: Powerful energy projectiles
- **Rockets**: Explosive area damage
- **Spread Shot**: Multiple projectiles in a cone

### Enemies

- **Scouts**: Fast, weak enemies
- **Fighters**: Balanced combat units
- **Dreadnoughts**: Slow, heavily armored
- **Interceptors**: Lightning-fast attackers
- **Bosses**: Unique multi-phase encounters

### Power-Ups

- **Health**: Restore hit points
- **Shield**: Temporary invulnerability
- **Weapon**: Change your armament
- **Speed Boost**: Increased movement speed
- **Mystery Boxes**: Random beneficial effects

## Online Leaderboard (Optional)

Stellar Siege supports an online leaderboard powered by GitHub Gist. This feature is optional and requires setup.

### For Players

If the game includes a `.env` file with leaderboard configuration, it will automatically connect to the online leaderboard. You can disable it by:

1. Deleting the `.env` file, or
2. Setting `GIST_ENABLED=false` in the `.env` file

### For Developers

To set up your own leaderboard:

1. Create a GitHub Gist at https://gist.github.com
2. Create a Personal Access Token with `gist` scope
3. Copy `.env.example` to `.env` and fill in your credentials:

```env
GIST_ID=your_gist_id_here
GH_GIST_TOKEN=your_github_token_here
GIST_ENABLED=true
```

See [ONLINE_LEADERBOARD.md](ONLINE_LEADERBOARD.md) for detailed setup instructions.

## Development

### Project Structure

```
stellar-siege/
├── game/
│   ├── components/      # Shared component types
│   ├── config/          # Game configuration
│   ├── core/            # Core game systems
│   ├── di/              # Dependency injection
│   ├── entities/        # Game entities (player, enemies, projectiles)
│   ├── interfaces/      # Interface definitions
│   ├── states/          # Game state machine
│   └── systems/         # Game systems (rendering, audio, spawning)
├── assets/              # Sprites and resources
├── config/              # Configuration files
├── .github/workflows/   # CI/CD pipelines
└── main.go              # Entry point
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

### Profiling

```bash
# CPU profiling
./stellar-siege -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Memory profiling
./stellar-siege -memprofile=mem.prof
go tool pprof mem.prof

# Live pprof server
./stellar-siege -pprof=:6060
# Visit http://localhost:6060/debug/pprof/
```

## Building Releases

### Manual Release

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

### Automated Release (GitHub Actions)

The project includes GitHub Actions workflows for automated builds and releases:

1. **CI/CD Pipeline**: Runs tests and builds on every push
2. **Release Pipeline**: Creates cross-platform releases on git tags

To create a release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

See [GITHUB_ACTIONS_SETUP.md](GITHUB_ACTIONS_SETUP.md) for complete CI/CD setup instructions.

## Documentation

- **[SECURITY.md](SECURITY.md)** - Security warnings and safe installation guide
- **[GITHUB_ACTIONS_SETUP.md](GITHUB_ACTIONS_SETUP.md)** - Complete CI/CD setup guide
- **[CI_CD_QUICK_REFERENCE.md](CI_CD_QUICK_REFERENCE.md)** - Quick reference for releases
- **[ONLINE_LEADERBOARD.md](ONLINE_LEADERBOARD.md)** - Leaderboard setup for players and developers
- **[LEADERBOARD_CONFIG_FIX.md](LEADERBOARD_CONFIG_FIX.md)** - Troubleshooting leaderboard issues

## Technologies

- **[Go](https://golang.org/)** - Programming language
- **[Ebiten](https://ebiten.org/)** - 2D game engine
- **GitHub Gist** - Optional online leaderboard backend

## Performance

The game is optimized for smooth 60 FPS gameplay:
- Spatial partitioning for collision detection
- Object pooling for entities
- Efficient rendering pipeline
- Minimal allocations in hot paths

## Contributing

Contributions are welcome! Please ensure:

1. Code follows Go conventions (`gofmt`, `go vet`)
2. Tests pass (`go test ./...`)
3. Performance critical code avoids allocations
4. New features include appropriate tests

## License

[Add your license here]

## Credits

Developed using Go and Ebiten game engine.

## Troubleshooting

### Security warnings when launching the game

See [SECURITY.md](SECURITY.md) for complete instructions on safely bypassing macOS Gatekeeper and Windows SmartScreen warnings.

### Game won't start on macOS

If you get a "damaged app" warning:
```bash
xattr -cr "Stellar Siege.app"
```

Or see [SECURITY.md](SECURITY.md) for detailed instructions with screenshots.

### Leaderboard not connecting

1. Check that `.env` file exists with valid credentials
2. Verify your GitHub token has `gist` scope
3. See [LEADERBOARD_CONFIG_FIX.md](LEADERBOARD_CONFIG_FIX.md) for detailed troubleshooting

### Build fails with CGO errors

Ensure you have the required platform dependencies installed (see Prerequisites above).

### Performance issues

1. Check if vsync is enabled (default)
2. Update graphics drivers
3. Try running with profiling to identify bottlenecks (see Profiling section)

## Support

For issues and questions:
- Check the documentation files
- Review GitHub Issues
- Check the troubleshooting section above

---

**Enjoy defending the frontier!** 🚀
