package indexing

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// CalculateSHA256 returns the lowercase hexadecimal SHA-256 checksum of a byte slice.
func CalculateSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// CalculateStringSHA256 returns the lowercase hexadecimal SHA-256 checksum of a string.
func CalculateStringSHA256(text string) string {
	return CalculateSHA256([]byte(text))
}

// GenerateSourceID creates a deterministic stable identifier for a knowledge source.
// Format: sha256(repository + ":" + path + ":" + commitSHA)
func GenerateSourceID(repository, path, commitSHA string) string {
	combined := fmt.Sprintf("%s:%s:%s", repository, path, commitSHA)
	return CalculateStringSHA256(combined)
}

// GenerateChunkID creates a deterministic stable identifier for a knowledge chunk.
// Format: sha256(sourceID + ":" + chunkIndex + ":" + chunkChecksum)
func GenerateChunkID(sourceID string, chunkIndex int, chunkChecksum string) string {
	combined := fmt.Sprintf("%s:%d:%s", sourceID, chunkIndex, chunkChecksum)
	return CalculateStringSHA256(combined)
}
