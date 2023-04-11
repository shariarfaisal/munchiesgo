package util

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomBrandType() string {
	brands := []string{"Restaurant", "Store", "Grocery", "Pharmacy"}
	n := len(brands)

	return brands[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(RandomInt(6, 10)))
}
