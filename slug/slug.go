package slug

import (
	"math/rand"
	"time"
)

func GenerateRandomSlug() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const slugLength = 6

	rand.Seed(time.Now().UnixNano())

	slug := make([]byte, slugLength)
	for i := range slug {
		slug[i] = charset[rand.Intn(len(charset))]
	}

	return string(slug)
}
