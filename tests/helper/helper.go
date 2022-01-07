package helper

import (
	"math/rand"
	"time"
)

func ApplyDelay() {
	time.Sleep(3 * time.Second)
}
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

func RandomString(len int) string {
	byts := make([]byte, len)
	for i := 0; i < len; i++ {
		byts[i] = byte(RandomInt(65, 90))
	}
	return string(byts)
}
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
