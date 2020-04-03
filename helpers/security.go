package helpers

import (
	"crypto/sha1"
	"encoding/hex"
)

func Encrypt(pwd string) string {
	h := sha1.New()
	h.Write([]byte(pwd))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}
