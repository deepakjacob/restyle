package util

import (
	"encoding/base64"
	"math/rand"
	"time"
)

// RandStr method generates random strings
func RandStr() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// RandNum Method generates random strings
func RandNum() int {
	rand.Seed(time.Now().UnixNano())
	min := 9999
	max := 99999
	return (rand.Intn(max-min+1) + min)
}
