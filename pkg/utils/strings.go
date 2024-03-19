package utils

import (
	"math/rand"
)

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func RandomTitle() string {
	return RandomString(50)
}

func RandomDescription() string {
	return RandomString(100)
}

func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomID() int {
	return RandomInt(0, 100)
}

func RandomRating() float64 {
	num := RandomInt(0, 5)
	return float64(num)
}
