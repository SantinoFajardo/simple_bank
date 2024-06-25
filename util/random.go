package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvxyz"

// This function will be called when the package is first use
func init() {
	rand.Seed(time.Now().UnixNano()) // Creat a seed with the current time in Nano secs
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c) // Add the bytes to the the sb var
	}

	return sb.String() // The bytes string is converted to a string
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(10)
}

// RandomMoney generate a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 100)
}

// RandomCurrency generate a random currency symbol
func RandomCurrency() string {
	availableCurrencies := []string{
		EUR, USD, ARS,
	}
	n := len(availableCurrencies)
	return availableCurrencies[rand.Intn(n)]
}
