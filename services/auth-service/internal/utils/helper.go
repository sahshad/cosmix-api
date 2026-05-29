package utils

import (
	crypto "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

var usernameRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func GenerateUsername(displayName string) string {
	rand.Seed(time.Now().UnixNano())

	username := strings.ToLower(displayName)

	username = usernameRegex.ReplaceAllString(username, "")

	if username == "" {
		username = "user"
	}

	suffix := rand.Intn(900000) + 100000

	return fmt.Sprintf("%s%d", username, suffix)
}

func GenerateSecureToken(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)

	if _, err := crypto.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
