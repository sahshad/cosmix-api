package utils

import (
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