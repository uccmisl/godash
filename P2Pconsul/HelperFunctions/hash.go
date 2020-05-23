package HelperFunctions

// Collaborative Code - Start

import (
	"crypto/sha256"
	"encoding/base64"
)

//takes a string input and produces a hashed version
func HashSha(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	s := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return s
}

// Collaborative Code - End
