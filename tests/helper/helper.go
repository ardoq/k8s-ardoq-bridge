package helper

import (
	"math/rand"
	"time"
)

func ApplyDelay(seconds ...time.Duration) {
	if len(seconds) > 0 {
		time.Sleep(seconds[0] * time.Second)
	} else {
		time.Sleep(4 * time.Second)
	}
}
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
