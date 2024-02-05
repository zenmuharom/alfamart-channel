package helper

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sync"
)

type Hasher struct {
	CounterMu sync.Mutex
	Counter   int
}

func (h *Hasher) GenerateHash(pid, key any) string {
	data := fmt.Sprintf("%v", key)

	// Increment the counter to ensure uniqueness
	h.CounterMu.Lock()
	h.Counter++
	h.CounterMu.Unlock()

	// Combine data with a timestamp to make it unique
	dataWithTimestamp := fmt.Sprintf("%s-%d", pid, data)

	// Calculate SHA-256 hash
	hash := sha256.Sum256([]byte(dataWithTimestamp))

	// Encode the hash to a shorter representation (base64 in this case)
	encodedHash := base64.URLEncoding.EncodeToString(hash[:])

	// Use the first N characters as the short hash (adjust N as needed)
	shortHash := encodedHash[:12] // Use the first 8 characters as an example

	return shortHash
}
