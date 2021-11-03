package data

import (
	"crypto/sha256"
	"encoding/hex"
)

// ByteDigest creates sha256 hash of content
func ByteDigest(content []byte) string {
	h := sha256.New()
	h.Write(content)
	sha := h.Sum(nil)
	return hex.EncodeToString(sha)
}

// Digest creates sha256 hash of content
func Digest(content string) string {
	return ByteDigest([]byte(content))
}
