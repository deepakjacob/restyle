package util

import (
	"crypto/rand"
	"encoding/base64"
)

// RandStr method generates random strings
func RandStr() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
