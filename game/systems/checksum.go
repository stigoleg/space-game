package systems

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// VerifyFileChecksum computes SHA256 hash of a file and compares with expected hash
func VerifyFileChecksum(filePath, expectedHash string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Compute SHA256
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, fmt.Errorf("failed to compute hash: %w", err)
	}

	actualHash := hex.EncodeToString(hash.Sum(nil))

	// Compare (case-insensitive)
	return strings.EqualFold(actualHash, expectedHash), nil
}

// ParseChecksumsFile parses checksums.txt content and returns map[filename]hash
func ParseChecksumsFile(content string) map[string]string {
	checksums := make(map[string]string)

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines, headers, and markdown code blocks
		if line == "" || strings.HasPrefix(line, "#") || line == "```" {
			continue
		}

		// Skip lines that don't look like checksums
		if strings.Contains(line, "Verify") || strings.Contains(line, "integrity") {
			continue
		}

		// Format: "hash  filename" (two spaces as separator from sha256sum)
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			hash := parts[0]
			filename := parts[1]

			// Validate hash format (64 hex characters for SHA256)
			if len(hash) == 64 {
				checksums[filename] = hash
			}
		}
	}

	return checksums
}

// FetchAndParseChecksums downloads checksums.txt and extracts the hash for a specific file
func FetchAndParseChecksums(checksumsURL, targetFilename string) (string, error) {
	// Download checksums.txt
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(checksumsURL)
	if err != nil {
		return "", fmt.Errorf("failed to download checksums: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("checksums download failed with status %d", resp.StatusCode)
	}

	// Read content
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read checksums: %w", err)
	}

	// Parse checksums
	checksums := ParseChecksumsFile(string(content))

	// Find hash for target file
	hash, ok := checksums[targetFilename]
	if !ok {
		return "", fmt.Errorf("checksum not found for %s", targetFilename)
	}

	return hash, nil
}

// ComputeFileHash computes SHA256 hash of a file
func ComputeFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to compute hash: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
