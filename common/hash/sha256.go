package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(bytes []byte) string {
	digest := sha256.Sum256(bytes)
	return hex.EncodeToString(digest[:])
}
