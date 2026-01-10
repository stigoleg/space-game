package game

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// Version is the current application version
	// This will be overridden at build time with -ldflags
	Version = "1.0.6"

	// GitHub repository info for update checks
	GitHubOwner = "stigoleg"
	GitHubRepo  = "space-game"
)

// SemanticVersion represents a semantic version (major.minor.patch)
type SemanticVersion struct {
	Major int
	Minor int
	Patch int
}

// ParseVersion parses a version string (e.g., "v1.2.3" or "1.2.3") into a SemanticVersion
func ParseVersion(v string) (SemanticVersion, error) {
	// Remove 'v' prefix if present
	v = strings.TrimPrefix(v, "v")

	// Split by dots
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return SemanticVersion{}, fmt.Errorf("invalid version format: %s (expected major.minor.patch)", v)
	}

	// Parse each component
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return SemanticVersion{}, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return SemanticVersion{}, fmt.Errorf("invalid minor version: %s", parts[1])
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return SemanticVersion{}, fmt.Errorf("invalid patch version: %s", parts[2])
	}

	return SemanticVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// IsNewerThan returns true if this version is newer than the other version
func (v SemanticVersion) IsNewerThan(other SemanticVersion) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	return v.Patch > other.Patch
}

// String returns the version as a string (e.g., "1.2.3")
func (v SemanticVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// StringWithV returns the version with 'v' prefix (e.g., "v1.2.3")
func (v SemanticVersion) StringWithV() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Equals returns true if this version equals the other version
func (v SemanticVersion) Equals(other SemanticVersion) bool {
	return v.Major == other.Major && v.Minor == other.Minor && v.Patch == other.Patch
}
