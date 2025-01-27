package main

import (
	"fmt"

	"golang.org/x/exp/rand"
)

func generateTokens(min, max int) []string {
	tokens := make([]string, rand.Intn(max-min+1)+min)
	for i := range tokens {
		tokens[i] = fmt.Sprintf("Token %d", i+1)
	}
	return tokens
}
