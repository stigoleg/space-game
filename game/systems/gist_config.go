package systems

import (
	"encoding/json"
	"os"
	"strconv"
)

// GistConfig represents the GitHub Gist configuration
type GistConfig struct {
	GistID      string `json:"gist_id"`
	GitHubToken string `json:"github_token"`
	Enabled     bool   `json:"enabled"`
}

// LoadGistConfig loads the Gist configuration from environment variables first,
// then falls back to JSON file if env vars are not set.
// Priority: Environment Variables > JSON Config File
func LoadGistConfig(filePath string) (*GistConfig, error) {
	// First, try to load from environment variables
	config := &GistConfig{
		GistID:      os.Getenv("GIST_ID"),
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
		Enabled:     parseEnvBool("GIST_ENABLED", false),
	}

	// If env vars are not set, try to load from JSON file
	if config.GistID == "" || config.GitHubToken == "" {
		if filePath == "" {
			filePath = "config/gist_config.json"
		}

		// Try to load JSON config
		if data, err := os.ReadFile(filePath); err == nil {
			var jsonConfig GistConfig
			if err := json.Unmarshal(data, &jsonConfig); err == nil {
				// Only use JSON values if env vars are not set
				if config.GistID == "" {
					config.GistID = jsonConfig.GistID
				}
				if config.GitHubToken == "" {
					config.GitHubToken = jsonConfig.GitHubToken
				}
				// Use JSON enabled flag only if env var wasn't set
				if !config.Enabled && jsonConfig.Enabled {
					config.Enabled = jsonConfig.Enabled
				}
			}
		}
	}

	return config, nil
}

// parseEnvBool parses a boolean environment variable
// Returns the defaultValue if the environment variable is not set
func parseEnvBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}
