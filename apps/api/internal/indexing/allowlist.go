package indexing

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AllowlistEntry represents a canonical repository approved for indexing.
type AllowlistEntry struct {
	Repository     string  `json:"repository"`
	Enabled        bool    `json:"enabled"`
	IndexingStatus string  `json:"indexingStatus"`
	LastIndexedAt  *string `json:"lastIndexedAt"`
	Notes          string  `json:"notes"`
}

type allowlistContainer struct {
	Repositories []AllowlistEntry `json:"repositories"`
}

// LoadAllowlist reads and parses the canonical data/ai/allowlist.json.
// The canonical allowlist is authoritative. It is never mutated by the indexing runtime.
func LoadAllowlist(dataDir string) ([]AllowlistEntry, error) {
	allowlistPath := filepath.Join(dataDir, "ai", "allowlist.json")
	data, err := os.ReadFile(allowlistPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read canonical allowlist at %s: %w", allowlistPath, err)
	}

	var container allowlistContainer
	if err := json.Unmarshal(data, &container); err != nil {
		return nil, fmt.Errorf("failed to parse canonical allowlist at %s: %w", allowlistPath, err)
	}

	return container.Repositories, nil
}

// IsRepositoryApproved verifies that the target repository is explicitly enabled in the canonical allowlist.
func IsRepositoryApproved(allowlist []AllowlistEntry, repo string) bool {
	norm := strings.TrimSpace(repo)
	for _, entry := range allowlist {
		if strings.EqualFold(strings.TrimSpace(entry.Repository), norm) && entry.Enabled {
			return true
		}
	}
	return false
}

// GetApprovedRepositories returns all repository identifiers that are enabled for indexing.
func GetApprovedRepositories(allowlist []AllowlistEntry) []string {
	var approved []string
	for _, entry := range allowlist {
		if entry.Enabled {
			approved = append(approved, strings.TrimSpace(entry.Repository))
		}
	}
	return approved
}
