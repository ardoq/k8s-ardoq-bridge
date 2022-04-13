package helper

import (
	"K8SArdoqBridge/app/controllers"
	"math/rand"
	"strings"
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
func CleanupSharedComponents(resourceType string) {
	controllers.Cache.Flush()
	_ = controllers.InitializeCache()
	for k := range controllers.Cache.Items() {
		if strings.HasPrefix(k, "Shared"+resourceType+"Component") {
			s := strings.Split(k, "/")
			_ = controllers.GenericDeleteSharedComponents(resourceType, s[1], s[2])
		}
	}
}
