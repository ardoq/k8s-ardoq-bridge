package helper

import (
	"math/rand"
	"time"
)

func RandomInt(min, max int32) int32 {
	return min + rand.Int31n(max-min)
}

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
